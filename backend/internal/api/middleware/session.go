package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"gorm.io/gorm"
)

func UseSession(accountRepository *repositories.AccountRepository, sessionRepository *repositories.SessionRepository) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {

		sessionID, err := c.Cookie("sessionID")

		if err != nil {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error": "missing session cookie"})
			c.Abort()
			return
		}

		if sessionID == "" {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error": "empty session cookie"})
			c.Abort()
			return
		}

		dbSession, err := sessionRepository.Get(sessionID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error": "session not found"})
			c.Abort()
			return
		}
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error looking up session"})
			c.Abort()
			return
		}

		now := time.Now()
		isExpired := dbSession.Expires.Before(now)

		if isExpired {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error": "session is expired"})
			c.Abort()
			return
		} else {
			account, err := accountRepository.GetBySession(sessionID)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.IndentedJSON(http.StatusForbidden, gin.H{"error": "no account associated with session"})
				c.Abort()
				return
			}
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error looking up account"})
				c.Abort()
				return
			}

			c.Set("accountID", account.AccountID)
			c.Set("accountName", account.AccountName)
			c.Set("apiKey", account.APIKey)
			c.Set("sessionID", sessionID)
			c.Next()
		}

	})
}
