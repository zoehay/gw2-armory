package repositories

import (
	dbmodels "github.com/zoehay/gw2armoury/backend/internal/db/models"
	"gorm.io/gorm"
)

type BagItemRepositoryInterface interface {
	Create(BagItem *dbmodels.DBBagItem) (*dbmodels.DBBagItem, error)
	DeleteByAccountID(accountID string) error
	DeleteByCharacterName(accountID string, characterName string) error
	GetByCharacterName(accountID string, characterName string) ([]dbmodels.DBBagItem, error)
	DeleteSharedInventory(accountID string) error
	GetIds() ([]int, error)
	GetDetailBagItemByCharacterName(accountID string, characterName string) ([]dbmodels.DBDetailBagItem, error)
	GetDetailBagItemByAccountID(accountID string) (detailBagItems []dbmodels.DBDetailBagItem, err error)
	GetDetailBagItemsWithSearch(accountID string, searchTerm string) (detailBagItems []dbmodels.DBDetailBagItem, err error)
}

type BagItemRepository struct {
	DB *gorm.DB
}

func NewBagItemRepository(db *gorm.DB) BagItemRepository {
	return BagItemRepository{
		DB: db,
	}
}

func (repository *BagItemRepository) Create(BagItem *dbmodels.DBBagItem) (*dbmodels.DBBagItem, error) {
	err := repository.DB.Create(&BagItem).Error
	if err != nil {
		return nil, err
	}

	return BagItem, nil
}

func (repository *BagItemRepository) DeleteByAccountID(accountID string) error {
	err := repository.DB.Where("db_bag_items.account_id = ?", accountID).Delete(&dbmodels.DBBagItem{}).Error

	return err
}

func (repository *BagItemRepository) DeleteByCharacterName(accountID string, characterName string) error {
	err := repository.DB.Where("account_id = ? AND character_name = ?", accountID, characterName).Delete(&dbmodels.DBBagItem{}).Error

	return err
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
	err := repository.DB.Where("db_bag_items.account_id = ? AND character_name = ?", accountID, "Shared Inventory").Delete(&dbmodels.DBBagItem{}).Error

	return err
}

func (repository *BagItemRepository) GetIds() ([]int, error) {
	var bagItemIds []int
	err := repository.DB.Model(&dbmodels.DBBagItem{}).Pluck("bag_item_id", &bagItemIds).Error
	if err != nil {
		return nil, err
	}

	return bagItemIds, nil

}

func (repository *BagItemRepository) GetDetailBagItemByCharacterName(accountID string, characterName string) ([]dbmodels.DBDetailBagItem, error) {
	var detailBagItems []dbmodels.DBDetailBagItem

	err := repository.DB.Table("db_bag_items").
		Select("db_bag_items.*, db_items.icon, db_items.name, db_items.description, db_items.rarity").
		Joins("left join db_items on db_bag_items.bag_item_id = db_items.id").
		Where("db_bag_items.account_id = ? AND db_bag_items.character_name = ?", accountID, characterName).
		Scan(&detailBagItems).Error

	if err != nil {
		return nil, err
	}

	return detailBagItems, nil
}

func (repository *BagItemRepository) GetDetailBagItemByAccountID(accountID string) (detailBagItems []dbmodels.DBDetailBagItem, err error) {
	err = repository.DB.Table("db_bag_items").
		Select("db_bag_items.*, db_items.name, db_items.description, db_items.type, db_items.rarity, db_items.vendor_value, db_items.icon").
		Joins("left join db_items on db_bag_items.bag_item_id = db_items.id").
		Where("db_bag_items.account_id = ?", accountID).
		Scan(&detailBagItems).Error

	if err != nil {
		return nil, err
	}

	return detailBagItems, nil

}

func (repository *BagItemRepository) GetDetailBagItemsWithSearch(accountID string, searchTerm string) (detailBagItems []dbmodels.DBDetailBagItem, err error) {
	err = repository.DB.Table("db_bag_items").
		Select("db_bag_items.*, db_items.name, db_items.description, db_items.type, db_items.rarity, db_items.vendor_value, db_items.icon").
		Joins("left join db_items on db_bag_items.bag_item_id = db_items.id").
		Where("db_bag_items.account_id = ?", accountID).
		Where("name ILIKE ? OR description ILIKE ? OR rarity ILIKE ?", searchTerm, searchTerm, searchTerm).
		Scan(&detailBagItems).Error

	if err != nil {
		return nil, err
	}

	return detailBagItems, nil

}
