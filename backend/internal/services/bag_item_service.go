package services

import (
	"errors"
	"fmt"

	dbmodels "github.com/zoehay/gw2-armory/backend/internal/db/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	gw2models "github.com/zoehay/gw2-armory/backend/internal/gw2_client/models"
	"github.com/zoehay/gw2-armory/backend/internal/gw2_client/providers"
)

type BagItemServiceInterface interface {
	GetAndStoreAllBagItems(accountID string, apiKey string) error
	GetAndStoreAllCharacters(accountID string, apiKey string) error
	GetAndStoreSharedInventory(accountID string, apiKey string) error
	ClearCharacterInventory(accountID string, characterName string) error
	ClearSharedInventory(accountID string) error
}

type BagItemService struct {
	BagItemRepository *repositories.BagItemRepository
	CharacterProvider providers.CharacterDataProvider
	AccountProvider   providers.AccountDataProvider
}

func NewBagItemService(bagItemRepository *repositories.BagItemRepository, characterProvider providers.CharacterDataProvider, accountProvider providers.AccountDataProvider) *BagItemService {
	return &BagItemService{
		BagItemRepository: bagItemRepository,
		CharacterProvider: characterProvider,
		AccountProvider:   accountProvider,
	}
}

func (service *BagItemService) GetAndStoreAllBagItems(accountID string, apiKey string) error {
	var errs []error
	if err := service.GetAndStoreSharedInventory(accountID, apiKey); err != nil {
		errs = append(errs, fmt.Errorf("GetAndStoreAllBagItems could not get account inventory: %s", err))
	}
	if err := service.GetAndStoreAllCharacters(accountID, apiKey); err != nil {
		errs = append(errs, fmt.Errorf("GetAndStoreAllBagItems could not get character inventory: %s", err))
	}
	return errors.Join(errs...)
}

func (service *BagItemService) GetAndStoreAllCharacters(accountID string, apiKey string) error {
	characters, err := service.CharacterProvider.GetAllCharacters(apiKey)
	if err != nil {
		return fmt.Errorf("service error using provider could not get characters: %s", err)
	}

	for _, character := range characters {
		items := collectCharacterBagItems(accountID, character)
		if err = service.BagItemRepository.ReplaceCharacterInventory(accountID, character.Name, items); err != nil {
			return fmt.Errorf("service error replacing inventory for character %s: %s", character.Name, err)
		}
	}
	return nil
}

func (service *BagItemService) GetAndStoreSharedInventory(accountID string, apiKey string) error {
	accountInventory, err := service.AccountProvider.GetAccountInventory(apiKey)
	if err != nil {
		return fmt.Errorf("service error using provider could not get account inventory: %s", err)
	}

	characterName := "Shared Inventory"
	items := make([]dbmodels.DBBagItem, 0, len(*accountInventory))
	for _, bagItem := range *accountInventory {
		items = append(items, bagItem.ToDBBagItem(accountID, &characterName))
	}

	if err = service.BagItemRepository.ReplaceSharedInventory(accountID, items); err != nil {
		return fmt.Errorf("service error replacing shared inventory for account %s: %s", accountID, err)
	}
	return nil
}

func (service *BagItemService) ClearCharacterInventory(accountID string, characterName string) error {
	if err := service.BagItemRepository.DeleteByCharacterName(accountID, characterName); err != nil {
		return fmt.Errorf("service error deleting bagitems for character %s: %s", characterName, err)
	}
	return nil
}

func (service *BagItemService) ClearSharedInventory(accountID string) error {
	if err := service.BagItemRepository.DeleteSharedInventory(accountID); err != nil {
		return fmt.Errorf("service error deleting shared inventory for account %s: %s", accountID, err)
	}
	return nil
}

func collectCharacterBagItems(accountID string, character gw2models.GW2Character) []dbmodels.DBBagItem {
	var items []dbmodels.DBBagItem
	if character.Bags != nil {
		for _, bag := range *character.Bags {
			for _, bagItem := range bag.Inventory {
				if bagItem != nil {
					items = append(items, bagItem.ToDBBagItem(accountID, &character.Name))
				}
			}
		}
	}
	for _, bagItem := range *character.Equipment {
		items = append(items, bagItem.ToDBBagItem(accountID, &character.Name))
	}
	return items
}
