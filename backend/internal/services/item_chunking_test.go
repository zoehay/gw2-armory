package services

import (
	"errors"
	"slices"
	"testing"

	apimodels "github.com/zoehay/gw2-armory/backend/internal/api/models"
)

type recordingItemService struct {
	calls [][]int
	err   error
}

func (r *recordingItemService) FetchAndStoreItemsByID(ids []int) error {
	r.calls = append(r.calls, slices.Clone(ids))
	return r.err
}
func (r *recordingItemService) FetchAndStoreAllItems() error             { return nil }
func (r *recordingItemService) GetAllItems() ([]apimodels.Item, error)   { return nil, nil }
func (r *recordingItemService) GetItemByID(int) (*apimodels.Item, error) { return nil, nil }

func sequentialIDs(n int) []int64 {
	ids := make([]int64, n)
	for i := range ids {
		ids[i] = int64(i + 1)
	}
	return ids
}

func TestSplitArray(t *testing.T) {
	tests := []struct {
		name      string
		length    int
		chunkSize int
		want      []int
	}{
		{"empty", 0, 200, nil},
		{"under one chunk", 3, 200, []int{3}},
		{"exactly one chunk", 200, 200, []int{200}},
		{"one over", 201, 200, []int{200, 1}},
		{"several chunks", 450, 200, []int{200, 200, 50}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]int, tt.length)
			for i := range arr {
				arr[i] = i
			}

			chunks := splitArray(arr, tt.chunkSize)

			var sizes []int
			var joined []int
			for _, chunk := range chunks {
				sizes = append(sizes, len(chunk))
				joined = append(joined, chunk...)
			}
			if !slices.Equal(sizes, tt.want) {
				t.Errorf("chunk sizes = %v, want %v", sizes, tt.want)
			}
			if len(arr) > 0 && !slices.Equal(joined, arr) {
				t.Error("chunks do not cover the input in order")
			}
		})
	}
}

func TestFetchMissingItemsUsesFullChunks(t *testing.T) {
	items := &recordingItemService{}
	service := &BagItemService{ItemService: items}

	service.FetchMissingItems(sequentialIDs(450))

	var sizes []int
	for _, call := range items.calls {
		sizes = append(sizes, len(call))
	}
	if want := []int{200, 200, 50}; !slices.Equal(sizes, want) {
		t.Errorf("request sizes = %v, want %v", sizes, want)
	}
}

func TestFetchMissingItemsDeduplicates(t *testing.T) {
	items := &recordingItemService{}
	service := &BagItemService{ItemService: items}

	service.FetchMissingItems([]int64{7, 3, 7, 3, 9})

	if len(items.calls) != 1 || !slices.Equal(items.calls[0], []int{7, 3, 9}) {
		t.Errorf("calls = %v, want one call with [7 3 9]", items.calls)
	}
}

func TestFetchMissingItemsNoIDsMakesNoRequest(t *testing.T) {
	items := &recordingItemService{}
	service := &BagItemService{ItemService: items}

	service.FetchMissingItems(nil)

	if len(items.calls) != 0 {
		t.Errorf("calls = %v, want none", items.calls)
	}
}

func TestFetchMissingItemsContinuesAfterChunkError(t *testing.T) {
	items := &recordingItemService{err: errors.New("gw2 unavailable")}
	service := &BagItemService{ItemService: items}

	service.FetchMissingItems(sequentialIDs(401))

	if len(items.calls) != 3 {
		t.Errorf("made %d requests, want 3 even when chunks fail", len(items.calls))
	}
}
