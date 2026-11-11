package services

import (
	"slices"
	"testing"
)

func TestMissingPermissions(t *testing.T) {
	tests := []struct {
		name    string
		granted []string
		want    []string
	}{
		{"all required", []string{"account", "characters", "inventories", "wallet"}, nil},
		{"extra permissions ignored", []string{"wallet", "pvp", "inventories", "account", "unlocks", "characters"}, nil},
		{"some missing", []string{"account", "inventories"}, []string{"characters", "wallet"}},
		{"none granted", nil, []string{"account", "characters", "inventories", "wallet"}},
		{"unknown permissions only", []string{"pvp", "wvw"}, []string{"account", "characters", "inventories", "wallet"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MissingPermissions(tt.granted)
			if !slices.Equal(got, tt.want) {
				t.Errorf("MissingPermissions(%v) = %v, want %v", tt.granted, got, tt.want)
			}
		})
	}
}
