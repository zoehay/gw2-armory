package providers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	gw2client "github.com/zoehay/gw2-armory/backend/internal/gw2_client"
)

func serveTokenInfo(t *testing.T, status int, body string) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	t.Cleanup(gw2client.SetBaseURL(server.URL + "/"))
}

func TestGetTokenInfoReturnsToken(t *testing.T) {
	serveTokenInfo(t, http.StatusOK, `{"id": "key-id", "name": "my key", "permissions": ["account", "wallet"]}`)

	token, err := (&AccountProvider{}).GetTokenInfo("key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token.Name == nil || *token.Name != "my key" {
		t.Errorf("token name = %v, want %q", token.Name, "my key")
	}
	if len(token.Permissions) != 2 {
		t.Errorf("permissions = %v, want 2 entries", token.Permissions)
	}
}

func TestGetTokenInfoInvalidKey(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusForbidden} {
		serveTokenInfo(t, status, `{"text": "Invalid access token"}`)

		_, err := (&AccountProvider{}).GetTokenInfo("bad-key")
		if !errors.Is(err, ErrInvalidKey) {
			t.Errorf("status %d: err = %v, want ErrInvalidKey", status, err)
		}
	}
}

func TestGetTokenInfoServerError(t *testing.T) {
	serveTokenInfo(t, http.StatusInternalServerError, `{"text": "ErrInternal"}`)

	_, err := (&AccountProvider{}).GetTokenInfo("key")
	if err == nil {
		t.Fatal("expected an error")
	}
	if errors.Is(err, ErrInvalidKey) {
		t.Error("a GW2 server error should not be reported as an invalid key")
	}
}
