package models

type Token struct {
	ID          *string  `json:"id"`
	Name        *string  `json:"name"`
	Permissions []string `json:"permissions"`
}
