package gw2models

import "github.com/zoehay/gw2-armory/backend/internal/api/models"

type GW2Token struct {
	ID          *string  `json:"id"`
	Name        *string  `json:"name"`
	Permissions []string `json:"permissions"`
}

func (gw2Token GW2Token) ToToken() models.Token {
	return models.Token{
		ID:          *&gw2Token.ID,
		Name:        gw2Token.Name,
		Permissions: gw2Token.Permissions,
	}
}
