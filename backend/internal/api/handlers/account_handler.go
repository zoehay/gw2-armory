package handlers

import (
	"errors"
	"log"
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
		respondError(c, http.StatusInternalServerError, "could not get account", err)
		return
	}

	c.IndentedJSON(http.StatusOK, account)
}

func (h AccountHandler) HandlePostAPIKeyRequest(c *gin.Context) {

	var accountRequest AccountRequest

	if err := c.BindJSON(&accountRequest); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request body", err)
		return
	}

	gw2Token, err := h.AccountService.FetchToken(accountRequest.APIKey)
	if errors.Is(err, services.ErrInvalidKey) {
		respondError(c, http.StatusBadRequest, "invalid API key", err)
		return
	}
	if err != nil {
		respondError(c, http.StatusBadGateway, "could not verify API key with the GW2 API", err)
		return
	}

	// checked before fetching the account so a rejected key costs one GW2 call
	if missing := services.MissingPermissions(gw2Token.Permissions); len(missing) > 0 {
		log.Printf("%s %s: rejected API key missing permissions %v", c.Request.Method, c.FullPath(), missing)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":               "API key is missing required permissions",
			"missing_permissions": missing,
		})
		return
	}

	gw2Account, err := h.AccountService.FetchAccount(accountRequest.APIKey)
	if err != nil {
		respondError(c, http.StatusBadGateway, "could not get account from the GW2 API", err)
		return
	}

	// token name is set by user, truncate if very long
	gw2TokenName := string(*gw2Token.Name)
	if len(*gw2Token.Name) > 25 {
		gw2TokenName = gw2TokenName[:25]
	}

	var requestAccount = &apimodels.Account{
		AccountID:      gw2Account.AccountID,
		AccountName:    accountRequest.AccountName,
		GW2AccountName: gw2Account.GW2AccountName,
		GW2TokenName:   &gw2TokenName,
		Permissions:    gw2Token.Permissions,
		APIKey:         &accountRequest.APIKey,
	}

	account, session, err := h.AccountService.GenerateOrUpdateAccount(requestAccount, gw2Account.AccountID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not create or update account", err)
		return
	}

	c.SetCookie("sessionID", session.SessionID, int(services.SessionDuration.Seconds()), "/", h.Domain, false, true)

	if h.AccountService.IsRecrawlDue(account.LastCrawl) {
		err = h.BagItemService.FetchAndStoreAllBagItems(account.AccountID, accountRequest.APIKey)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "could not get inventory", err)
			return
		}
		err = h.AccountService.UpdateLastCrawl(account.AccountID)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "could not update account", err)
			return
		}
	}

	c.IndentedJSON(http.StatusOK, account)
}

func (h AccountHandler) Delete(c *gin.Context) {
	accountID := c.MustGet("accountID").(string)
	sessionID := c.MustGet("sessionID").(string)

	if err := h.AccountService.DeleteAccount(accountID, sessionID); err != nil {
		respondError(c, http.StatusInternalServerError, "could not delete account", err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"deleted": accountID})
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

type AccountRequest struct {
	AccountName *string
	APIKey      string
}

type APIKeyRequest struct {
	APIKey string
}
