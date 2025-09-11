package services

import (
	"errors"
	"fmt"

	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	gw2models "github.com/zoehay/gw2-armory/backend/internal/gw2_client/models"
	"github.com/zoehay/gw2-armory/backend/internal/gw2_client/providers"
)

type BagItemServiceInterface interface {
	GetAndStoreAllBagItems(accountID string, apiKey string) error

	GetAndStoreAllCharacters(accountID string, apiKey string) error
	GetAndStoreSharedInventory(accountID string, apiKey string) error

	StoreCharacterInventory(accountID string, character gw2models.GW2Character) error
	StoreSharedInventory(accountID string, accountInventory *[]gw2models.GW2BagItem) error

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
	err := service.GetAndStoreSharedInventory(accountID, apiKey)
	if err != nil {
		errs = append(errs, fmt.Errorf("GetAndStoreAllBagItems could not get account inventory: %s", err))
	}

	err = service.GetAndStoreAllCharacters(accountID, apiKey)
	if err != nil {
		errs = append(errs, fmt.Errorf("GetAndStoreAllBagItems could not get character inventory: %s", err))
	}

	return errors.Join(errs...)
}

func (service *BagItemService) GetAndStoreAllCharacters(accountID string, apiKey string) error {
	characters, err := service.CharacterProvider.GetAllCharacters(apiKey)
	if err != nil {
		return fmt.Errorf("service error using provider could not get characters: %s", err)
	}

	tx := service.BagItemRepository.DB.Begin()

	defer func() {
		r := recover()
		if r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, character := range characters {
		err = service.ClearCharacterInventory(accountID, character.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = service.StoreCharacterInventory(accountID, character)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (service *BagItemService) GetAndStoreSharedInventory(accountID string, apiKey string) error {

	accountInventory, err := service.AccountProvider.GetAccountInventory(apiKey)
	if err != nil {
		return fmt.Errorf("service error using provider could not get account inventory: %s", err)
	}

	tx := service.BagItemRepository.DB.Begin()

	defer func() {
		r := recover()
		if r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err = service.ClearSharedInventory(accountID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = service.StoreSharedInventory(accountID, accountInventory)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

func (service *BagItemService) StoreCharacterInventory(accountID string, character gw2models.GW2Character) error {
	apiBags := character.Bags

	if apiBags != nil {
		for _, bag := range *apiBags {
			for _, bagItem := range bag.Inventory {
				if bagItem != nil {
					dbBagItem := bagItem.ToDBBagItem(accountID, &character.Name)
					_, err := service.BagItemRepository.Create(&dbBagItem)
					if err != nil {
						return fmt.Errorf("service error using create bagitem %d for character %s: %s", bagItem.ID, character.Name, err)
					}
				}
			}
		}
	}

	equipment := character.Equipment
	for _, bagItem := range *equipment {
		dbBagItem := bagItem.ToDBBagItem(accountID, &character.Name)
		_, err := service.BagItemRepository.Create(&dbBagItem)
		if err != nil {
			return fmt.Errorf("service error using create bagitem %d for character %s: %s", bagItem.ID, character.Name, err)
		}
	}

	return nil
}

func (service *BagItemService) StoreSharedInventory(accountID string, accountInventory *[]gw2models.GW2BagItem) error {
	characterName := "Shared Inventory"
	for _, bagItem := range *accountInventory {
		dbBagItem := bagItem.ToDBBagItem(accountID, &characterName)
		_, err := service.BagItemRepository.Create(&dbBagItem)
		if err != nil {
			return fmt.Errorf("service error using create bagitem %d for account %s: %s", bagItem.ID, accountID, err)
		}
	}
	return nil
}

func (service *BagItemService) ClearCharacterInventory(accountID string, characterName string) error {
	err := service.BagItemRepository.DeleteByCharacterName(accountID, characterName)
	if err != nil {
		return fmt.Errorf("service error using delete bagitems for character %s: %s", characterName, err)
	}

	return nil
}

func (service *BagItemService) ClearSharedInventory(accountID string) error {
	err := service.BagItemRepository.DeleteSharedInventory(accountID)
	if err != nil {
		return fmt.Errorf("service error using delete bagitems for account %s: %s", accountID, err)
	}

	return nil
}
