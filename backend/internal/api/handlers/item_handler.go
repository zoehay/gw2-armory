package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zoehay/gw2-armory/backend/internal/services"
)

type ItemHandler struct {
	ItemService services.ItemServiceInterface
}

func NewItemHandler(itemService services.ItemServiceInterface) *ItemHandler {
	return &ItemHandler{
		ItemService: itemService,
	}
}

func (h ItemHandler) GetAllItems(c *gin.Context) {
	items, err := h.ItemService.GetAllItems()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not get items", err)
		return
	}
	c.IndentedJSON(http.StatusOK, items)
}

func (h ItemHandler) GetItemByID(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid item id", err)
		return
	}

	item, err := h.ItemService.GetItemByID(itemID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not get item", err)
		return
	}
	c.IndentedJSON(http.StatusOK, item)
}
