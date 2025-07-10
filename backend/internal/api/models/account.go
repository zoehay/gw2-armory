package models

import "time"

type Account struct {
	AccountID      string     `json:"id"`
	LastCrawl      *time.Time `json:"last_crawl"`
	AccountName    *string    `json:"name"`
	GW2AccountName *string    `json:"gw2_name"`
	GW2TokenName   *string    `json:"gw2_token_name"`
	APIKey         *string    `json:"api_key"`
	Password       *string    `json:"password"`
	SessionID      *string    `json:"session_id"`
	Session        *Session   `json:"session"`
}
