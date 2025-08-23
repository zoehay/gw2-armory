package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type BagItem struct {
	CharacterName string                  `json:"character_name"`
	BagItemID     uint                    `json:"id"`
	Count         uint                    `json:"count"`
	Charges       *uint                   `json:"charges,omitempty"`
	Infusions     *[]int64                `json:"infusions,omitempty"`
	Upgrades      *[]int64                `json:"upgrades,omitempty"`
	Skin          *uint                   `json:"skin,omitempty"`
	Stats         *map[string]interface{} `json:"stats,omitempty" gorm:"type:json"`
	Dyes          *[]int64                `json:"dyes,omitempty" gorm:"type:integer[]"`
	Binding       *string                 `json:"binding,omitempty"`
	BoundTo       *string                 `json:"bound_to,omitempty"`
	Slot          *string                 `json:"slot,omitempty"`
	Location      *string                 `json:"location,omitempty"`

	// fields from full item details optional in case not in db
	Name        *string                 `json:"name,omitempty"`
	Icon        *string                 `json:"icon,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Type        *string                 `json:"type,omitempty"`
	Rarity      *string                 `json:"rarity,omitempty"`
	VendorValue *uint                   `json:"vendor_value,omitempty"`
	Details     *map[string]interface{} `json:"details,omitempty" gorm:"type:json"`
}

func (item BagItem) IsEquipment() bool {
	if item.Slot != nil && *item.Slot != "" {
		return true
	} else {
		return false
	}
}

type DetailsMap map[string]interface{}

func (detailsMap *DetailsMap) Scan(value interface{}) error {
	if value == nil {
		*detailsMap = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, detailsMap)
	return err

}

func (detailsMap DetailsMap) Value() (driver.Value, error) {
	if len(detailsMap) == 0 {
		return nil, nil
	}
	return json.Marshal(detailsMap)
}

type DetailsMapArray []map[string]interface{}

func (detailsMapArray *DetailsMapArray) Scan(value interface{}) error {
	if value == nil {
		*detailsMapArray = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, detailsMapArray)
	return err

}

func (detailsMapArray DetailsMapArray) Value() (driver.Value, error) {
	if len(detailsMapArray) == 0 {
		return nil, nil
	}
	return json.Marshal(detailsMapArray)
}
