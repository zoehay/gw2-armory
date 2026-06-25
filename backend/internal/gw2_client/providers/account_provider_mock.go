package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	gw2models "github.com/zoehay/gw2-armory/backend/internal/gw2_client/models"
)

func testDataPath(filename string) string {
	_, sourceFile, _, _ := runtime.Caller(0)
	backendDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(sourceFile))))
	return filepath.Join(backendDir, "test_data", filename)
}

type AccountProviderMock struct{}

func (accountProvider *AccountProviderMock) GetAccount(apiKey string) (*gw2models.GW2Account, error) {
	account, err := accountProvider.ReadAccountFromFile(testDataPath("account_test_data.txt"))
	if err != nil {
		return nil, fmt.Errorf("error reading from test data file: %s", err)
	}
	return account, nil
}

func (accountProvider *AccountProviderMock) ReadAccountFromFile(filepath string) (*gw2models.GW2Account, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var account gw2models.GW2Account
	err = json.Unmarshal(content, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (accountProvider *AccountProviderMock) GetAccountInventory(apiKey string) (*[]gw2models.GW2BagItem, error) {
	accountInventory, err := accountProvider.ReadAccountInventoryFromFile(testDataPath("account_inventory_test_data.txt"))
	if err != nil {
		return nil, fmt.Errorf("error reading from test data file: %s", err)
	}
	return accountInventory, nil
}

func (accountProvider *AccountProviderMock) ReadAccountInventoryFromFile(filepath string) (*[]gw2models.GW2BagItem, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var accountInventory *[]gw2models.GW2BagItem
	err = json.Unmarshal(content, &accountInventory)
	if err != nil {
		return nil, err
	}

	return accountInventory, nil
}

func (accountProvider *AccountProviderMock) GetBankInventory(apiKey string) (*[]gw2models.GW2BagItem, error) {
	bankInventory, err := accountProvider.ReadBankInventoryFromFile(testDataPath("account_bank_test_data.txt"))
	if err != nil {
		return nil, fmt.Errorf("error reading from test data file: %s", err)
	}
	return bankInventory, nil
}

func (accountProvider *AccountProviderMock) ReadBankInventoryFromFile(filepath string) (*[]gw2models.GW2BagItem, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var bankInventory *[]gw2models.GW2BagItem
	err = json.Unmarshal(content, &bankInventory)
	if err != nil {
		return nil, err
	}

	return bankInventory, nil
}

func (accountProvider *AccountProviderMock) GetTokenInfo(apiKey string) (*gw2models.GW2Token, error) {
	token, err := accountProvider.ReadTokenInfoFromFile(testDataPath("token_info_test_data.txt"))
	if err != nil {
		return nil, fmt.Errorf("error reading from test data file: %s", err)
	}
	return token, nil
}

func (accountProvider *AccountProviderMock) ReadTokenInfoFromFile(filepath string) (*gw2models.GW2Token, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var token gw2models.GW2Token
	err = json.Unmarshal(content, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
