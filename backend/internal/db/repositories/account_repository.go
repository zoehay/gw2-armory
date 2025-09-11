package repositories

import (
	"time"

	dbmodels "github.com/zoehay/gw2-armory/backend/internal/db/models"
	"gorm.io/gorm"
)

type AccountRepositoryInterface interface {
	GetBySession(sessionID string) (*dbmodels.DBAccount, error)
	GetByID(id string) (*dbmodels.DBAccount, error)
	// GetByUserID(userID string) ([]*dbmodels.DBAccount, error)
	GetByName(name string) (*dbmodels.DBAccount, error)
	Create(account *dbmodels.DBAccount) (*dbmodels.DBAccount, error)
	UpdateSession(accountID string, session *dbmodels.DBSession) (*dbmodels.DBAccount, error)
	UpdateLastCrawl(accountID string) error
	Update(existingAccount *dbmodels.DBAccount, updateAccount *dbmodels.DBAccount) (*dbmodels.DBAccount, error)
	DeleteAPIKey(accountID string) error
	UpdateAPIKey(accountID string, apiKey string) (*dbmodels.DBAccount, error)
	DeleteAccount(accountID string) error
}

type AccountRepository struct {
	DB *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return AccountRepository{
		DB: db,
	}
}

func (repository *AccountRepository) GetBySession(sessionID string) (*dbmodels.DBAccount, error) {
	var account dbmodels.DBAccount

	err := repository.DB.Where("session_id = ?", sessionID).First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil

}

func (repository *AccountRepository) GetByID(id string) (*dbmodels.DBAccount, error) {
	var account dbmodels.DBAccount

	err := repository.DB.Where("account_id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil

}

// func (repository *AccountRepository) GetByUserID(id string) ([]dbmodels.DBAccount, error) {
// 	var accounts []dbmodels.DBAccount

// 	err := repository.DB.Where("user_id = ?", id).Find(&accounts).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return accounts, nil

// }

func (repository *AccountRepository) GetByName(name string) (*dbmodels.DBAccount, error) {
	var account dbmodels.DBAccount

	err := repository.DB.Where("account_name = ?", name).First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil

}

func (repository *AccountRepository) Create(account *dbmodels.DBAccount) (*dbmodels.DBAccount, error) {

	err := repository.DB.Create(&account).Error
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repository *AccountRepository) UpdateSession(accountID string, session *dbmodels.DBSession) (updatedAccount *dbmodels.DBAccount, err error) {
	var account dbmodels.DBAccount

	err = repository.DB.Model(&account).Where("account_id = ?", accountID).Update("session_id", session.SessionID).Error
	if err != nil {
		return nil, err
	}

	err = repository.DB.Where("account_id = ?", accountID).First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (repository *AccountRepository) UpdateLastCrawl(accountID string) error {
	var account dbmodels.DBAccount

	err := repository.DB.Model(&account).Where("account_id = ?", accountID).Update("last_crawl", time.Now()).Error
	if err != nil {
		return err
	}

	return nil
}

func (repository *AccountRepository) Update(existingAccount *dbmodels.DBAccount, updateAccount *dbmodels.DBAccount) (*dbmodels.DBAccount, error) {
	var account dbmodels.DBAccount

	err := repository.DB.Model(&account).Where("account_id = ?", existingAccount.AccountID).Updates(updateAccount).Error
	if err != nil {
		return nil, err
	}

	err = repository.DB.Where("account_id = ?", existingAccount.AccountID).First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (repository *AccountRepository) DeleteAPIKey(accountID string) error {
	var account dbmodels.DBAccount
	err := repository.DB.Model(&account).Where("account_id = ?", accountID).Update("api_key", "").Error

	return err
}

func (repository *AccountRepository) UpdateAPIKey(accountID string, apiKey string) (*dbmodels.DBAccount, error) {
	var account dbmodels.DBAccount
	err := repository.DB.Model(&account).Where("account_id = ?", accountID).Update("api_key", apiKey).Error
	if err != nil {
		return nil, err
	}

	err = repository.DB.Where("account_id = ?", accountID).First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (repository *AccountRepository) DeleteAccount(accountID string) error {
	var account dbmodels.DBAccount
	err := repository.DB.Model(&account).Where("account_id = ?", accountID).Delete(&account).Error

	return err
}
