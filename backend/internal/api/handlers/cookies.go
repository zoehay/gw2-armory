package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const sessionCookieName = "sessionID"

// setSessionCookie sets or, with a negative maxAge, clears the session cookie.
// Setting and clearing share attributes because browsers only clear a cookie whose path and domain match.
// No Domain makes it host-only, so it isn't sent to subdomains
func (h AccountHandler) setSessionCookie(c *gin.Context, value string, maxAge int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		Secure:   h.SecureCookies,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
