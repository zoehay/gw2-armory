package dbmodels

import (
	"github.com/lib/pq"
	"github.com/zoehay/gw2armoury/backend/internal/api/models"
)

type DBItem struct {
	ID           uint `gorm:"primaryKey"`
	ChatLink     string
	Name         string
	Icon         *string
	Description  *string
	Type         string
	Rarity       string
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
		Type:         &dbItem.Name,
		Rarity:       &dbItem.Rarity,
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

// type GormArmorDetails struct {
// 	Type        string
// 	WeightClass string
// 	Defense     int
// 	// InfusionSlots       *[]InfusionSlotType `json:"infusion_slots"`
// 	AttributeAdjustment float32 `sql:"type:decimal(10,2);"`
// 	// InfixUpgrade        *ApiInfixUpgrade `json:"infix_upgrade"`
// 	SuffixItemId          *int
// 	SecondarySuffixItemId string
// 	// StatChoices *[]int `json:"stat_choices"`
// }

// type GameType string

// type Flag string

// const (
// 	AccountBindOnUse
//     AccountBound
//     Attuned
//     BulkConsume
//     DeleteWarning
//     HideSuffix
//     Infused
//     MonsterOnly
//     NoMysticForge
//     NoSalvage
//     NoSell
//     NotUpgradeable
//     NoUnderwater
//     SoulbindOnAcquire
//     SoulBindOnUse
//     Tonic
//     Unique
// )

// type ArmorDetails struct {
// 	Type: string
// 	WeightClass: string
// 	Defense: int
// 	InfusionSlots: InfusionSlot
// 	AttributeAdjustment: int
// 	InfixUpgrade: InfixUpgrade
// 	SuffixItemId: int
// 	SecondarySuffixItemId: string
// 	StatChoices: []int
// }

// func GormItemToItem(gormItem GormItem) models.Item {
// 	return models.Item{
// 		ID:           gormItem.ID,
// 		ChatLink:     gormItem.Icon,
// 		Name:         gormItem.Icon,
// 		Icon:         gormItem.Icon,
// 		Description:  gormItem.Icon,
// 		Type:         gormItem.Icon,
// 		Rarity:       gormItem.Icon,
// 		Level:        gormItem.Level,
// 		VendorValue:  gormItem.VendorValue,
// 		DefaultSkin:  *gormItem.DefaultSkin,
// 		Flags:        gormItem.Flags,
// 		GameTypes:    gormItem.GameTypes,
// 		Restrictions: gormItem.Restrictions,
// 		UpgradesInto: *gormItem.UpgradesInto,
// 		UpgradesFrom: *gormItem.UpgradesFrom,
// 		// Details : gormItem.Details,
// 	}
// }

// func ItemToGormItem(item models.Item) GormItem {
// 	return GormItem{
// 		ID:           item.ID,
// 		ChatLink:     item.Icon,
// 		Name:         item.Icon,
// 		Icon:         item.Icon,
// 		Description:  item.Icon,
// 		Type:         item.Icon,
// 		Rarity:       item.Icon,
// 		Level:        item.Level,
// 		VendorValue:  item.VendorValue,
// 		DefaultSkin:  &item.DefaultSkin,
// 		Flags:        item.Flags,
// 		GameTypes:    item.GameTypes,
// 		Restrictions: item.Restrictions,
// 		UpgradesInto: &item.UpgradesInto,
// 		UpgradesFrom: &item.UpgradesFrom,
// 		// Details : item.Details,
// 	}
// }
