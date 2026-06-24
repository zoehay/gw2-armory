package handlerroutes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/api/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type AccountAuthTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Cookie     *http.Cookie
}

func TestAccountAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AccountAuthTestSuite))
}

func (s *AccountAuthTestSuite) SetupSuite() {
	router, repository, _, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Router = router
	s.Repository = repository

	userJson := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)
	s.Require().Equal(http.StatusOK, w.Code, "Setup: POST /apikeys must succeed")

	cookies := w.Result().Cookies()
	s.Require().NotEmpty(cookies, "Setup: expected sessionID cookie from POST /apikeys")
	s.Cookie = cookies[0]
}

func (s *AccountAuthTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *AccountAuthTestSuite) TestGetAccountInfo() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/account/info", nil)
	req.AddCookie(s.Cookie)
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	account, err := testutils.UnmarshalToType[models.Account](w)
	s.Require().NoError(err, "Failed to unmarshal response")
	assert.Equal(s.T(), "gw2apiaccountidstring", account.AccountID)
}

func (s *AccountAuthTestSuite) TestLoginAndAccessProtectedRoute() {
	// Login returns session_id in body — the handler does not set a cookie
	loginJson := `{"AccountName":"Name forAccount", "Password":"anypassword"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(loginJson))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	account, err := testutils.UnmarshalToType[models.Account](w)
	s.Require().NoError(err, "Failed to unmarshal login response")
	s.Require().NotNil(account.SessionID, "Login must return a session_id in the response body")

	// Use the session_id from the body to access a protected route
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/account/info", nil)
	req2.AddCookie(&http.Cookie{Name: "sessionID", Value: *account.SessionID})
	s.Router.ServeHTTP(w2, req2)

	assert.Equal(s.T(), http.StatusOK, w2.Code)

	info, err := testutils.UnmarshalToType[models.Account](w2)
	s.Require().NoError(err, "Failed to unmarshal account info response")
	assert.Equal(s.T(), "gw2apiaccountidstring", info.AccountID)
}

func (s *AccountAuthTestSuite) TestLoginNonExistentAccount() {
	loginJson := `{"AccountName":"doesnotexist", "Password":"anypassword"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(loginJson))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}
