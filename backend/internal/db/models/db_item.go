package dbmodels

import (
	"github.com/lib/pq"
	"github.com/zoehay/gw2-armory/backend/internal/api/models"
)

type DBItem struct {
	ID           uint `gorm:"primaryKey"`
	ChatLink     string
	Name         string
	Icon         *string
	Description  *string
	Type         string
	Rarity       string
	Category     *string
	Level        uint
	VendorValue  uint
	DefaultSkin  *uint
	Flags        *pq.StringArray         `gorm:"type:text[]"`
	GameTypes    *pq.StringArray         `gorm:"type:text[]"`
	Restrictions *pq.StringArray         `gorm:"type:text[]"`
	UpgradesInto *models.DetailsMapArray `gorm:"type:json"`
	UpgradesFrom *models.DetailsMapArray `gorm:"type:json"`
	Details      *models.DetailsMap      `gorm:"type:json"`
}

func (dbItem DBItem) ToItem() models.Item {
	var upgradesInto []map[string]interface{}
	if dbItem.UpgradesInto != nil {
		for _, upgrade := range *dbItem.UpgradesInto {
			upgradesInto = append(upgradesInto, (map[string]interface{})(upgrade))
		}
	}

	var upgradesFrom []map[string]interface{}
	if dbItem.UpgradesFrom != nil {
		for _, upgrade := range *dbItem.UpgradesFrom {
			upgradesFrom = append(upgradesFrom, (map[string]interface{})(upgrade))
		}
	}

	return models.Item{
		ID:           dbItem.ID,
		ChatLink:     &dbItem.ChatLink,
		Name:         &dbItem.Name,
		Icon:         dbItem.Icon,
		Description:  dbItem.Description,
		Type:         &dbItem.Type,
		Rarity:       &dbItem.Rarity,
		Category:     dbItem.Category,
		Level:        &dbItem.Level,
		VendorValue:  &dbItem.VendorValue,
		DefaultSkin:  dbItem.DefaultSkin,
		Flags:        (*[]string)(dbItem.Flags),
		GameTypes:    (*[]string)(dbItem.GameTypes),
		Restrictions: (*[]string)(dbItem.Restrictions),
		UpgradesInto: &upgradesInto,
		UpgradesFrom: &upgradesFrom,
		Details:      (*map[string]interface{})(dbItem.Details),
	}
}
