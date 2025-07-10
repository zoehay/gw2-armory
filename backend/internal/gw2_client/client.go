package gw2client

import (
	"fmt"
	"net/http"
)

const baseUrl = "https://api.guildwars2.com/v2/"

func GetItemsById(ids string) (*http.Response, error) {
	url := baseUrl + "items?ids=" + ids
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return res, err
}

func GetItemIds() (*http.Response, error) {
	url := baseUrl + "items"
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return res, err
}

func GetAllCharacters(apiKey string) (*http.Response, error) {
	url := baseUrl + "characters?ids=all&access_token=" + apiKey
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return res, err
}

func GetTokenInfo(apiKey string) (*http.Response, error) {
	url := baseUrl + "tokeninfo?access_token=" + apiKey
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return res, err
}

func GetAccount(apiKey string) (*http.Response, error) {
	url := baseUrl + "account?access_token=" + apiKey
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return res, err
}

func GetAccountInventory(apiKey string) (*http.Response, error) {
	url := baseUrl + "account/inventory?v=latest&access_token=" + apiKey
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return res, err
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
