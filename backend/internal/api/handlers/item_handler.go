package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zoehay/gw2-armory/backend/internal/api/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
)

type ItemHandler struct {
	ItemRepository repositories.ItemRepositoryInterface
}

func NewItemHandler(itemRepository repositories.ItemRepositoryInterface) *ItemHandler {
	return &ItemHandler{
		ItemRepository: itemRepository,
	}
}

func (itemHandler ItemHandler) GetAllItems(c *gin.Context) {
	dbItems, err := itemHandler.ItemRepository.GetAll()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]models.Item, len(dbItems))
	for i := range dbItems {
		items[i] = dbItems[i].ToItem()
	}

	c.IndentedJSON(http.StatusOK, items)
}

func (itemHandler ItemHandler) GetItemByID(c *gin.Context) {
	stringId := c.Params.ByName("id")
	itemId, err := strconv.Atoi(stringId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	item, err := itemHandler.ItemRepository.GetById(itemId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, item.ToItem())
}
