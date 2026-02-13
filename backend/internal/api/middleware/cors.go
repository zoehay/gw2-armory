package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOriginsString := os.Getenv("CORS_ALLOW_ORIGIN")
		allowedOrigins := strings.Split(allowedOriginsString, " ")

		originIsAllowed := func(origin string, allowedOrigins []string) bool {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		}

		origin := c.Request.Header.Get("Origin")
		if originIsAllowed(origin, allowedOrigins) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
		}

		c.Next()
	}
}
