package gw2client

import (
	"fmt"
	"net/http"
	"time"
)

const baseUrl = "https://api.guildwars2.com/v2/"

var httpClient = &http.Client{Timeout: 20 * time.Second}

func get(path string) (*http.Response, error) {
	return httpClient.Get(baseUrl + path)
}

func GetItemsById(ids string) (*http.Response, error) {
	return get("items?ids=" + ids)
}

func GetItemIds() (*http.Response, error) {
	return get("items")
}

func GetAllCharacters(apiKey string) (*http.Response, error) {
	return get("characters?ids=all&access_token=" + apiKey)
}

func GetTokenInfo(apiKey string) (*http.Response, error) {
	return get("tokeninfo?access_token=" + apiKey)
}

func GetAccount(apiKey string) (*http.Response, error) {
	return get("account?access_token=" + apiKey)
}

func GetAccountInventory(apiKey string) (*http.Response, error) {
	return get("account/inventory?v=latest&access_token=" + apiKey)
}

func GetBankInventory(apiKey string) (*http.Response, error) {
	return get("account/bank?v=latest&access_token=" + apiKey)
}

func GetMaterialsInventory(apiKey string) (*http.Response, error) {
	return get("account/materials?v=latest&access_token=" + apiKey)
}

func GetMaterialCategories() (*http.Response, error) {
	return get("materials?ids=all")
}

func Get(baseUrl string, params map[string]string, headers http.Header) (*http.Response, error) {
	// req, err := http.NewRequest(http.MethodGet, baseUrl, nil)
	// if err != nil {
	// 	return nil, err
	// }

	// query := url.Values{}
	// // query := req.URL.Query()

	// for key, value := range params {
	// 	fmt.Println(key, value)
	// 	if key == "ids" {
	// 		continue
	// 	}
	// 	query.Set(key, value)
	// }
	// fmt.Println(query)
	// req.URL.RawQuery = query.Encode()
	// if value, ok := params["ids"]; ok {
	// 	query.Add("ids", value)
	// }
	// fmt.Println(req.URL.RawQuery)

	req, err := http.NewRequest(http.MethodGet, baseUrl, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range params {
		req.URL.RawQuery += key + value
	}

	fmt.Println(req.URL)

	// req.Header = headers
	// req.Header.Add("Content-Type", `application/json;charset=utf-8`)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil

}
