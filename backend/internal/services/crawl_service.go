package services

import (
	"fmt"

	"golang.org/x/sync/singleflight"
)

type CrawlServiceInterface interface {
	CrawlIfDue(accountID string, apiKey string) error
}

// CrawlService runs at most one inventory crawl per account at a time. Concurrent
// crawls of one account would each replace its inventory and could duplicate items
type CrawlService struct {
	accounts AccountServiceInterface
	bagItems BagItemServiceInterface
	crawls   singleflight.Group
}

func NewCrawlService(accounts AccountServiceInterface, bagItems BagItemServiceInterface) *CrawlService {
	return &CrawlService{accounts: accounts, bagItems: bagItems}
}

// CrawlIfDue crawls the account unless it was crawled recently. Callers arriving
// during a crawl wait for it and share its result
func (service *CrawlService) CrawlIfDue(accountID string, apiKey string) error {
	_, err, _ := service.crawls.Do(accountID, func() (any, error) {
		// read here, a crawl may have finished since the caller loaded the account
		account, err := service.accounts.GetAccountByID(accountID)
		if err != nil {
			return nil, fmt.Errorf("could not load account for crawl: %w", err)
		}
		if !service.accounts.IsRecrawlDue(account.LastCrawl) {
			return nil, nil
		}
		if err := service.bagItems.FetchAndStoreAllBagItems(accountID, apiKey); err != nil {
			return nil, fmt.Errorf("could not crawl inventory: %w", err)
		}
		if err := service.accounts.UpdateLastCrawl(accountID); err != nil {
			return nil, fmt.Errorf("could not update last crawl: %w", err)
		}
		return nil, nil
	})
	return err
}
