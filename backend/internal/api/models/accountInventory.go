package models

type AccountInventory struct {
	AccountID          string       `json:"id"`
	SharedInventory    *[]BagItem   `json:"shared_inventory,omitempty"`
	BankInventory      *[]BagItem   `json:"bank_inventory,omitempty"`
	MaterialsInventory *[]BagItem   `json:"materials_inventory,omitempty"`
	Characters         *[]Character `json:"characters,omitempty"`
}
