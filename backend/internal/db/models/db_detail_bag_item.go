package dbmodels

import (
	"github.com/lib/pq"
	"github.com/zoehay/gw2armoury/backend/internal/api/models"
)

type DBDetailBagItem struct {
	AccountID     string             `json:"account_id"`
	CharacterName string             `json:"character_name"`
	BagItemID     uint               `json:"id"`
	Count         uint               `json:"count"`
	Charges       *uint              `json:"charges,omitempty"`
	Infusions     *pq.Int64Array     `json:"infusions,omitempty" gorm:"type:integer[]"`
	Upgrades      *pq.Int64Array     `json:"upgrades,omitempty" gorm:"type:integer[]"`
	Skin          *uint              `json:"skin,omitempty"`
	Stats         *models.DetailsMap `json:"stats,omitempty" gorm:"type:json"`
	Dyes          *pq.Int64Array     `json:"dyes,omitempty" gorm:"type:integer[]"`
	Binding       *string            `json:"binding,omitempty"`
	BoundTo       *string            `json:"bound_to,omitempty"`
	Slot          *string            `json:"slot"`
	Location      *string            `json:"location"`

	// fields from full item details optional in case not in db
	Name        *string            `json:"name"`
	Icon        *string            `json:"icon"`
	Description *string            `json:"description"`
	Type        *string            `json:"type"`
	Rarity      *string            `json:"rarity"`
	VendorValue *uint              `json:"vendor_value"`
	Details     *models.DetailsMap `json:"details" gorm:"type:json"`
}

func (dbDetailBagItem DBDetailBagItem) ToBagItem() models.BagItem {
	return models.BagItem{
		CharacterName: dbDetailBagItem.CharacterName,
		BagItemID:     dbDetailBagItem.BagItemID,
		Count:         dbDetailBagItem.Count,
		Charges:       dbDetailBagItem.Charges,
		Infusions:     (*[]int64)(dbDetailBagItem.Infusions),
		Upgrades:      (*[]int64)(dbDetailBagItem.Upgrades),
		Skin:          dbDetailBagItem.Skin,
		Stats:         (*map[string]interface{})(dbDetailBagItem.Stats),
		Dyes:          (*[]int64)(dbDetailBagItem.Dyes),
		Binding:       dbDetailBagItem.Binding,
		BoundTo:       dbDetailBagItem.BoundTo,
		Slot:          dbDetailBagItem.Slot,
		Location:      dbDetailBagItem.Location,

		Name:        dbDetailBagItem.Name,
		Icon:        dbDetailBagItem.Icon,
		Description: dbDetailBagItem.Description,
		Type:        dbDetailBagItem.Type,
		Rarity:      dbDetailBagItem.Rarity,
		VendorValue: dbDetailBagItem.VendorValue,
		Details:     (*map[string]interface{})(dbDetailBagItem.Details),
	}
}

func DBDetailBagItemsToAccountInventory(dbIconBagItems []DBDetailBagItem, accountID string) (accountInventory models.AccountInventory, itemsNotInDB []int64) {

	characterNameMap := map[string]models.Character{}
	var sharedInventory []models.BagItem
	var characters []models.Character

	for _, item := range dbIconBagItems {
		if ItemNotInDB(item) {
			itemsNotInDB = append(itemsNotInDB, int64(item.BagItemID))
		}

		item := item.ToBagItem()
		name := item.CharacterName

		if name == "Shared Inventory" {
			sharedInventory = append(sharedInventory, item)
		} else {
			entry, ok := characterNameMap[name]
			isEquipment := item.IsEquipment()

			if ok {
				if isEquipment {
					entry.Equipment = append(entry.Equipment, item)
					characterNameMap[name] = entry
				} else {
					entry.Inventory = append(entry.Inventory, item)
					characterNameMap[name] = entry
				}
			} else {
				newCharacter := &models.Character{
					Name:      name,
					Equipment: []models.BagItem{},
					Inventory: []models.BagItem{},
				}
				if isEquipment {
					newCharacter.Equipment = append(newCharacter.Equipment, item)
				} else {
					newCharacter.Inventory = append(newCharacter.Inventory, item)
				}
				characterNameMap[name] = *newCharacter
			}
		}

	}

	for character := range characterNameMap {
		characters = append(characters, characterNameMap[character])
	}

	accountInventory.AccountID = accountID
	accountInventory.SharedInventory = &sharedInventory
	accountInventory.Characters = &characters

	return accountInventory, itemsNotInDB
}

func ItemNotInDB(item DBDetailBagItem) bool {
	if (item.Name) == nil {
		return true
	} else {
		return false
	}
}
