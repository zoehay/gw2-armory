package repositories

import (
	dbmodels "github.com/zoehay/gw2-armory/backend/internal/db/models"
	"gorm.io/gorm"
)

type BagItemRepositoryInterface interface {
	Create(BagItem *dbmodels.DBBagItem) (*dbmodels.DBBagItem, error)
	DeleteByAccountID(accountID string) error
	DeleteByCharacterName(accountID string, characterName string) error
	GetByCharacterName(accountID string, characterName string) ([]dbmodels.DBBagItem, error)
	DeleteSharedInventory(accountID string) error
	GetIds() ([]int, error)
	GetDetailBagItemByCharacterName(accountID string, characterName string) ([]dbmodels.DBBagItem, error)
	GetDetailBagItemByAccountID(accountID string) ([]dbmodels.DBBagItem, error)
	GetDetailBagItemsWithSearch(accountID string, searchTerm string) ([]dbmodels.DBBagItem, error)
}

type BagItemRepository struct {
	DB *gorm.DB
}

func NewBagItemRepository(db *gorm.DB) BagItemRepository {
	return BagItemRepository{
		DB: db,
	}
}

func (repository *BagItemRepository) Create(bagItem *dbmodels.DBBagItem) (*dbmodels.DBBagItem, error) {
	for i := range bagItem.Infusions {
		if err := repository.DB.FirstOrCreate(&bagItem.Infusions[i], dbmodels.DBItem{ID: bagItem.Infusions[i].ID}).Error; err != nil {
			return nil, err
		}
	}
	for i := range bagItem.Upgrades {
		if err := repository.DB.FirstOrCreate(&bagItem.Upgrades[i], dbmodels.DBItem{ID: bagItem.Upgrades[i].ID}).Error; err != nil {
			return nil, err
		}
	}
	err := repository.DB.Omit("Infusions.*", "Upgrades.*").Create(bagItem).Error
	if err != nil {
		return nil, err
	}
	return bagItem, nil
}

func (repository *BagItemRepository) DeleteByAccountID(accountID string) error {
	repository.DB.Exec(`DELETE FROM db_bag_item_infusions WHERE db_bag_item_id IN (SELECT id FROM db_bag_items WHERE account_id = ?)`, accountID)
	repository.DB.Exec(`DELETE FROM db_bag_item_upgrades WHERE db_bag_item_id IN (SELECT id FROM db_bag_items WHERE account_id = ?)`, accountID)
	return repository.DB.Where("account_id = ?", accountID).Delete(&dbmodels.DBBagItem{}).Error
}

func (repository *BagItemRepository) DeleteByCharacterName(accountID string, characterName string) error {
	repository.DB.Exec(`DELETE FROM db_bag_item_infusions WHERE db_bag_item_id IN (SELECT id FROM db_bag_items WHERE account_id = ? AND character_name = ?)`, accountID, characterName)
	repository.DB.Exec(`DELETE FROM db_bag_item_upgrades WHERE db_bag_item_id IN (SELECT id FROM db_bag_items WHERE account_id = ? AND character_name = ?)`, accountID, characterName)
	return repository.DB.Where("account_id = ? AND character_name = ?", accountID, characterName).Delete(&dbmodels.DBBagItem{}).Error
}

func (repository *BagItemRepository) GetByCharacterName(accountID string, characterName string) ([]dbmodels.DBBagItem, error) {
	var bagItems []dbmodels.DBBagItem
	err := repository.DB.Where("account_id = ? AND character_name = ?", accountID, characterName).Find(&bagItems).Error
	if err != nil {
		return nil, err
	}
	return bagItems, nil
}

func (repository *BagItemRepository) DeleteSharedInventory(accountID string) error {
	repository.DB.Exec(`DELETE FROM db_bag_item_infusions WHERE db_bag_item_id IN (SELECT id FROM db_bag_items WHERE account_id = ? AND character_name = 'Shared Inventory')`, accountID)
	repository.DB.Exec(`DELETE FROM db_bag_item_upgrades WHERE db_bag_item_id IN (SELECT id FROM db_bag_items WHERE account_id = ? AND character_name = 'Shared Inventory')`, accountID)
	return repository.DB.Where("account_id = ? AND character_name = ?", accountID, "Shared Inventory").Delete(&dbmodels.DBBagItem{}).Error
}

func (repository *BagItemRepository) GetIds() ([]int, error) {
	var bagItemIds []int
	err := repository.DB.Model(&dbmodels.DBBagItem{}).Pluck("bag_item_id", &bagItemIds).Error
	if err != nil {
		return nil, err
	}
	return bagItemIds, nil
}

func (repository *BagItemRepository) GetDetailBagItemByCharacterName(accountID string, characterName string) ([]dbmodels.DBBagItem, error) {
	var bagItems []dbmodels.DBBagItem
	err := repository.DB.
		Preload("Item").
		Preload("Infusions").
		Preload("Upgrades").
		Where("account_id = ? AND character_name = ?", accountID, characterName).
		Find(&bagItems).Error
	if err != nil {
		return nil, err
	}
	return bagItems, nil
}

func (repository *BagItemRepository) GetDetailBagItemByAccountID(accountID string) ([]dbmodels.DBBagItem, error) {
	var bagItems []dbmodels.DBBagItem
	err := repository.DB.
		Preload("Item").
		Preload("Infusions").
		Preload("Upgrades").
		Where("account_id = ?", accountID).
		Find(&bagItems).Error
	if err != nil {
		return nil, err
	}
	return bagItems, nil
}

func (repository *BagItemRepository) GetDetailBagItemsWithSearch(accountID string, searchTerm string) ([]dbmodels.DBBagItem, error) {
	var bagItems []dbmodels.DBBagItem
	err := repository.DB.
		Preload("Item").
		Preload("Infusions").
		Preload("Upgrades").
		Select("db_bag_items.*").
		Joins("JOIN db_items ON db_bag_items.bag_item_id = db_items.id").
		Where("db_bag_items.account_id = ? AND (db_items.name ILIKE ? OR db_items.description ILIKE ? OR db_items.rarity ILIKE ?)",
			accountID, searchTerm, searchTerm, searchTerm).
		Find(&bagItems).Error
	if err != nil {
		return nil, err
	}
	return bagItems, nil
}
