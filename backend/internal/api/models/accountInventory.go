package models

type AccountInventory struct {
	AccountID       string       `json:"id"`
	SharedInventory *[]BagItem   `json:"shared_inventory,omitempty"`
	Characters      *[]Character `json:"characters,omitempty"`
}
