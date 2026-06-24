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

const guestAccountID = "gw2apiaccountidstring"

type CreateGuestAccountTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
}

func TestCreateGuestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(CreateGuestAccountTestSuite))
}

func (s *CreateGuestAccountTestSuite) SetupSuite() {
	router, repository, _, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Router = router
	s.Repository = repository
}

func (s *CreateGuestAccountTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *CreateGuestAccountTestSuite) TestCreateGuestWithNewAPIKey() {
	userJson := `{"APIKey":"stringthatisapikey"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	account, err := testutils.UnmarshalToType[models.Account](w)
	s.Require().NoError(err, "Failed to unmarshal response")
	assert.Equal(s.T(), guestAccountID, account.AccountID)

	dbAccount, err := s.Repository.AccountRepository.GetByID(guestAccountID)
	s.Require().NoError(err, "Failed to get account from db")
	assert.Equal(s.T(), dbAccount.SessionID, account.SessionID, "SessionID in db matches returned account")

	cookies := w.Result().Cookies()
	s.Require().NotEmpty(cookies, "Expected sessionID cookie to be set")
	cookieSessionID := cookies[0].Value
	assert.Equal(s.T(), dbAccount.SessionID, &cookieSessionID, "SessionID in db matches returned cookie")
}

func (s *CreateGuestAccountTestSuite) TestReturningGuestSessionRenewal() {
	// First call creates the account (or finds the existing one from a prior test)
	userJson := `{"APIKey":"stringthatisapikey"}`
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	req1.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w1, req1)
	s.Require().Equal(http.StatusOK, w1.Code)

	// Second call exercises the returning-guest path (existing account, no password)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	req2.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w2, req2)

	assert.Equal(s.T(), http.StatusOK, w2.Code)

	account, err := testutils.UnmarshalToType[models.Account](w2)
	s.Require().NoError(err, "Failed to unmarshal response")
	assert.Equal(s.T(), guestAccountID, account.AccountID)
	assert.NotNil(s.T(), account.SessionID, "Renewed session should be present in response")

	cookies := w2.Result().Cookies()
	s.Require().NotEmpty(cookies, "Expected sessionID cookie on session renewal")

	// The renewed session should grant access to protected routes
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/account/info", nil)
	req3.AddCookie(cookies[0])
	s.Router.ServeHTTP(w3, req3)
	assert.Equal(s.T(), http.StatusOK, w3.Code)
}

func (s *CreateGuestAccountTestSuite) TestCreateGuestWithMalformedBody() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader("not valid json"))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}
