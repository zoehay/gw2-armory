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
		assert.Equal(s.T(), dbAccount.SessionID, account.SessionID, "SessionID in db matches returned account")
		assert.Nil(s.T(), dbAccount.Password, "Password is not stored")

		cookies := w.Result().Cookies()
		s.Require().NotEmpty(cookies, "Expected sessionID cookie to be set")
		cookieSessionID := cookies[0].Value
		assert.Equal(s.T(), dbAccount.SessionID, &cookieSessionID, "SessionID in db matches returned cookie")
		firstSessionID = dbAccount.SessionID
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
		assert.Equal(s.T(), firstSessionID, account.SessionID, "Existing session is renewed, not replaced")
	})
}

func (s *CreateAccountTestSuite) TestCreateAccountWithMalformedBody() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader("not valid json"))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}
