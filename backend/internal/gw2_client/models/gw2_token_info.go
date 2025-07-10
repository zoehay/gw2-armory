package gw2models

type GW2TokenInfo struct {
	ID          *string  `json:"id"`
	Name        *string  `json:"name"`
	Permissions []string `json:"permissions"`
}
