package dbmodels

import "github.com/zoehay/gw2-armory/backend/internal/api/models"

type DBCharacter struct {
	Name       string
	Race       string
	Gender     string
	Flags      []string
	Profession string
	Level      int
	Guild      *string
	Age        int
	Created    string
	Deaths     int
	// Crafting   []*DetailsMap
	Title *int
	// Backstory []string
	// WvwAbilities []string
	// BuildTabsUnlocked     int
	// ActiveBuildTab        int
	// BuildTabs             []*DetailsMap
	// EquipmentTabsUnlocked int
	// ActiveEquipmentTab    int
	Equipment []*models.DetailsMap
	// EquipmentTabs []*DetailsMap
	// Recipes []int
	// Training []string
	// Bags []GormBag
}
