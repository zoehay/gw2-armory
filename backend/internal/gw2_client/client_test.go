package gw2client

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestKeyedRequestsUseBearerHeader(t *testing.T) {
	const apiKey = "secret-api-key"

	type seen struct {
		url  string
		auth string
	}
	var requests []seen
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, seen{url: r.URL.String(), auth: r.Header.Get("Authorization")})
	}))
	defer server.Close()

	original := baseUrl
	baseUrl = server.URL + "/"
	defer func() { baseUrl = original }()

	calls := map[string]func(string) (*http.Response, error){
		"GetAllCharacters":      GetAllCharacters,
		"GetTokenInfo":          GetTokenInfo,
		"GetAccount":            GetAccount,
		"GetAccountInventory":   GetAccountInventory,
		"GetBankInventory":      GetBankInventory,
		"GetMaterialsInventory": GetMaterialsInventory,
	}

	for name, call := range calls {
		requests = nil
		res, err := call(apiKey)
		if err != nil {
			t.Fatalf("%s: request failed: %v", name, err)
		}
		res.Body.Close()

		if len(requests) != 1 {
			t.Fatalf("%s: expected 1 request, got %d", name, len(requests))
		}
		if got := requests[0].auth; got != "Bearer "+apiKey {
			t.Errorf("%s: Authorization header = %q, want bearer token", name, got)
		}
		if strings.Contains(requests[0].url, apiKey) {
			t.Errorf("%s: API key leaked into URL %q", name, requests[0].url)
		}
	}
}

func TestKeyedRequestErrorsDoNotContainKey(t *testing.T) {
	const apiKey = "secret-api-key"

	original := baseUrl
	// nothing listens here, so the request fails with a *url.Error that includes the URL
	baseUrl = "http://127.0.0.1:1/"
	defer func() { baseUrl = original }()

	_, err := GetAccount(apiKey)
	if err == nil {
		t.Fatal("expected a connection error")
	}
	if strings.Contains(err.Error(), apiKey) {
		t.Errorf("API key leaked into error: %v", err)
	}
}
