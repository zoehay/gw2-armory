package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	dbmodels "github.com/zoehay/gw2armoury/backend/internal/db/models"
	"github.com/zoehay/gw2armoury/backend/internal/db/repositories"
	gw2models "github.com/zoehay/gw2armoury/backend/internal/gw2_client/models"
	"github.com/zoehay/gw2armoury/backend/internal/gw2_client/providers"

	"gorm.io/gorm"
)

type AccountServiceInterface interface {
	GetAccount(apiKey string) (*gw2models.GW2Account, error)
	GetTokenInfo(apiKey string) (*gw2models.GW2TokenInfo, error)
	GenerateOrUpdateAccount(requestAccount *dbmodels.DBAccount, gw2AccountID string) (*dbmodels.DBAccount, *dbmodels.DBSession, error)
	RenewOrGenerateSession(account *dbmodels.DBAccount) (*dbmodels.DBAccount, *dbmodels.DBSession, error)
	generateNewSession(account *dbmodels.DBAccount) (updatedAccount *dbmodels.DBAccount, newSession *dbmodels.DBSession, err error)
	generateSessionID() (sessionID string, err error)
	IsRecrawlDue(lastCrawl *time.Time) bool
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

func (service *AccountService) GetAccount(apiKey string) (*gw2models.GW2Account, error) {
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

	return account, nil
}

func (service *AccountService) GetTokenInfo(apiKey string) (*gw2models.GW2TokenInfo, error) {
	token, err := service.AccountProvider.GetTokenInfo(apiKey)
	if err != nil {
		return nil, fmt.Errorf("service error using provider could not get account id: %s", err)
	}
	if token.ID == nil || token.Name == nil {
		return nil, fmt.Errorf("service error no token id or name: %s", err)
	}

	return token, nil
}

func (service *AccountService) GenerateOrUpdateAccount(requestAccount *dbmodels.DBAccount, gw2AccountID string) (*dbmodels.DBAccount, *dbmodels.DBSession, error) {
	var account *dbmodels.DBAccount

	existingAccount, err := service.AccountRepository.GetByID(gw2AccountID)
	if err != nil {
		// new user
		if errors.Is(err, gorm.ErrRecordNotFound) {
			account, err = service.AccountRepository.Create(requestAccount)
			if err != nil {
				return nil, nil, fmt.Errorf("account repository create error: %s", err)
			}
		} else {
			// db error
			if err != nil {
				return nil, nil, fmt.Errorf("error accessing account db: %s", err)
			}
		}
	} else {
		// returning user
		// TODO replace password with user
		if existingAccount.Password != nil {
			// existing full account
			return nil, nil, fmt.Errorf("error existing account for account id: %s", gw2AccountID)
		} else {
			// returning user has not previously set a password
			if requestAccount.Password != nil {
				// existing guest account, accountRequest has password so upgrade to full account
				account, err = service.AccountRepository.Update(existingAccount, requestAccount)
				// TODO add password encryption
				if err != nil {
					return nil, nil, fmt.Errorf("account repository update account error: %s", err)
				}
			} else {
				// existing guest account, no password in request so update api key
				// account, err = handler.AccountRepository.UpdateAPIKey(existingAccount.AccountID, *newAccount.APIKey)
				// if err != nil {
				// 	return nil, fmt.Errorf("account repository update apikey error: %s", err)
				// }
				account = existingAccount
			}
		}
	}

	account, session, err := service.RenewOrGenerateSession(account)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating or updating session: %s", err.Error())
	}

	return account, session, nil
}

func (service *AccountService) RenewOrGenerateSession(account *dbmodels.DBAccount) (*dbmodels.DBAccount, *dbmodels.DBSession, error) {
	var session *dbmodels.DBSession
	var err error

	if account.SessionID != nil {
		session, err = service.SessionRepository.Renew(*account.SessionID)
		if err != nil {
			return nil, nil, fmt.Errorf("error renewing session for existing account: %w", err)
		}
	} else {
		account, session, err = service.generateNewSession(account)
		if err != nil {
			return nil, nil, fmt.Errorf("error generating new session for existing account: %w", err)
		}
	}

	return account, session, nil
}

func (service *AccountService) generateNewSession(account *dbmodels.DBAccount) (updatedAccount *dbmodels.DBAccount, newSession *dbmodels.DBSession, err error) {
	newSessionID, err := service.generateSessionID()
	if err != nil {
		return nil, nil, err
	}

	var session = &dbmodels.DBSession{
		SessionID: newSessionID,
		Expires:   time.Now().Add(3600 * time.Second),
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
