package models

type Item struct {
	ID           uint                      `json:"id"`
	ChatLink     *string                   `json:"chat_link,omitempty"`
	Name         *string                   `json:"name,omitempty"`
	Icon         *string                   `json:"icon,omitempty"`
	Description  *string                   `json:"description,omitempty"`
	Type         *string                   `json:"type,omitempty"`
	Rarity       *string                   `json:"rarity,omitempty"`
	Level        *uint                     `json:"level,omitempty"`
	VendorValue  *uint                     `json:"vendor_value,omitempty"`
	DefaultSkin  *uint                     `json:"default_skin,omitempty"`
	Flags        *[]string                 `json:"flags,omitempty" gorm:"type:text[]"`
	GameTypes    *[]string                 `json:"game_types,omitempty" gorm:"type:text[]"`
	Restrictions *[]string                 `json:"restrictions,omitempty" gorm:"type:text[]"`
	UpgradesInto *[]map[string]interface{} `json:"upgrades_into,omitempty" gorm:"type:json"`
	UpgradesFrom *[]map[string]interface{} `json:"upgrades_from,omitempty" gorm:"type:json"`
	Details      *map[string]interface{}   `json:"details,omitempty" gorm:"type:json"`
}
