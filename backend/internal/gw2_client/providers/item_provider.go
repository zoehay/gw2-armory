package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	gw2client "github.com/zoehay/gw2-armory/backend/internal/gw2_client"
	gw2models "github.com/zoehay/gw2-armory/backend/internal/gw2_client/models"
)

type ItemDataProvider interface {
	GetItemsByIDs(intArrIds []int) ([]gw2models.GW2Item, error)
	GetAllItemIDs() ([]int, error)
	GetMaterialCategories() ([]gw2models.GW2MaterialCategory, error)
}

type ItemProvider struct{}

func (itemProvider *ItemProvider) GetItemsByIDs(intArrIds []int) ([]gw2models.GW2Item, error) {
	idString := strings.Join(IntArrToStringArr(intArrIds), ",")
	res, err := gw2client.GetItemsById(idString)

	if err != nil {
		return nil, fmt.Errorf("provider get error: %s", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("provider io.ReadAll error: %s", err)
	}

	var result []gw2models.GW2Item
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("provider json.Unmarshal error: %s", err)
	}

	return result, nil
}

func (itemProvider *ItemProvider) GetAllItemIDs() ([]int, error) {
	res, err := gw2client.GetItemIds()

	if err != nil {
		return nil, fmt.Errorf("provider get error: %s", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("provider io.ReadAll error: %s", err)
	}

	var result []int

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("provider json.Unmarshal error: %s", err)
	}

	return result, nil

}

func (itemProvider *ItemProvider) GetMaterialCategories() ([]gw2models.GW2MaterialCategory, error) {
	res, err := gw2client.GetMaterialCategories()

	if err != nil {
		return nil, fmt.Errorf("provider get error: %s", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("provider io.ReadAll error: %s", err)
	}

	var result []gw2models.GW2MaterialCategory
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("provider json.Unmarshal error: %s", err)
	}

	return result, nil
}

func IntArrToStringArr(intArr []int) []string {
	var stringArr []string
	for _, i := range intArr {
		stringArr = append(stringArr, strconv.Itoa(i))
	}
	return stringArr
}
