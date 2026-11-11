package gw2client

import (
	"net/http"
	"time"
)

var baseUrl = "https://api.guildwars2.com/v2/"

var httpClient = &http.Client{Timeout: 20 * time.Second}

// SetBaseURL points the client at another server and returns a func that restores the original.
// For tests in other packages
func SetBaseURL(url string) (restore func()) {
	original := baseUrl
	baseUrl = url
	return func() { baseUrl = original }
}

func get(path string) (*http.Response, error) {
	return httpClient.Get(baseUrl + path)
}

func getWithKey(path string, apiKey string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, baseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	return httpClient.Do(req)
}

func GetItemsById(ids string) (*http.Response, error) {
	return get("items?ids=" + ids)
}

func GetItemIds() (*http.Response, error) {
	return get("items")
}

func GetAllCharacters(apiKey string) (*http.Response, error) {
	return getWithKey("characters?ids=all", apiKey)
}

func GetTokenInfo(apiKey string) (*http.Response, error) {
	return getWithKey("tokeninfo", apiKey)
}

func GetAccount(apiKey string) (*http.Response, error) {
	return getWithKey("account", apiKey)
}

func GetAccountInventory(apiKey string) (*http.Response, error) {
	return getWithKey("account/inventory?v=latest", apiKey)
}

func GetBankInventory(apiKey string) (*http.Response, error) {
	return getWithKey("account/bank?v=latest", apiKey)
}

func GetMaterialsInventory(apiKey string) (*http.Response, error) {
	return getWithKey("account/materials?v=latest", apiKey)
}

func GetMaterialCategories() (*http.Response, error) {
	return get("materials?ids=all")
}
