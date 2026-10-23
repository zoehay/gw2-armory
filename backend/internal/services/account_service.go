package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	apimodels "github.com/zoehay/gw2-armory/backend/internal/api/models"
	dbmodels "github.com/zoehay/gw2-armory/backend/internal/db/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/gw2_client/providers"

	"gorm.io/gorm"
)

// SessionDuration is how long a new or renewed session, and its cookie, stays valid
const SessionDuration = time.Hour

type AccountServiceInterface interface {
	FetchAccount(apiKey string) (*apimodels.Account, error)
	FetchToken(apiKey string) (*apimodels.Token, error)
	GetAccountByID(accountID string) (*apimodels.Account, error)
	GenerateOrUpdateAccount(requestAccount *apimodels.Account, gw2AccountID string) (*apimodels.Account, *apimodels.Session, error)
	RenewOrGenerateSession(account *dbmodels.DBAccount) (*apimodels.Account, *apimodels.Session, error)
	generateNewSession(account *dbmodels.DBAccount) (updatedAccount *dbmodels.DBAccount, newSession *dbmodels.DBSession, err error)
	generateSessionID() (sessionID string, err error)
	IsRecrawlDue(lastCrawl *time.Time) bool
	UpdateLastCrawl(accountID string) error
	Logout(sessionID string) error
	DeleteAccount(accountID string, sessionID string) error
}

type AccountService struct {
	AccountRepository *repositories.AccountRepository
	AccountProvider   providers.AccountDataProvider
	SessionRepository *repositories.SessionRepository
}

func NewAccountService(accountRepository *repositories.AccountRepository, accountProvider providers.AccountDataProvider, sessionRepository *repositories.SessionRepository) *AccountService {
	return &AccountService{
		AccountRepository: accountRepository,
		AccountProvider:   accountProvider,
		SessionRepository: sessionRepository,
	}
}

func (service *AccountService) FetchAccount(apiKey string) (*apimodels.Account, error) {
	account, err := service.AccountProvider.GetAccount(apiKey)
	if err != nil {
		return nil, fmt.Errorf("service error using provider could not get account id: %s", err)
	}
	if account.ID == nil {
		return nil, fmt.Errorf("service error no account id: %s", err)
	}
	if account.Name == nil {
		return nil, fmt.Errorf("service error no account id: %s", err)
	}

	result := account.ToAccount()
	return &result, nil
}

func (service *AccountService) FetchToken(apiKey string) (*apimodels.Token, error) {
	token, err := service.AccountProvider.GetTokenInfo(apiKey)
	if err != nil {
		return nil, fmt.Errorf("service error using provider could not get account id: %s", err)
	}
	if token.ID == nil || token.Name == nil {
		return nil, fmt.Errorf("service error no token id or name: %s", err)
	}

	result := token.ToToken()
	return &result, nil
}

func (service *AccountService) GetAccountByID(accountID string) (*apimodels.Account, error) {
	account, err := service.AccountRepository.GetByID(accountID)
	if err != nil {
		return nil, err
	}
	result := account.DBAccountToAccount()
	return &result, nil
}

func (service *AccountService) GenerateOrUpdateAccount(requestAccount *apimodels.Account, gw2AccountID string) (*apimodels.Account, *apimodels.Session, error) {
	dbRequestAccount := &dbmodels.DBAccount{
		AccountID:      requestAccount.AccountID,
		AccountName:    requestAccount.AccountName,
		GW2AccountName: requestAccount.GW2AccountName,
		GW2TokenName:   requestAccount.GW2TokenName,
		APIKey:         requestAccount.APIKey,
	}

	var account *dbmodels.DBAccount

	existingAccount, err := service.AccountRepository.GetByID(gw2AccountID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("error accessing account db: %s", err)
		}
		// new account
		account, err = service.AccountRepository.Create(dbRequestAccount)
		if err != nil {
			return nil, nil, fmt.Errorf("account repository create error: %s", err)
		}
	} else {
		// returning account, store the request key in case it was rotated
		account = existingAccount
		if requestAccount.APIKey != nil && (account.APIKey == nil || *account.APIKey != *requestAccount.APIKey) {
			account, err = service.AccountRepository.UpdateAPIKey(gw2AccountID, *requestAccount.APIKey)
			if err != nil {
				return nil, nil, fmt.Errorf("account repository update api key error: %s", err)
			}
		}
	}

	apiAccount, apiSession, err := service.RenewOrGenerateSession(account)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating or updating session: %s", err.Error())
	}

	return apiAccount, apiSession, nil
}

