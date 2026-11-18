package providers

import (
	"net/http"
	"testing"
)

func TestGetItemsByIDsReturnsItems(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusPartialContent} {
		serveGW2Response(t, status, `[{"id": 24, "name": "Sealed Package of Snowballs"}]`)

		items, err := (&ItemProvider{}).GetItemsByIDs([]int{24, 1})
		if err != nil {
			t.Fatalf("status %d: unexpected error: %v", status, err)
		}
		if len(items) != 1 || items[0].ID != 24 {
			t.Errorf("status %d: items = %v, want item 24", status, items)
		}
	}
}

func TestGetItemsByIDsAllUnknownIsEmpty(t *testing.T) {
	serveGW2Response(t, http.StatusNotFound, `{"text": "all ids provided are invalid"}`)

	items, err := (&ItemProvider{}).GetItemsByIDs([]int{1, 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("items = %v, want none", items)
	}
}

func TestGetItemsByIDsFailureIsAnError(t *testing.T) {
	for _, status := range []int{http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusServiceUnavailable} {
		serveGW2Response(t, status, `{"text": "unavailable"}`)

		if _, err := (&ItemProvider{}).GetItemsByIDs([]int{24}); err == nil {
			t.Errorf("status %d: expected an error", status)
		}
	}
}
