package dbmodels

import (
	"github.com/lib/pq"
	"github.com/zoehay/gw2-armory/backend/internal/api/models"
)

type DBBagItem struct {
	ID            uint               `gorm:"primaryKey;autoIncrement"`
	AccountID     string
	CharacterName *string
	Source        string
	BagItemID     uint
	Item          DBItem             `gorm:"foreignKey:BagItemID"`
	Count         uint
	Charges       *uint
	Infusions     []DBItem           `gorm:"many2many:db_bag_item_infusions;"`
	Upgrades      []DBItem           `gorm:"many2many:db_bag_item_upgrades;"`
	Skin          *uint
	Stats         *models.DetailsMap `gorm:"type:json"`
	Dyes          *pq.Int64Array     `gorm:"type:integer[]"`
	Binding       *string
	BoundTo       *string
	Slot          *string
	Location      *string
}

func (b DBBagItem) ToBagItem() models.BagItem {
	characterName := ""
	if b.CharacterName != nil {
		characterName = *b.CharacterName
	}

	var infusionIDs *[]int64
	var infusionDetails *[]map[string]interface{}
	if len(b.Infusions) > 0 {
		ids := make([]int64, len(b.Infusions))
		details := make([]map[string]interface{}, len(b.Infusions))
		for i, inf := range b.Infusions {
			ids[i] = int64(inf.ID)
			var name interface{}
			if inf.Name != "" {
				name = inf.Name
			}
			var rarity interface{}
			if inf.Rarity != "" {
				rarity = inf.Rarity
			}
			details[i] = map[string]interface{}{
				"id":      inf.ID,
				"name":    name,
				"icon":    inf.Icon,
				"rarity":  rarity,
				"details": (*map[string]interface{})(inf.Details),
			}
		}
		infusionIDs = &ids
		infusionDetails = &details
	}

	var upgradeIDs *[]int64
	var upgradeDetails *[]map[string]interface{}
	if len(b.Upgrades) > 0 {
		ids := make([]int64, len(b.Upgrades))
		details := make([]map[string]interface{}, len(b.Upgrades))
		for i, upg := range b.Upgrades {
			ids[i] = int64(upg.ID)
			var name interface{}
			if upg.Name != "" {
				name = upg.Name
			}
			var rarity interface{}
			if upg.Rarity != "" {
				rarity = upg.Rarity
			}
			details[i] = map[string]interface{}{
				"id":      upg.ID,
				"name":    name,
				"icon":    upg.Icon,
				"rarity":  rarity,
				"details": (*map[string]interface{})(upg.Details),
			}
		}
		upgradeIDs = &ids
		upgradeDetails = &details
	}

	bagItem := models.BagItem{
		CharacterName:   characterName,
		Source:          b.Source,
		BagItemID:       b.BagItemID,
		Count:           b.Count,
		Charges:         b.Charges,
		Infusions:       infusionIDs,
		Upgrades:        upgradeIDs,
		Skin:            b.Skin,
		Stats:           (*map[string]interface{})(b.Stats),
		Dyes:            (*[]int64)(b.Dyes),
		Binding:         b.Binding,
		BoundTo:         b.BoundTo,
		Slot:            b.Slot,
		Location:        b.Location,
		InfusionDetails: infusionDetails,
		UpgradeDetails:  upgradeDetails,
	}

	if b.Item.Name != "" {
		bagItem.Name = &b.Item.Name
		bagItem.Icon = b.Item.Icon
		bagItem.Description = b.Item.Description
		bagItem.Type = &b.Item.Type
		bagItem.Rarity = &b.Item.Rarity
		bagItem.VendorValue = &b.Item.VendorValue
		bagItem.Details = (*map[string]interface{})(b.Item.Details)
	}

	return bagItem
}

func DBBagItemsToAccountInventory(bagItems []DBBagItem, accountID string) (accountInventory models.AccountInventory, itemsNotInDB []int64) {
	characterNameMap := map[string]models.Character{}
	var sharedInventory []models.BagItem
	var characters []models.Character

	for _, dbItem := range bagItems {
		if dbItem.Item.Name == "" {
			itemsNotInDB = append(itemsNotInDB, int64(dbItem.BagItemID))
		}
		for _, infusion := range dbItem.Infusions {
			if infusion.Name == "" {
				itemsNotInDB = append(itemsNotInDB, int64(infusion.ID))
			}
		}
		for _, upgrade := range dbItem.Upgrades {
			if upgrade.Name == "" {
				itemsNotInDB = append(itemsNotInDB, int64(upgrade.ID))
			}
		}

		item := dbItem.ToBagItem()
		name := item.CharacterName

		if item.Source == "shared" {
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
