package tests

import (
	"fmt"

	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
)

func GetItemsInCharacterBags(itemService *services.ItemService, bagItemRepository repositories.BagItemRepository) error {
	itemIds, err := bagItemRepository.GetIds()
	fmt.Println("number of items", len(itemIds))
	if err != nil {
		return err
	}

	if len(itemIds) > 500 {
		return fmt.Errorf("dev setup more than 500 items, maybe wait")
	}

	err = itemService.GetAndStoreEachByIDs(itemIds)
	if err != nil {
		return fmt.Errorf("dev setup error getting and storing items : %s", err)
	}

	return nil
}
