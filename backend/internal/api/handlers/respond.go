package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

// respondError logs the detailed error and sends the client only a generic message.
// Wrapped errors can carry internals like GW2 request URLs, which include the API key.
func respondError(c *gin.Context, status int, message string, err error) {
	log.Printf("%s %s: %s: %v", c.Request.Method, c.FullPath(), message, err)
	c.IndentedJSON(status, gin.H{"error": message})
}
