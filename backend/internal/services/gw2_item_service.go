package services

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/zoehay/gw2armoury/backend/internal/db/repositories"
	"github.com/zoehay/gw2armoury/backend/internal/gw2_client/providers"
	"gorm.io/gorm"
)

type ItemServiceInterface interface {
	GetAndStoreItemsByID(ids []int) error
	GetAndStoreAllItems() error
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

func (service *ItemService) GetAndStoreItemsByID(ids []int) error {
	apiItems, err := service.ItemProvider.GetItemsByIDs(ids)
	if err != nil {
		return fmt.Errorf("service error using provider: %s", err)
	}

	var errs []error
	for _, item := range apiItems {
		dbItem := item.ToDBItem()
		_, err := service.ItemRepository.Create(&dbItem)
		if err != nil {
			errs = append(errs, fmt.Errorf("GetAndStoreItemsByID: %s", err))
		}
	}

	return errors.Join(errs...)
}

func (service *ItemService) GetAndStoreAllItems() error {
	// allItemIDs, err := service.ItemProvider.GetAllItemIDs()

	// if err != nil {
	// 	return fmt.Errorf("service e rror getting all itemIds: %s", err)
	// }
	allItemIDs := []int{4, 5, 6} // not pre filling db during development

	itemIDChunks := SplitArray(allItemIDs, 50)

	var errs []error

	for _, idChunk := range itemIDChunks {
		err := service.GetAndStoreItemsByID(idChunk)
		if err != nil {
			errs = append(errs, fmt.Errorf("service error getting and storing items in chunk %d: %s", idChunk, err))
		}
	}

	return errors.Join(errs...)
}

func SplitArray(arr []int, chunkSize int) [][]int {
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
