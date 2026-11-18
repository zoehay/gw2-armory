package services

import (
	"errors"
	"slices"
	"sync"
	"testing"
	"time"
)

type stubItemService struct {
	ItemServiceInterface

	mu     sync.Mutex
	calls  [][]int
	handle func(ids []int) ([]int, error)
}

func (s *stubItemService) FetchAndStoreItemsByID(ids []int) ([]int, error) {
	s.mu.Lock()
	s.calls = append(s.calls, slices.Clone(ids))
	s.mu.Unlock()
	if s.handle != nil {
		return s.handle(ids)
	}
	return ids, nil
}

func (s *stubItemService) recordedCalls() [][]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.calls)
}

type fakeClock struct {
	mu      sync.Mutex
	current time.Time
}

func (c *fakeClock) now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.current
}

func (c *fakeClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.current = c.current.Add(d)
}

func newTestFetcher(items *stubItemService) (*MissingItemFetcher, *fakeClock) {
	clock := &fakeClock{current: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	fetcher := NewMissingItemFetcher(items)
	fetcher.now = clock.now
	return fetcher, clock
}

func enqueueAndWait(fetcher *MissingItemFetcher, ids ...int64) {
	fetcher.Enqueue(ids)
	fetcher.Wait()
}

func sequentialIDs(n int) []int64 {
	ids := make([]int64, n)
	for i := range ids {
		ids[i] = int64(i + 1)
	}
	return ids
}

func TestFetcherUsesFullBatches(t *testing.T) {
	items := &stubItemService{}
	fetcher, _ := newTestFetcher(items)

	enqueueAndWait(fetcher, sequentialIDs(450)...)

	var sizes []int
	for _, call := range items.recordedCalls() {
		sizes = append(sizes, len(call))
	}
	if want := []int{200, 200, 50}; !slices.Equal(sizes, want) {
		t.Errorf("request sizes = %v, want %v", sizes, want)
	}
}

func TestFetcherDeduplicates(t *testing.T) {
	items := &stubItemService{}
	fetcher, _ := newTestFetcher(items)

	enqueueAndWait(fetcher, 7, 3, 7, 3, 9)

	calls := items.recordedCalls()
	if len(calls) != 1 || !slices.Equal(calls[0], []int{7, 3, 9}) {
		t.Errorf("calls = %v, want one call with [7 3 9]", calls)
	}
}

func TestFetcherNoIDsMakesNoRequest(t *testing.T) {
	items := &stubItemService{}
	fetcher, _ := newTestFetcher(items)

	enqueueAndWait(fetcher)

	if calls := items.recordedCalls(); len(calls) != 0 {
		t.Errorf("calls = %v, want none", calls)
	}
}

func TestFetcherEnqueueDoesNotWaitForFetch(t *testing.T) {
	release := make(chan struct{})
	items := &stubItemService{handle: func(ids []int) ([]int, error) {
		<-release
		return ids, nil
	}}
	fetcher, _ := newTestFetcher(items)

	done := make(chan struct{})
	go func() {
		fetcher.Enqueue([]int64{1})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Error("Enqueue blocked on the GW2 call")
	}
	close(release)
	fetcher.Wait()
}

func TestFetcherSkipsIDsAlreadyInFlight(t *testing.T) {
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	items := &stubItemService{handle: func(ids []int) ([]int, error) {
		started <- struct{}{}
		<-release
		return ids, nil
	}}
	fetcher, _ := newTestFetcher(items)

	fetcher.Enqueue([]int64{1, 2})
	<-started
	fetcher.Enqueue([]int64{2, 3})
	close(release)
	fetcher.Wait()

	calls := items.recordedCalls()
	if len(calls) != 2 || !slices.Equal(calls[0], []int{1, 2}) || !slices.Equal(calls[1], []int{3}) {
		t.Errorf("calls = %v, want [[1 2] [3]]", calls)
	}
}

func TestFetcherSkipsUnknownIDsForADay(t *testing.T) {
	items := &stubItemService{handle: func(ids []int) ([]int, error) {
		return nil, nil
	}}
	fetcher, clock := newTestFetcher(items)

	enqueueAndWait(fetcher, 1)
	clock.advance(unknownItemRetryAfter - time.Second)
	enqueueAndWait(fetcher, 1)
	if calls := items.recordedCalls(); len(calls) != 1 {
		t.Fatalf("made %d requests before the skip expired, want 1", len(calls))
	}

	clock.advance(time.Second)
	enqueueAndWait(fetcher, 1)
	if calls := items.recordedCalls(); len(calls) != 2 {
		t.Errorf("made %d requests after the skip expired, want 2", len(calls))
	}
}

func TestFetcherOnlySkipsIDsNotStored(t *testing.T) {
	items := &stubItemService{handle: func(ids []int) ([]int, error) {
		return []int{1}, nil
	}}
	fetcher, _ := newTestFetcher(items)

	enqueueAndWait(fetcher, 1, 2)
	enqueueAndWait(fetcher, 1, 2)

	calls := items.recordedCalls()
	if len(calls) != 2 || !slices.Equal(calls[1], []int{1}) {
		t.Errorf("calls = %v, want the stored ID requested again and the unknown ID skipped", calls)
	}
}

func TestFetcherRetriesFailedIDsAfterAMinute(t *testing.T) {
	items := &stubItemService{handle: func(ids []int) ([]int, error) {
		return nil, errors.New("gw2 unavailable")
	}}
	fetcher, clock := newTestFetcher(items)

	enqueueAndWait(fetcher, 1)
	clock.advance(failedItemRetryAfter - time.Second)
	enqueueAndWait(fetcher, 1)
	if calls := items.recordedCalls(); len(calls) != 1 {
		t.Fatalf("made %d requests before the cooldown expired, want 1", len(calls))
	}

	clock.advance(time.Second)
	enqueueAndWait(fetcher, 1)
	if calls := items.recordedCalls(); len(calls) != 2 {
		t.Errorf("made %d requests after the cooldown expired, want 2", len(calls))
	}
}

func TestFetcherContinuesAfterBatchError(t *testing.T) {
	items := &stubItemService{handle: func(ids []int) ([]int, error) {
		return nil, errors.New("gw2 unavailable")
	}}
	fetcher, _ := newTestFetcher(items)

	enqueueAndWait(fetcher, sequentialIDs(401)...)

	if calls := items.recordedCalls(); len(calls) != 3 {
		t.Errorf("made %d requests, want 3 even when batches fail", len(calls))
	}
}

func TestFetcherSurvivesPanic(t *testing.T) {
	items := &stubItemService{handle: func(ids []int) ([]int, error) {
		if ids[0] == 1 {
			panic("boom")
		}
		return ids, nil
	}}
	fetcher, _ := newTestFetcher(items)

	enqueueAndWait(fetcher, 1)
	enqueueAndWait(fetcher, 2)

	if calls := items.recordedCalls(); len(calls) != 2 {
		t.Errorf("made %d requests, want the worker to keep working after a panic", len(calls))
	}
}
