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

func (s *AccountAuthTestSuite) TestNameLoginRoutesAreNotRegistered() {
	// Accounts are not password protected, so there must be no way to get a session by account name
	for _, path := range []string{"/login", "/signup"} {
		loginJson := `{"AccountName":"Name forAccount", "Password":"anypassword"}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", path, strings.NewReader(loginJson))
		req.Header.Set("Content-Type", "application/json")
		s.Router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusNotFound, w.Code, path)
	}
}
