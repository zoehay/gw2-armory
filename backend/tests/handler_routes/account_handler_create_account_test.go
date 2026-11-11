package handlerroutes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/api/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type CreateAccountTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
}

func TestCreateAccountTestSuite(t *testing.T) {
	suite.Run(t, new(CreateAccountTestSuite))
}

func (s *CreateAccountTestSuite) SetupSuite() {
	router, repository, _, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Router = router
	s.Repository = repository
}

func (s *CreateAccountTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *CreateAccountTestSuite) TestCreateAccount() {
	var firstSessionID *string

	s.Run("creates account with valid request", func() {
		userJson := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey"}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
		req.Header.Set("Content-Type", "application/json")
		s.Router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusOK, w.Code)

		account, err := testutils.UnmarshalToType[models.Account](w)
		s.Require().NoError(err, "Failed to unmarshal response")
		assert.Equal(s.T(), "gw2apiaccountidstring", account.AccountID)

		dbAccount, err := s.Repository.AccountRepository.GetByID("gw2apiaccountidstring")
		s.Require().NoError(err, "Failed to get account from db")
		assert.Nil(s.T(), dbAccount.Password, "Password is not stored")

		cookies := w.Result().Cookies()
		s.Require().NotEmpty(cookies, "Expected sessionID cookie to be set")
		cookieSessionID := cookies[0].Value
		assert.Equal(s.T(), dbAccount.SessionID, &cookieSessionID, "SessionID in db matches returned cookie")
		firstSessionID = dbAccount.SessionID

		assert.NotContains(s.T(), w.Body.String(), "stringthatisapikey", "API key is not in the response")
		assert.NotContains(s.T(), w.Body.String(), cookieSessionID, "SessionID is not in the response")
	})

	s.Run("returning account with same api key keeps its session", func() {
		userJson := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey"}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
		req.Header.Set("Content-Type", "application/json")
		s.Router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusOK, w.Code)

		account, err := testutils.UnmarshalToType[models.Account](w)
		s.Require().NoError(err, "Failed to unmarshal response")
		assert.Equal(s.T(), "gw2apiaccountidstring", account.AccountID)

		cookies := w.Result().Cookies()
		s.Require().NotEmpty(cookies, "Expected sessionID cookie to be set")
		assert.Equal(s.T(), *firstSessionID, cookies[0].Value, "Existing session is renewed, not replaced")

		session, err := s.Repository.SessionRepository.Get(*firstSessionID)
		s.Require().NoError(err, "Failed to get session from db")
		assert.WithinDuration(s.T(), time.Now().Add(services.SessionDuration), session.Expires, time.Minute, "Renewed session lasts as long as the cookie")
	})

	s.Run("returning account with new api key stores the new key", func() {
		userJson := `{"AccountName":"Name forAccount", "APIKey":"rotatedapikey"}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
		req.Header.Set("Content-Type", "application/json")
		s.Router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusOK, w.Code)

		dbAccount, err := s.Repository.AccountRepository.GetByID("gw2apiaccountidstring")
		s.Require().NoError(err, "Failed to get account from db")
		s.Require().NotNil(dbAccount.APIKey)
		assert.Equal(s.T(), "rotatedapikey", *dbAccount.APIKey)
		assert.Equal(s.T(), firstSessionID, dbAccount.SessionID, "Updating the key keeps the session")
	})
}

func (s *CreateAccountTestSuite) TestCreateAccountWithMalformedBody() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader("not valid json"))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *CreateAccountTestSuite) TestCreateAccountStoresPermissions() {
	allPermissions := []string{"tradingpost", "characters", "wvw", "pvp", "progression", "wallet", "guilds", "builds", "account", "inventories", "unlocks"}

	postKey := func() *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(`{"APIKey":"stringthatisapikey"}`))
		req.Header.Set("Content-Type", "application/json")
		s.Router.ServeHTTP(w, req)
		s.Require().Equal(http.StatusOK, w.Code)
		return w
	}

	s.Run("submit stores and returns the key's permissions", func() {
		w := postKey()

		account, err := testutils.UnmarshalToType[models.Account](w)
		s.Require().NoError(err, "Failed to unmarshal response")
		assert.ElementsMatch(s.T(), allPermissions, account.Permissions)

		dbAccount, err := s.Repository.AccountRepository.GetByID("gw2apiaccountidstring")
		s.Require().NoError(err, "Failed to get account from db")
		assert.ElementsMatch(s.T(), allPermissions, []string(dbAccount.Permissions))

		infoW := httptest.NewRecorder()
		infoReq, _ := http.NewRequest("GET", "/account/info", nil)
		infoReq.AddCookie(w.Result().Cookies()[0])
		s.Router.ServeHTTP(infoW, infoReq)
		s.Require().Equal(http.StatusOK, infoW.Code)

		info, err := testutils.UnmarshalToType[models.Account](infoW)
		s.Require().NoError(err, "Failed to unmarshal account info")
		assert.ElementsMatch(s.T(), allPermissions, info.Permissions)
	})

	s.Run("returning account without stored permissions gets them on next submit", func() {
		err := s.Repository.AccountRepository.DB.Exec("UPDATE db_accounts SET permissions = NULL WHERE account_id = ?", "gw2apiaccountidstring").Error
		s.Require().NoError(err, "Failed to clear permissions")

		postKey()

		dbAccount, err := s.Repository.AccountRepository.GetByID("gw2apiaccountidstring")
		s.Require().NoError(err, "Failed to get account from db")
		assert.ElementsMatch(s.T(), allPermissions, []string(dbAccount.Permissions))
	})
}
