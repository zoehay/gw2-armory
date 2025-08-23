package gw2models

import (
	"github.com/lib/pq"
	"github.com/zoehay/gw2armoury/backend/internal/api/models"
	dbmodels "github.com/zoehay/gw2armoury/backend/internal/db/models"
)

type ItemRequest struct {
	IDs string `json:"ids"`
}

type ItemResponse struct {
	Items []GW2Item
}

type GW2Item struct {
	ID           uint                      `json:"id"`
	ChatLink     string                    `json:"chat_link"`
	Name         string                    `json:"name"`
	Icon         *string                   `json:"icon,omitempty"`
	Description  *string                   `json:"description,omitempty"`
	Type         string                    `json:"type"`
	Rarity       string                    `json:"rarity"`
	Level        uint                      `json:"level"`
	VendorValue  uint                      `json:"vendor_value"`
	DefaultSkin  *uint                     `json:"default_skin,omitempty"`
	Flags        *[]string                 `json:"flags,omitempty"`
	GameTypes    *[]string                 `json:"game_types,omitempty"`
	Restrictions *[]string                 `json:"restrictions,omitempty"`
	UpgradesInto *[]map[string]interface{} `json:"upgrades_into,omitempty"`
	UpgradesFrom *[]map[string]interface{} `json:"upgrades_from,omitempty"`
	Details      *map[string]interface{}   `json:"details,omitempty"`
}

func (gw2Item GW2Item) ToDBItem() dbmodels.DBItem {
	var upgradesInto models.DetailsMapArray
	if gw2Item.UpgradesInto != nil && len(*gw2Item.UpgradesInto) > 0 {
		for _, upgrade := range *gw2Item.UpgradesInto {
			upgradesInto = append(upgradesInto, (models.DetailsMap)(upgrade))
		}
	}

	var upgradesFrom models.DetailsMapArray
	if gw2Item.UpgradesFrom != nil && len(*gw2Item.UpgradesFrom) > 0 {
		for _, upgrade := range *gw2Item.UpgradesFrom {
			upgradesFrom = append(upgradesFrom, (models.DetailsMap)(upgrade))
		}
	}

	return dbmodels.DBItem{
		ID:           gw2Item.ID,
		ChatLink:     gw2Item.ChatLink,
		Name:         gw2Item.Name,
		Icon:         gw2Item.Icon,
		Description:  gw2Item.Description,
		Type:         gw2Item.Type,
		Rarity:       gw2Item.Rarity,
		Level:        gw2Item.Level,
		VendorValue:  gw2Item.VendorValue,
		DefaultSkin:  gw2Item.DefaultSkin,
		Flags:        (*pq.StringArray)(gw2Item.Flags),
		GameTypes:    (*pq.StringArray)(gw2Item.GameTypes),
		Restrictions: (*pq.StringArray)(gw2Item.Restrictions),
		UpgradesInto: &upgradesInto,
		UpgradesFrom: &upgradesFrom,
		Details:      (*models.DetailsMap)(gw2Item.Details),
	}
}

type Upgrade struct {
	Upgrade string `json:"upgrade"`
	ItemID  uint   `json:"item_id"`
}

// func (gw2Item GW2Item) ToItem(gw2Item gw2Item) models.Item {
// 	return models.Item{
// 		ID:           gw2Item.ID,
// 		ChatLink:     gw2Item.ChatLink,
// 		Name:         gw2Item.Name,
// 		Icon:         gw2Item.Icon,
// 		Description:  gw2Item.Description,
// 		Type:         gw2Item.Type,
// 		Rarity:       gw2Item.Rarity,
// 		Level:        gw2Item.Level,
// 		VendorValue:  gw2Item.VendorValue,
// 		DefaultSkin:  *gw2Item.DefaultSkin,
// 		Flags:        gw2Item.Flags,
// 		GameTypes:    gw2Item.GameTypes,
// 		Restrictions: gw2Item.Restrictions,
// 		UpgradesInto: *gw2Item.UpgradesInto,
// 		UpgradesFrom: *gw2Item.UpgradesFrom,
// 		Details:      gw2Item.Details,
// 	}
// }

type ItemError struct {
	Text string `json:"text"`
}
