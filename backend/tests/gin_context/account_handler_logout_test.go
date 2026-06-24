package gincontext_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/api/handlers"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type LogoutTestSuite struct {
	suite.Suite
	Repository     *repositories.Repository
	AccountHandler *handlers.AccountHandler
}

func TestLogoutTestSuite(t *testing.T) {
	suite.Run(t, new(LogoutTestSuite))
}

func (s *LogoutTestSuite) SetupSuite() {
	_, repository, service, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Repository = repository
	s.AccountHandler = handlers.NewAccountHandler("localhost", service.AccountService, service.BagItemService)
}

func (s *LogoutTestSuite) SetupTest() {
	err := s.Repository.AccountRepository.DB.Exec("TRUNCATE TABLE db_bag_item_infusions, db_bag_item_upgrades, db_bag_items, db_sessions, db_accounts CASCADE").Error
	s.Require().NoError(err, "Error truncating tables")
}

func (s *LogoutTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *LogoutTestSuite) TestLogout() {
	// Create an account to establish a session
	userJson := `{"APIKey":"stringthatisapikey"}`
	req1, _ := http.NewRequest("POST", "/addkey", strings.NewReader(userJson))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request = req1
	s.AccountHandler.HandlePostAPIKeyRequest(c1)
	s.Require().Equal(http.StatusOK, w1.Code, "Account creation must succeed")

	cookies := w1.Result().Cookies()
	s.Require().NotEmpty(cookies, "Expected sessionID cookie after account creation")
	sessionCookie := cookies[0]
	sessionID := sessionCookie.Value

	// Logout with the session cookie
	req2, _ := http.NewRequest("POST", "/logout", nil)
	req2.AddCookie(sessionCookie)
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = req2
	s.AccountHandler.Logout(c2)
	c2.Writer.WriteHeaderNow()

	assert.Equal(s.T(), http.StatusNoContent, w2.Code)

	// Cookie should be cleared in the response
	responseCookies := w2.Result().Cookies()
	s.Require().NotEmpty(responseCookies, "Expected cleared sessionID cookie in logout response")
	clearedCookie := responseCookies[0]
	assert.Equal(s.T(), "sessionID", clearedCookie.Name)
	assert.Empty(s.T(), clearedCookie.Value, "Cookie value should be cleared")

	// Session should be deleted from DB
	_, err := s.Repository.SessionRepository.Get(sessionID)
	assert.Error(s.T(), err, "Session should not exist after logout")
}

func (s *LogoutTestSuite) TestLogoutWithNoSession() {
	req, _ := http.NewRequest("POST", "/logout", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	s.AccountHandler.Logout(c)
	c.Writer.WriteHeaderNow()

	assert.Equal(s.T(), http.StatusNoContent, w.Code)
}
