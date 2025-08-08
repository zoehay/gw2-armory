package models

type Item struct {
	ID           uint                      `json:"id"`
	ChatLink     *string                   `json:"chat_link"`
	Name         *string                   `json:"name"`
	Icon         *string                   `json:"icon"`
	Description  *string                   `json:"description"`
	Type         *string                   `json:"type"`
	Rarity       *string                   `json:"rarity"`
	Level        *uint                     `json:"level"`
	VendorValue  *uint                     `json:"vendor_value"`
	DefaultSkin  *uint                     `json:"default_skin,omitempty"`
	Flags        *[]string                 `json:"flags" gorm:"type:text[]"`
	GameTypes    *[]string                 `json:"game_types" gorm:"type:text[]"`
	Restrictions *[]string                 `json:"restrictions" gorm:"type:text[]"`
	UpgradesInto *[]map[string]interface{} `json:"upgrades_into,omitempty" gorm:"type:json"`
	UpgradesFrom *[]map[string]interface{} `json:"upgrades_from,omitempty" gorm:"type:json"`
	Details      *map[string]interface{}   `json:"details,omitempty" gorm:"type:json"`
}
