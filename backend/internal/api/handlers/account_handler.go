package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apimodels "github.com/zoehay/gw2-armory/backend/internal/api/models"
	"github.com/zoehay/gw2-armory/backend/internal/services"
)

type AccountHandler struct {
	Domain         string
	AccountService services.AccountServiceInterface
	BagItemService services.BagItemServiceInterface
}

func NewAccountHandler(domain string, accountService services.AccountServiceInterface, bagItemService services.BagItemServiceInterface) *AccountHandler {
	return &AccountHandler{
		Domain:         domain,
		AccountService: accountService,
		BagItemService: bagItemService,
	}
}

func (h AccountHandler) GetAccount(c *gin.Context) {
	accountID := c.MustGet("accountID").(string)
	account, err := h.AccountService.GetAccountByID(accountID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, account)
}

func (h AccountHandler) HandlePostAPIKeyRequest(c *gin.Context) {

	var accountRequest AccountRequest

	if err := c.BindJSON(&accountRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"request body bind json error": err.Error()})
		return
	}

	gw2Token, err := h.AccountService.FetchToken(accountRequest.APIKey)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{"error could not get token info from gw2 api": err.Error()})
		return
	}

	gw2Account, err := h.AccountService.FetchAccount(accountRequest.APIKey)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{"error could not get account id from gw2 api": err.Error()})
		return
	}

	// token name is set by user, truncate if very long
	gw2TokenName := string(*gw2Token.Name)
	if len(*gw2Token.Name) > 25 {
		gw2TokenName = gw2TokenName[:25]
	}

	var requestAccount = &apimodels.Account{
		AccountID:      *&gw2Account.AccountID,
		AccountName:    accountRequest.AccountName,
		GW2AccountName: gw2Account.GW2AccountName,
		GW2TokenName:   &gw2TokenName,
		APIKey:         &accountRequest.APIKey,
		Password:       accountRequest.Password,
	}

	account, session, err := h.AccountService.GenerateOrUpdateAccount(requestAccount, *&gw2Account.AccountID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error generating or updating account": err.Error()})
		return
	}

	c.SetCookie("sessionID", session.SessionID, 3600, "/", h.Domain, false, true)

	if h.AccountService.IsRecrawlDue(account.LastCrawl) {
		err = h.BagItemService.FetchAndStoreAllBagItems(account.AccountID, accountRequest.APIKey)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error getting inventory after guest creation": err.Error()})
			return
		}
		err = h.AccountService.UpdateLastCrawl(account.AccountID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error updating account last crawl": err.Error()})
			return
		}
	}

	c.IndentedJSON(http.StatusOK, account)
}

func (h AccountHandler) Delete(c *gin.Context) {

	// use request later for User with multiple Accounts
	var deleteKeyRequest DeleteKeyRequest

	if err := c.BindJSON(&deleteKeyRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"request body bind json error": err.Error()})
		return
	}

	accountID := c.MustGet("accountID").(string)
	sessionID := c.MustGet("sessionID").(string)

	if err := h.AccountService.DeleteAccount(accountID, sessionID); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error deleting account": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"API key deleted": deleteKeyRequest.APIKey})
}

func (h AccountHandler) Login(c *gin.Context) {
	var accountLogin AccountLogin

	if err := c.BindJSON(&accountLogin); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	account, _, err := h.AccountService.Login(accountLogin.AccountName, accountLogin.Password)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, account)
}

func (h AccountHandler) Logout(c *gin.Context) {
	sessionID, err := c.Cookie("sessionID")
	if err == nil {
		if err = h.AccountService.Logout(sessionID); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.SetCookie("sessionID", "", -1, "/", h.Domain, false, true)
	c.Status(http.StatusNoContent)
}

type AccountLogin struct {
	AccountName string
	Password    string
}

type AccountRequest struct {
	AccountName *string
	APIKey      string
	Password    *string
}

type CreateRequest struct {
	AccountName string
	APIKey      string
	Password    string
}

type APIKeyRequest struct {
	APIKey string
}

type DeleteKeyRequest struct {
	APIKey string
}
