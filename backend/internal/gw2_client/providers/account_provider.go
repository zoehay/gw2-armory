package providers

import (
	"encoding/json"
	"fmt"
	"io"

	gw2client "github.com/zoehay/gw2-armory/backend/internal/gw2_client"
	gw2models "github.com/zoehay/gw2-armory/backend/internal/gw2_client/models"
)

type AccountDataProvider interface {
	GetAccount(apiKey string) (*gw2models.GW2Account, error)
	GetAccountInventory(apiKey string) (*[]gw2models.GW2BagItem, error)
	GetTokenInfo(apiKey string) (*gw2models.GW2TokenInfo, error)
}

type AccountProvider struct{}

func (accountProvider *AccountProvider) GetAccount(apiKey string) (*gw2models.GW2Account, error) {
	res, err := gw2client.GetAccount(apiKey)

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

	var result *gw2models.GW2Account

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("provider json.Unmarshal error: %s", err)
	}

	return result, nil

}

func (accountProvider *AccountProvider) GetAccountInventory(apiKey string) (*[]gw2models.GW2BagItem, error) {
	res, err := gw2client.GetAccountInventory(apiKey)

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

	var result *[]gw2models.GW2BagItem

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("provider json.Unmarshal error: %s", err)
	}

	return result, nil

}

func (accountProvider *AccountProvider) GetTokenInfo(apiKey string) (*gw2models.GW2TokenInfo, error) {
	res, err := gw2client.GetTokenInfo(apiKey)

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

	var result *gw2models.GW2TokenInfo

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("provider json.Unmarshal error: %s", err)
	}

	return result, nil

}
