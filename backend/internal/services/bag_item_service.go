package services

import (
	"errors"
	"fmt"

	apimodels "github.com/zoehay/gw2-armory/backend/internal/api/models"
	dbmodels "github.com/zoehay/gw2-armory/backend/internal/db/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	gw2models "github.com/zoehay/gw2-armory/backend/internal/gw2_client/models"
	"github.com/zoehay/gw2-armory/backend/internal/gw2_client/providers"
)

type BagItemServiceInterface interface {
	FetchAndStoreAllBagItems(accountID string, apiKey string) error
	FetchAndStoreAllCharacters(accountID string, apiKey string) error
	FetchAndStoreSharedInventory(accountID string, apiKey string) error
	FetchAndStoreBankInventory(accountID string, apiKey string) error
	ClearCharacterInventory(accountID string, characterName string) error
	ClearSharedInventory(accountID string) error
	ClearBankInventory(accountID string) error
	GetBagItemsByCharacter(accountID string, characterName string) ([]apimodels.BagItem, error)
	GetBagItemsByAccount(accountID string) ([]apimodels.BagItem, error)
	GetAccountInventory(accountID string) (apimodels.AccountInventory, []int64, error)
	GetFilteredAccountInventory(accountID string, searchTerm string) (apimodels.AccountInventory, error)
	FetchMissingItems(itemIDs []int64)
}

type BagItemService struct {
	BagItemRepository *repositories.BagItemRepository
	CharacterProvider providers.CharacterDataProvider
	AccountProvider   providers.AccountDataProvider
	ItemService       ItemServiceInterface
}

func NewBagItemService(bagItemRepository *repositories.BagItemRepository, characterProvider providers.CharacterDataProvider, accountProvider providers.AccountDataProvider, itemService ItemServiceInterface) *BagItemService {
	return &BagItemService{
		BagItemRepository: bagItemRepository,
		CharacterProvider: characterProvider,
		AccountProvider:   accountProvider,
		ItemService:       itemService,
	}
}

func (service *BagItemService) FetchAndStoreAllBagItems(accountID string, apiKey string) error {
	var errs []error
	if err := service.FetchAndStoreSharedInventory(accountID, apiKey); err != nil {
		errs = append(errs, fmt.Errorf("FetchAndStoreAllBagItems could not get account inventory: %s", err))
	}
	if err := service.FetchAndStoreAllCharacters(accountID, apiKey); err != nil {
		errs = append(errs, fmt.Errorf("FetchAndStoreAllBagItems could not get character inventory: %s", err))
	}
	return errors.Join(errs...)
}

func (service *BagItemService) FetchAndStoreAllCharacters(accountID string, apiKey string) error {
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

func (service *BagItemService) FetchAndStoreSharedInventory(accountID string, apiKey string) error {
	accountInventory, err := service.AccountProvider.GetAccountInventory(apiKey)
	if err != nil {
		return fmt.Errorf("service error using provider could not get account inventory: %s", err)
	}

	items := make([]dbmodels.DBBagItem, 0, len(*accountInventory))
	for _, bagItem := range *accountInventory {
		items = append(items, bagItem.ToDBBagItem(accountID, nil, "shared"))
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

func (service *BagItemService) FetchAndStoreBankInventory(accountID string, apiKey string) error {
	bankInventory, err := service.AccountProvider.GetBankInventory(apiKey)
	if err != nil {
		return fmt.Errorf("service error using provider could not get bank inventory: %s", err)
	}

	items := make([]dbmodels.DBBagItem, 0, len(*bankInventory))
	for _, bagItem := range *bankInventory {
		items = append(items, bagItem.ToDBBagItem(accountID, nil, "bank"))
	}

	if err = service.BagItemRepository.ReplaceBankInventory(accountID, items); err != nil {
		return fmt.Errorf("service error replacing bank inventory for account %s: %s", accountID, err)
	}
	return nil
}

func (service *BagItemService) ClearBankInventory(accountID string) error {
	if err := service.BagItemRepository.DeleteBankInventory(accountID); err != nil {
		return fmt.Errorf("service error deleting bank inventory for account %s: %s", accountID, err)
	}
	return nil
}

func (service *BagItemService) GetBagItemsByCharacter(accountID string, characterName string) ([]apimodels.BagItem, error) {
	dbItems, err := service.BagItemRepository.GetDetailBagItemByCharacterName(accountID, characterName)
	if err != nil {
		return nil, err
	}
	items := make([]apimodels.BagItem, len(dbItems))
	for i := range dbItems {
		items[i] = dbItems[i].ToBagItem()
	}
	return items, nil
}

func (service *BagItemService) GetBagItemsByAccount(accountID string) ([]apimodels.BagItem, error) {
	dbItems, err := service.BagItemRepository.GetDetailBagItemByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	items := make([]apimodels.BagItem, len(dbItems))
	for i := range dbItems {
		items[i] = dbItems[i].ToBagItem()
	}
	return items, nil
}

func (service *BagItemService) GetAccountInventory(accountID string) (apimodels.AccountInventory, []int64, error) {
	dbItems, err := service.BagItemRepository.GetDetailBagItemByAccountID(accountID)
	if err != nil {
		return apimodels.AccountInventory{}, nil, err
	}
	inventory, itemsNotInDB := dbmodels.DBBagItemsToAccountInventory(dbItems, accountID)
	return inventory, itemsNotInDB, nil
}

func (service *BagItemService) GetFilteredAccountInventory(accountID string, searchTerm string) (apimodels.AccountInventory, error) {
	dbItems, err := service.BagItemRepository.GetDetailBagItemsWithSearch(accountID, searchTerm)
	if err != nil {
		return apimodels.AccountInventory{}, err
	}
	inventory, _ := dbmodels.DBBagItemsToAccountInventory(dbItems, accountID)
	return inventory, nil
}

func (service *BagItemService) FetchMissingItems(itemIDs []int64) {
	noDuplicates := removeDuplicates(itemIDs)
	chunks := splitArray(noDuplicates, 10)
	for _, chunk := range chunks {
		if err := service.ItemService.FetchAndStoreItemsByID(chunk); err != nil {
			fmt.Printf("error fetching missing items: %v\n", err)
		}
	}
}

func removeDuplicates(inputIDs []int64) []int {
	intMap := make(map[int64]bool)
	var noDuplicates []int
	for _, id := range inputIDs {
		if _, value := intMap[id]; !value {
			intMap[id] = true
			noDuplicates = append(noDuplicates, int(id))
		}
	}
	return noDuplicates
}

func splitArray(arr []int, chunkSize int) [][]int {
	var result [][]int
	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		result = append(result, arr[i:end])
	}
	return result
}

func collectCharacterBagItems(accountID string, character gw2models.GW2Character) []dbmodels.DBBagItem {
	var items []dbmodels.DBBagItem
	if character.Bags != nil {
		for _, bag := range *character.Bags {
			for _, bagItem := range bag.Inventory {
				if bagItem != nil {
					items = append(items, bagItem.ToDBBagItem(accountID, &character.Name, "character"))
				}
			}
		}
	}
	for _, bagItem := range *character.Equipment {
		items = append(items, bagItem.ToDBBagItem(accountID, &character.Name, "character"))
	}
	return items
}
