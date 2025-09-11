package dbmodels

import (
	"github.com/lib/pq"
	"github.com/zoehay/gw2-armory/backend/internal/api/models"
)

// MVP do not care about which bag an item is in
type DBBagItem struct {
	AccountID     string
	CharacterName *string
	BagItemID     uint
	Count         uint
	Charges       *uint
	Infusions     *pq.Int64Array `gorm:"type:integer[]"`
	Upgrades      *pq.Int64Array `gorm:"type:integer[]"`
	Skin          *uint
	Stats         *models.DetailsMap `gorm:"type:json"`
	Dyes          *pq.Int64Array     `gorm:"type:integer[]"`
	Binding       *string
	BoundTo       *string
	Slot          *string
	Location      *string
}

// func GORMBagItemToBagItem(gormBagItem GORMBagItem) models.BagItem {

// }
