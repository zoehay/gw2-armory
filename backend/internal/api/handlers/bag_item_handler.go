package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zoehay/gw2armoury/backend/internal/api/models"
	dbmodels "github.com/zoehay/gw2armoury/backend/internal/db/models"
	"github.com/zoehay/gw2armoury/backend/internal/db/repositories"
	"github.com/zoehay/gw2armoury/backend/internal/services"
)

type BagItemHandler struct {
	BagItemRepository repositories.BagItemRepositoryInterface
	ItemService       services.ItemServiceInterface
}

func NewBagItemHandler(bagItemRepository repositories.BagItemRepositoryInterface, itemService services.ItemServiceInterface) *BagItemHandler {
	return &BagItemHandler{
		BagItemRepository: bagItemRepository,
		ItemService:       itemService,
	}
}

func (bagItemHandler BagItemHandler) GetByCharacter(c *gin.Context) {
	characterName := c.Params.ByName("charactername")
	value, exists := c.Get("accountID")
	if !exists {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "could not find Gin Context accountID"})
		return
	}
	accountID := value.(string)

	dbItems, err := bagItemHandler.BagItemRepository.GetDetailBagItemByCharacterName(accountID, characterName)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]models.BagItem, len(dbItems))
	for i := range dbItems {
		items[i] = dbItems[i].ToBagItem()
	}

	c.IndentedJSON(http.StatusOK, items)
}

func (bagItemHandler BagItemHandler) GetByAccount(c *gin.Context) {
	value, exists := c.Get("accountID")
	if !exists {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "could not find Gin Context accountID"})
		return
	}

	accountID := value.(string)
	dbItems, err := bagItemHandler.BagItemRepository.GetDetailBagItemByAccountID(accountID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]models.BagItem, len(dbItems))
	for i := range dbItems {
		items[i] = dbItems[i].ToBagItem()
	}

	c.IndentedJSON(http.StatusOK, items)
}

func (bagItemHandler BagItemHandler) GetAccountInventory(c *gin.Context) {
	value, exists := c.Get("accountID")
	if !exists {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "could not find Gin Context accountID"})
		return
	}
	accountID := value.(string)

	detailBagItems, err := bagItemHandler.BagItemRepository.GetDetailBagItemByAccountID(accountID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error getting account inventory": err.Error()})
		return
	}

	accountInventory, itemsNotInDB := dbmodels.DBDetailBagItemsToAccountInventory(detailBagItems, accountID)

	c.IndentedJSON(http.StatusOK, accountInventory)

	noDuplicates := removeDuplicates(itemsNotInDB)
	itemIDChunks := SplitArray(noDuplicates, 10)
	var errs []error
	for _, idChunk := range itemIDChunks {
		err = bagItemHandler.ItemService.GetAndStoreItemsByID(idChunk)
		if err != nil {
			errs = append(errs, fmt.Errorf("service error getting and storing items in chunk %d: %s", idChunk, err))
		}
	}
	fmt.Println(errs)
}

func (bagItemHandler BagItemHandler) GetFilteredAccountInventory(c *gin.Context) {
	value, exists := c.Get("accountID")
	if !exists {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "could not find Gin Context accountID"})
		return
	}
	accountID := value.(string)

	var searchRequest SearchRequest
	if err := c.BindJSON(&searchRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"request body bind json error": err.Error()})
		return
	}

	searchString := fmt.Sprintf("%%%v%%", searchRequest.SearchTerm)
	detailBagItems, err := bagItemHandler.BagItemRepository.GetDetailBagItemsWithSearch(accountID, searchString)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error getting account inventory": err.Error()})
		return
	}

	accountInventory, _ := dbmodels.DBDetailBagItemsToAccountInventory(detailBagItems, accountID)

	c.IndentedJSON(http.StatusOK, accountInventory)

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

type SearchRequest struct {
	SearchTerm string
}
