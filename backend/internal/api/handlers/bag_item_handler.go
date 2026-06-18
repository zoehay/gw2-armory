package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zoehay/gw2-armory/backend/internal/services"
)

type BagItemHandler struct {
	BagItemService services.BagItemServiceInterface
}

func NewBagItemHandler(bagItemService services.BagItemServiceInterface) *BagItemHandler {
	return &BagItemHandler{
		BagItemService: bagItemService,
	}
}

func (h BagItemHandler) GetByCharacter(c *gin.Context) {
	accountID, ok := getAccountID(c)
	if !ok {
		return
	}
	characterName := c.Params.ByName("charactername")

	items, err := h.BagItemService.GetBagItemsByCharacter(accountID, characterName)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, items)
}

func (h BagItemHandler) GetByAccount(c *gin.Context) {
	accountID, ok := getAccountID(c)
	if !ok {
		return
	}

	items, err := h.BagItemService.GetBagItemsByAccount(accountID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, items)
}

func (h BagItemHandler) GetAccountInventory(c *gin.Context) {
	accountID, ok := getAccountID(c)
	if !ok {
		return
	}

	inventory, missingIDs, err := h.BagItemService.GetAccountInventory(accountID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error getting account inventory": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, inventory)

	if len(missingIDs) > 0 {
		h.BagItemService.FetchMissingItems(missingIDs)
	}
}

func (h BagItemHandler) GetFilteredAccountInventory(c *gin.Context) {
	accountID, ok := getAccountID(c)
	if !ok {
		return
	}

	var searchRequest SearchRequest
	if err := c.BindJSON(&searchRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"request body bind json error": err.Error()})
		return
	}

	inventory, err := h.BagItemService.GetFilteredAccountInventory(accountID, searchRequest.SearchTerm)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error getting account inventory": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, inventory)
}

func getAccountID(c *gin.Context) (string, bool) {
	value, exists := c.Get("accountID")
	if !exists {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "could not find Gin Context accountID"})
		return "", false
	}
	return value.(string), true
}

type SearchRequest struct {
	SearchTerm string
}
