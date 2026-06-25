package repositories

import (
	"fmt"

	dbmodels "github.com/zoehay/gw2-armory/backend/internal/db/models"
	"gorm.io/gorm"
)

type BagItemRepositoryInterface interface {
	Create(BagItem *dbmodels.DBBagItem) (*dbmodels.DBBagItem, error)
	DeleteByAccountID(accountID string) error
	DeleteByCharacterName(accountID string, characterName string) error
	DeleteSharedInventory(accountID string) error
	ReplaceCharacterInventory(accountID string, characterName string, items []dbmodels.DBBagItem) error
	ReplaceSharedInventory(accountID string, items []dbmodels.DBBagItem) error
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
	if err := repository.createItem(repository.DB, bagItem); err != nil {
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
	return repository.deleteCharacterInventory(repository.DB, accountID, characterName)
}

func (repository *BagItemRepository) DeleteSharedInventory(accountID string) error {
	return repository.deleteSharedInventory(repository.DB, accountID)
}

func (repository *BagItemRepository) ReplaceCharacterInventory(accountID string, characterName string, items []dbmodels.DBBagItem) error {
	tx := repository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	if err := repository.deleteCharacterInventory(tx, accountID, characterName); err != nil {
		tx.Rollback()
		return err
	}
	for i := range items {
		if err := repository.createItem(tx, &items[i]); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (repository *BagItemRepository) ReplaceSharedInventory(accountID string, items []dbmodels.DBBagItem) error {
	tx := repository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	if err := repository.deleteSharedInventory(tx, accountID); err != nil {
		tx.Rollback()
		return err
	}
	for i := range items {
		if err := repository.createItem(tx, &items[i]); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
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
	likePattern := fmt.Sprintf("%%%s%%", searchTerm)
	var bagItems []dbmodels.DBBagItem
	err := repository.DB.
		Preload("Item").
		Preload("Infusions").
		Preload("Upgrades").
		Select("db_bag_items.*").
		Joins("JOIN db_items ON db_bag_items.bag_item_id = db_items.id").
		Where("db_bag_items.account_id = ? AND (db_items.name ILIKE ? OR db_items.description ILIKE ? OR db_items.rarity ILIKE ?)",
			accountID, likePattern, likePattern, likePattern).
		Find(&bagItems).Error
	if err != nil {
		return nil, err
	}
	return bagItems, nil
}

func (repository *BagItemRepository) createItem(db *gorm.DB, bagItem *dbmodels.DBBagItem) error {
	for i := range bagItem.Infusions {
		if err := db.FirstOrCreate(&bagItem.Infusions[i], dbmodels.DBItem{ID: bagItem.Infusions[i].ID}).Error; err != nil {
			return err
		}
	}
	for i := range bagItem.Upgrades {
		if err := db.FirstOrCreate(&bagItem.Upgrades[i], dbmodels.DBItem{ID: bagItem.Upgrades[i].ID}).Error; err != nil {
			return err
		}
	}
	return db.Omit("Infusions.*", "Upgrades.*").Create(bagItem).Error
}

func (repository *BagItemRepository) deleteCharacterInventory(db *gorm.DB, accountID string, characterName string) error {
	if err := db.Exec(`DELETE FROM db_bag_item_infusions WHERE db_bag_item_id IN (SELECT id FROM db_bag_items WHERE account_id = ? AND character_name = ?)`, accountID, characterName).Error; err != nil {
		return err
	}
	if err := db.Exec(`DELETE FROM db_bag_item_upgrades WHERE db_bag_item_id IN (SELECT id FROM db_bag_items WHERE account_id = ? AND character_name = ?)`, accountID, characterName).Error; err != nil {
		return err
	}
	return db.Where("account_id = ? AND character_name = ?", accountID, characterName).Delete(&dbmodels.DBBagItem{}).Error
}

func (repository *BagItemRepository) deleteSharedInventory(db *gorm.DB, accountID string) error {
	if err := db.Exec(`DELETE FROM db_bag_item_infusions WHERE db_bag_item_id IN (SELECT id FROM db_bag_items WHERE account_id = ? AND source = 'shared')`, accountID).Error; err != nil {
		return err
	}
	if err := db.Exec(`DELETE FROM db_bag_item_upgrades WHERE db_bag_item_id IN (SELECT id FROM db_bag_items WHERE account_id = ? AND source = 'shared')`, accountID).Error; err != nil {
		return err
	}
	return db.Where("account_id = ? AND source = ?", accountID, "shared").Delete(&dbmodels.DBBagItem{}).Error
}
