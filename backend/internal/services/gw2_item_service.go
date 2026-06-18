package services

import (
	"errors"
	"fmt"
	"strconv"

	apimodels "github.com/zoehay/gw2-armory/backend/internal/api/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/gw2_client/providers"
	"gorm.io/gorm"
)

type ItemServiceInterface interface {
	FetchAndStoreItemsByID(ids []int) error
	FetchAndStoreAllItems() error
	GetAllItems() ([]apimodels.Item, error)
	GetItemByID(id int) (*apimodels.Item, error)
}

type ItemService struct {
	ItemRepository *repositories.ItemRepository
	ItemProvider   providers.ItemDataProvider
}

func NewItemService(itemRepository *repositories.ItemRepository, itemProvider providers.ItemDataProvider) *ItemService {
	return &ItemService{
		ItemRepository: itemRepository,
		ItemProvider:   itemProvider,
	}
}

func (service *ItemService) FetchAndStoreItemsByID(ids []int) error {
	apiItems, err := service.ItemProvider.GetItemsByIDs(ids)
	if err != nil {
		return fmt.Errorf("service error using provider: %s", err)
	}

	var errs []error
	for _, item := range apiItems {
		dbItem := item.ToDBItem()
		_, err := service.ItemRepository.Create(&dbItem)
		if err != nil {
			errs = append(errs, fmt.Errorf("FetchAndStoreItemsByID: %s", err))
		}
	}

	return errors.Join(errs...)
}

func (service *ItemService) FetchAndStoreAllItems() error {
	// allItemIDs, err := service.ItemProvider.GetAllItemIDs()

	// if err != nil {
	// 	return fmt.Errorf("service e rror getting all itemIds: %s", err)
	// }
	allItemIDs := []int{4, 5, 6} // not pre filling db during development

	itemIDChunks := splitArray(allItemIDs, 50)

	var errs []error

	for _, idChunk := range itemIDChunks {
		err := service.FetchAndStoreItemsByID(idChunk)
		if err != nil {
			errs = append(errs, fmt.Errorf("service error getting and storing items in chunk %d: %s", idChunk, err))
		}
	}

	return errors.Join(errs...)
}

func (service *ItemService) GetAllItems() ([]apimodels.Item, error) {
	dbItems, err := service.ItemRepository.GetAll()
	if err != nil {
		return nil, err
	}
	items := make([]apimodels.Item, len(dbItems))
	for i := range dbItems {
		items[i] = dbItems[i].ToItem()
	}
	return items, nil
}

func (service *ItemService) GetItemByID(id int) (*apimodels.Item, error) {
	dbItem, err := service.ItemRepository.GetById(id)
	if err != nil {
		return nil, err
	}
	item := dbItem.ToItem()
	return &item, nil
}

func (service *ItemService) GetAndStoreEachByIDs(itemIds []int) error {
	apiItems, err := service.ItemProvider.GetItemsByIDs(itemIds)
	if err != nil {
		return fmt.Errorf("provider error requesting items: %s", err)
	}

	var duplicateKeyErrorIDs []uint
	for _, item := range apiItems {
		dbItem := item.ToDBItem()
		_, err := service.ItemRepository.Create(&dbItem)
		if err != nil {
			if isDuplicateKeyError(err) {
				duplicateKeyErrorIDs = append(duplicateKeyErrorIDs, item.ID)
			} else {
				return fmt.Errorf("gorm error adding item id %d: %s", item.ID, err)
			}
		}
	}

	if len(duplicateKeyErrorIDs) != 0 {
		fmt.Printf("skipped adding duplicate values %#v\n", duplicateKeyErrorIDs)
	}

	return nil
}

func isDuplicateKeyError(err error) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey)
}

func IntArrToStringArr(intArr []int) []string {
	var stringArr []string
	for _, num := range intArr {
		stringArr = append(stringArr, strconv.Itoa(num))
	}
	return stringArr

}
