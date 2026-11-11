package handlerroutes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type SessionCookieTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
}

func TestSessionCookieTestSuite(t *testing.T) {
	suite.Run(t, new(SessionCookieTestSuite))
}

func (s *SessionCookieTestSuite) SetupSuite() {
	router, repository, _, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Router = router
	s.Repository = repository
}

func (s *SessionCookieTestSuite) SetupTest() {
	err := s.Repository.AccountRepository.DB.Exec("TRUNCATE TABLE db_bag_item_infusions, db_bag_item_upgrades, db_bag_items, db_sessions, db_accounts CASCADE").Error
	s.Require().NoError(err, "Error truncating tables")
}

func (s *SessionCookieTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *SessionCookieTestSuite) login() *http.Cookie {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(`{"APIKey":"stringthatisapikey"}`))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)
	s.Require().Equal(http.StatusOK, w.Code)
	return s.sessionCookie(w)
}

func (s *SessionCookieTestSuite) sessionCookie(w *httptest.ResponseRecorder) *http.Cookie {
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "sessionID" {
			return cookie
		}
	}
	s.FailNow("expected a sessionID cookie in the response")
	return nil
}

// assertCookieFlags checks the attributes shared by setting and clearing the session cookie
func (s *SessionCookieTestSuite) assertCookieFlags(cookie *http.Cookie) {
	assert.True(s.T(), cookie.Secure, "Secure")
	assert.True(s.T(), cookie.HttpOnly, "HttpOnly")
	assert.Equal(s.T(), http.SameSiteLaxMode, cookie.SameSite, "SameSite=Lax")
	assert.Equal(s.T(), "/", cookie.Path, "Path")
	assert.Empty(s.T(), cookie.Domain, "No Domain, so the cookie is host-only")
}

func (s *SessionCookieTestSuite) TestLoginSetsSessionCookieFlags() {
	cookie := s.login()

	s.assertCookieFlags(cookie)
	assert.NotEmpty(s.T(), cookie.Value)
	assert.Equal(s.T(), int(services.SessionDuration.Seconds()), cookie.MaxAge, "Cookie lasts as long as the session")
}

func (s *SessionCookieTestSuite) assertClearsCookie(w *httptest.ResponseRecorder) {
	cookie := s.sessionCookie(w)
	s.assertCookieFlags(cookie)
	assert.Empty(s.T(), cookie.Value)
	assert.Less(s.T(), cookie.MaxAge, 0, "Negative MaxAge tells the browser to delete the cookie")
}

func (s *SessionCookieTestSuite) accountInfoStatus(cookie *http.Cookie) int {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/account/info", nil)
	req.AddCookie(cookie)
	s.Router.ServeHTTP(w, req)
	return w.Code
}

func (s *SessionCookieTestSuite) TestLogoutClearsCookieAndSession() {
	cookie := s.login()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/account/logout", nil)
	req.AddCookie(cookie)
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusNoContent, w.Code)
	s.assertClearsCookie(w)
	assert.Equal(s.T(), http.StatusForbidden, s.accountInfoStatus(cookie), "Old session no longer works")
}

func (s *SessionCookieTestSuite) TestLogoutWithoutSessionStillClearsCookie() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/account/logout", nil)
	req.AddCookie(&http.Cookie{Name: "sessionID", Value: "unknownsession"})
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusNoContent, w.Code)
	s.assertClearsCookie(w)
}

func (s *SessionCookieTestSuite) TestDeleteClearsCookie() {
	cookie := s.login()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/account/delete", nil)
	req.AddCookie(cookie)
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	s.assertClearsCookie(w)
}
