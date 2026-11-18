package services

import (
	"slices"
	"testing"
)

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