func (service *AccountService) RenewOrGenerateSession(account *dbmodels.DBAccount) (*apimodels.Account, *apimodels.Session, error) {
	var session *dbmodels.DBSession
	var err error

	if account.SessionID != nil {
		session, err = service.SessionRepository.Renew(*account.SessionID, time.Now().Add(SessionDuration))
		if err != nil {
			return nil, nil, fmt.Errorf("error renewing session for existing account: %w", err)
		}
	} else {
		account, session, err = service.generateNewSession(account)
		if err != nil {
			return nil, nil, fmt.Errorf("error generating new session for existing account: %w", err)
		}
	}

	apiAccount := account.DBAccountToAccount()
	return &apiAccount, (*apimodels.Session)(session), nil
}

func (service *AccountService) generateNewSession(account *dbmodels.DBAccount) (updatedAccount *dbmodels.DBAccount, newSession *dbmodels.DBSession, err error) {
	newSessionID, err := service.generateSessionID()
	if err != nil {
		return nil, nil, err
	}

	var session = &dbmodels.DBSession{
		SessionID: newSessionID,
		Expires:   time.Now().Add(SessionDuration),
	}

	newSession, err = service.SessionRepository.Create(session)
	if err != nil {
		return nil, nil, err
	}

	updatedAccount, err = service.AccountRepository.UpdateSession(account.AccountID, newSession)
	if err != nil {
		return nil, nil, err
	}

	return updatedAccount, newSession, nil
}

func (service *AccountService) generateSessionID() (sessionID string, err error) {
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return "", err
	}
	sessionID = base64.RawURLEncoding.EncodeToString(b)
	return sessionID, nil
}

func (service *AccountService) IsRecrawlDue(lastCrawl *time.Time) bool {
	minHoursSinceCrawl := float64(1)
	var elapsed float64

	if lastCrawl != nil {
		t := time.Now()
		elapsed = t.Sub(*lastCrawl).Hours()
	}

	return (elapsed >= minHoursSinceCrawl || lastCrawl == nil)
}

func (service *AccountService) UpdateLastCrawl(accountID string) error {
	return service.AccountRepository.UpdateLastCrawl(accountID)
}

func (service *AccountService) Logout(sessionID string) error {
	return service.SessionRepository.Delete(sessionID)
}

func (service *AccountService) DeleteAccount(accountID string, sessionID string) error {
	tx := service.AccountRepository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Exec(`DELETE FROM db_bag_item_infusions WHERE db_bag_item_id IN (SELECT id FROM db_bag_items WHERE account_id = ?)`, accountID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting bag item infusions: %w", err)
	}
	if err := tx.Exec(`DELETE FROM db_bag_item_upgrades WHERE db_bag_item_id IN (SELECT id FROM db_bag_items WHERE account_id = ?)`, accountID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting bag item upgrades: %w", err)
	}
	if err := tx.Where("account_id = ?", accountID).Delete(&dbmodels.DBBagItem{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting bag items: %w", err)
	}
	if err := tx.Where("account_id = ?", accountID).Delete(&dbmodels.DBAccount{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting account: %w", err)
	}
	if err := tx.Where("session_id = ?", sessionID).Delete(&dbmodels.DBSession{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting session: %w", err)
	}

	return tx.Commit().Error
}
