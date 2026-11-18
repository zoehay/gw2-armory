package services

import (
	"fmt"
	"log"
	"slices"
	"sync"
	"time"
)

const (
	// unknownItemRetryAfter is how long to skip an ID GW2 didn't return
	unknownItemRetryAfter = 24 * time.Hour
	// failedItemRetryAfter keeps a GW2 outage from causing retries on every page load
	failedItemRetryAfter = time.Minute
)

// MissingItemFetcher stores items found in inventories but not in the database.
// A single background worker drains the queue, so requests from different users
// share batches and only one GW2 call runs at a time
type MissingItemFetcher struct {
	items ItemServiceInterface
	now   func() time.Time

	mu        sync.Mutex
	queue     []int
	claimed   map[int]struct{} // queued or being fetched
	skipUntil map[int]time.Time
	running   bool
	worker    sync.WaitGroup
}

func NewMissingItemFetcher(items ItemServiceInterface) *MissingItemFetcher {
	return &MissingItemFetcher{
		items:     items,
		now:       time.Now,
		claimed:   make(map[int]struct{}),
		skipUntil: make(map[int]time.Time),
	}
}

// Enqueue queues the IDs that need fetching and returns without waiting for them
func (f *MissingItemFetcher) Enqueue(itemIDs []int64) {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := f.now()
	for _, itemID := range itemIDs {
		id := int(itemID)
		if _, ok := f.claimed[id]; ok {
			continue
		}
		if until, ok := f.skipUntil[id]; ok {
			if now.Before(until) {
				continue
			}
			delete(f.skipUntil, id)
		}
		f.claimed[id] = struct{}{}
		f.queue = append(f.queue, id)
	}

	if len(f.queue) > 0 && !f.running {
		f.running = true
		f.worker.Add(1)
		go f.run()
	}
}

// Wait blocks until the worker has drained the queue
func (f *MissingItemFetcher) Wait() {
	f.worker.Wait()
}

func (f *MissingItemFetcher) run() {
	defer f.worker.Done()
	for {
		batch := f.nextBatch()
		if batch == nil {
			return
		}
		storedIDs, err := f.fetch(batch)
		if err != nil {
			log.Printf("error fetching missing items %v: %v", batch, err)
		}
		f.finish(batch, storedIDs, err)
	}
}

func (f *MissingItemFetcher) nextBatch() []int {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.queue) == 0 {
		f.running = false
		return nil
	}
	size := min(len(f.queue), maxItemIDsPerRequest)
	batch := slices.Clone(f.queue[:size])
	f.queue = f.queue[size:]
	return batch
}

// fetch recovers panics because gin.Recovery doesn't cover this goroutine
func (f *MissingItemFetcher) fetch(batch []int) (storedIDs []int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return f.items.FetchAndStoreItemsByID(batch)
}

func (f *MissingItemFetcher) finish(batch []int, storedIDs []int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	retryAfter := unknownItemRetryAfter
	if err != nil {
		retryAfter = failedItemRetryAfter
	}
	until := f.now().Add(retryAfter)

	for _, id := range batch {
		delete(f.claimed, id)
		if !slices.Contains(storedIDs, id) {
			f.skipUntil[id] = until
		}
	}
}
