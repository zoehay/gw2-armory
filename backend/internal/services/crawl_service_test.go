package services

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	apimodels "github.com/zoehay/gw2-armory/backend/internal/api/models"
)

type stubCrawlAccounts struct {
	AccountServiceInterface

	mu        sync.Mutex
	lastCrawl map[string]*time.Time
	updateErr error
}

func (s *stubCrawlAccounts) GetAccountByID(accountID string) (*apimodels.Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &apimodels.Account{AccountID: accountID, LastCrawl: s.lastCrawl[accountID]}, nil
}

func (s *stubCrawlAccounts) IsRecrawlDue(lastCrawl *time.Time) bool {
	return lastCrawl == nil
}

func (s *stubCrawlAccounts) UpdateLastCrawl(accountID string) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	if s.lastCrawl == nil {
		s.lastCrawl = make(map[string]*time.Time)
	}
	s.lastCrawl[accountID] = &now
	return nil
}

type stubCrawlBagItems struct {
	BagItemServiceInterface

	crawls atomic.Int32
	handle func(accountID string) error
}

func (s *stubCrawlBagItems) FetchAndStoreAllBagItems(accountID string, apiKey string) error {
	s.crawls.Add(1)
	if s.handle != nil {
		return s.handle(accountID)
	}
	return nil
}

func TestCrawlIfDueCrawlsAndRecordsTheCrawl(t *testing.T) {
	accounts := &stubCrawlAccounts{}
	bagItems := &stubCrawlBagItems{}
	service := NewCrawlService(accounts, bagItems)

	if err := service.CrawlIfDue("account", "key"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := bagItems.crawls.Load(); got != 1 {
		t.Errorf("crawled %d times, want 1", got)
	}
	if accounts.lastCrawl["account"] == nil {
		t.Error("last crawl was not recorded")
	}
}

func TestCrawlIfDueSkipsRecentCrawl(t *testing.T) {
	bagItems := &stubCrawlBagItems{}
	service := NewCrawlService(&stubCrawlAccounts{}, bagItems)

	for range 2 {
		if err := service.CrawlIfDue("account", "key"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if got := bagItems.crawls.Load(); got != 1 {
		t.Errorf("crawled %d times, want 1", got)
	}
}

func TestCrawlIfDueConcurrentCallsCrawlOnce(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	bagItems := &stubCrawlBagItems{handle: func(string) error {
		once.Do(func() { close(started) })
		<-release
		return nil
	}}
	service := NewCrawlService(&stubCrawlAccounts{}, bagItems)

	var callers sync.WaitGroup
	for range 5 {
		callers.Go(func() {
			if err := service.CrawlIfDue("account", "key"); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}

	<-started
	time.Sleep(50 * time.Millisecond) // let the other callers reach the crawl
	close(release)
	callers.Wait()

	if got := bagItems.crawls.Load(); got != 1 {
		t.Errorf("crawled %d times, want 1", got)
	}
}

func TestCrawlIfDueSharesErrorWithWaitingCallers(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	crawlErr := errors.New("gw2 unavailable")
	bagItems := &stubCrawlBagItems{handle: func(string) error {
		once.Do(func() { close(started) })
		<-release
		return crawlErr
	}}
	service := NewCrawlService(&stubCrawlAccounts{}, bagItems)

	errs := make([]error, 3)
	var callers sync.WaitGroup
	for i := range errs {
		callers.Go(func() {
			errs[i] = service.CrawlIfDue("account", "key")
		})
	}

	<-started
	time.Sleep(50 * time.Millisecond) // let the other callers reach the crawl
	close(release)
	callers.Wait()

	for i, err := range errs {
		if !errors.Is(err, crawlErr) {
			t.Errorf("caller %d: err = %v, want the crawl error", i, err)
		}
	}
}

func TestCrawlIfDueFailedCrawlIsNotRecorded(t *testing.T) {
	accounts := &stubCrawlAccounts{}
	bagItems := &stubCrawlBagItems{handle: func(string) error {
		return errors.New("gw2 unavailable")
	}}
	service := NewCrawlService(accounts, bagItems)

	if err := service.CrawlIfDue("account", "key"); err == nil {
		t.Fatal("expected an error")
	}
	if accounts.lastCrawl["account"] != nil {
		t.Error("a failed crawl should stay due")
	}
}

func TestCrawlIfDueReportsLastCrawlUpdateError(t *testing.T) {
	updateErr := errors.New("db unavailable")
	service := NewCrawlService(&stubCrawlAccounts{updateErr: updateErr}, &stubCrawlBagItems{})

	if err := service.CrawlIfDue("account", "key"); !errors.Is(err, updateErr) {
		t.Errorf("err = %v, want the update error", err)
	}
}

func TestCrawlIfDueDifferentAccountsRunIndependently(t *testing.T) {
	blockedStarted := make(chan struct{})
	release := make(chan struct{})
	bagItems := &stubCrawlBagItems{handle: func(accountID string) error {
		if accountID == "blocked" {
			close(blockedStarted)
			<-release
		}
		return nil
	}}
	service := NewCrawlService(&stubCrawlAccounts{}, bagItems)

	blockedDone := make(chan error, 1)
	go func() { blockedDone <- service.CrawlIfDue("blocked", "key") }()
	<-blockedStarted

	otherDone := make(chan error, 1)
	go func() { otherDone <- service.CrawlIfDue("other", "key") }()

	select {
	case err := <-otherDone:
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("one account's crawl blocked another account")
	}

	close(release)
	if err := <-blockedDone; err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
