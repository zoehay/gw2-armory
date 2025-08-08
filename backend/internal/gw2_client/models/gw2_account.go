package gw2models

import "github.com/zoehay/gw2armoury/backend/internal/api/models"

type GW2Account struct {
	ID   *string `json:"id"`
	Age  *int    `json:"age"`
	Name *string `json:"name"`
}

func (gw2Account GW2Account) ToAccount() models.Account {
	return models.Account{
		AccountID:      *gw2Account.ID,
		GW2AccountName: gw2Account.Name,
	}
}
