package services

import (
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/gw2_client/providers"
)

type Service struct {
	AccountService *AccountService
	BagItemService *BagItemService
	ItemService    *ItemService
}

func NewService(repository *repositories.Repository, mocks bool) *Service {
	var accountProvider providers.AccountDataProvider
	var characterProvider providers.CharacterDataProvider
	var itemProvider providers.ItemDataProvider

	if mocks {
		accountProvider = &providers.AccountProviderMock{}
		characterProvider = &providers.CharacterProviderMock{}
		itemProvider = &providers.ItemProviderMock{}

	} else {
		accountProvider = &providers.AccountProvider{}
		characterProvider = &providers.CharacterProvider{}
		itemProvider = &providers.ItemProvider{}
	}

	accountService := NewAccountService(&repository.AccountRepository, accountProvider, &repository.SessionRepository)
	bagItemService := NewBagItemService(&repository.BagItemRepository, characterProvider, accountProvider)
	itemService := NewItemService(&repository.ItemRepository, itemProvider)

	return &Service{
		AccountService: accountService,
		BagItemService: bagItemService,
		ItemService:    itemService,
	}
}
