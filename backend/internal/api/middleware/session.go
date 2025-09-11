package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
)

func UseSession(accountRepository *repositories.AccountRepository, sessionRepository *repositories.SessionRepository) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {

		sessionID, err := c.Cookie("sessionID")
		// fmt.Printf("Cookie value: %s \n", sessionID)

		if err != nil {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error getting session cookie": err.Error()})
			c.Abort()
			return
		}

		if sessionID == "" {
			c.IndentedJSON(http.StatusForbidden, gin.H{"session error": err.Error()})
			c.Abort()
			return
		}

		dbSession, err := sessionRepository.Get(sessionID)
		if err != nil {
			c.IndentedJSON(http.StatusForbidden, gin.H{"sessionID not in database": err.Error()})
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
			if err != nil {
				c.IndentedJSON(http.StatusForbidden, gin.H{"error": "no account associated with session"})
				c.Abort()
				return
			}

			c.Set("accountID", account.AccountID)
			c.Set("accountName", account.AccountName)
			c.Set("apiKey", account.APIKey)
			c.Set("sessionID", sessionID)
			c.Next()
		}

		// TODO Add token
		// token := context.GetHeader("Authorization")
		// err := ValidateToken(token)
		// if err != nil {
		// 	context.JSON(401, gin.H{"error": err.Error()})
		// 	context.Abort()
		// 	return
		// }

	})
}
