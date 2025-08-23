package models

import "time"

type Account struct {
	AccountID      string     `json:"id"`
	LastCrawl      *time.Time `json:"last_crawl,omitempty"`
	AccountName    *string    `json:"name,omitempty"`
	GW2AccountName *string    `json:"gw2_name,omitempty"`
	GW2TokenName   *string    `json:"gw2_token_name,omitempty"`
	APIKey         *string    `json:"api_key,omitempty"`
	Password       *string    `json:"password,omitempty"`
	SessionID      *string    `json:"session_id,omitempty"`
	Session        *Session   `json:"session,omitempty"`
}
