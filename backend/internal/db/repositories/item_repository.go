package repositories

import (
	dbmodels "github.com/zoehay/gw2-armory/backend/internal/db/models"
	"gorm.io/gorm"
)

type ItemRepositoryInterface interface {
	Create(item *dbmodels.DBItem) (*dbmodels.DBItem, error)
	CreateMany(items []*dbmodels.DBItem) error
	GetAll() ([]dbmodels.DBItem, error)
	GetFirst() (*dbmodels.DBItem, error)
	GetById(id int) (*dbmodels.DBItem, error)
}

type ItemRepository struct {
	DB *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return ItemRepository{
		DB: db,
	}
}

func (repository *ItemRepository) Create(item *dbmodels.DBItem) (*dbmodels.DBItem, error) {

	err := repository.DB.Create(&item).Error
	if err != nil {
		return nil, err
	}

	return item, nil

}

func (repository *ItemRepository) CreateMany(items []*dbmodels.DBItem) error {

	err := repository.DB.Create(&items).Error
	if err != nil {
		return err
	}

	return nil

}

func (repository *ItemRepository) GetAll() ([]dbmodels.DBItem, error) {
	var items []dbmodels.DBItem

	err := repository.DB.Find(&items).Error
	if err != nil {
		return nil, err
	}

	for i := range items {
		items[i].ToItem()
	}

	return items, nil

}

func (repository *ItemRepository) GetFirst() (*dbmodels.DBItem, error) {
	var item dbmodels.DBItem

	err := repository.DB.First(&dbmodels.DBItem{}).Error
	if err != nil {
		return nil, err
	}

	item.ToItem()
	return &item, nil

}

func (repository *ItemRepository) GetById(id int) (*dbmodels.DBItem, error) {
	var item dbmodels.DBItem

	err := repository.DB.First(&item, id).Error
	if err != nil {
		return nil, err
	}

	item.ToItem()
	return &item, nil

}

func (repository *ItemRepository) GetByIds(ids []int) (*[]dbmodels.DBItem, error) {
	var items []dbmodels.DBItem

	err := repository.DB.Where(ids).Find(&items).Error
	if err != nil {
		return nil, err
	}

	for i := range items {
		items[i].ToItem()
	}

	return &items, nil

}
