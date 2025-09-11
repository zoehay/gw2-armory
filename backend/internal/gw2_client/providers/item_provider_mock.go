package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	gw2models "github.com/zoehay/gw2-armory/backend/internal/gw2_client/models"
)

type ItemProviderMock struct{}

func (itemProvider *ItemProviderMock) GetItemsByIDs(intArrIds []int) ([]gw2models.GW2Item, error) {
	wd, _ := os.Getwd()
	isTesting := strings.Contains(wd, "test")
	leadingFilepath := ""

	if isTesting {
		leadingFilepath = "../."
	}

	filepath := fmt.Sprintf("%s./test_data/item_test_data.txt", leadingFilepath)
	apiItems, err := itemProvider.readItemFromFile(filepath)

	if err != nil {
		return nil, fmt.Errorf("error reading from test data file: %s", err)
	}

	return apiItems, nil
}

func (itemProvider *ItemProviderMock) GetAllItemIDs() ([]int, error) {
	return mockAllItemIds, nil
}

func (service *ItemProviderMock) readItemFromFile(filepath string) ([]gw2models.GW2Item, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var items []gw2models.GW2Item
	err = json.Unmarshal(content, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

var mockAllItemIds = []int{
	24,
}
