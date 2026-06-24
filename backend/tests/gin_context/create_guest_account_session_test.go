package gincontext_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/api/handlers"
	"github.com/zoehay/gw2-armory/backend/internal/api/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"

	"github.com/gin-gonic/gin"
)

type CreateGuestAccountSessionTestSuite struct {
	suite.Suite
	Repository     *repositories.Repository
	AccountHandler *handlers.AccountHandler
}

func TestCreateGuestAccountSessionSuite(t *testing.T) {
	suite.Run(t, new(CreateGuestAccountSessionTestSuite))
}

func (s *CreateGuestAccountSessionTestSuite) SetupSuite() {
	_, repository, service, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Repository = repository
	s.AccountHandler = handlers.NewAccountHandler("localhost", service.AccountService, service.BagItemService)
}

func (s *CreateGuestAccountSessionTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *CreateGuestAccountSessionTestSuite) TestCreateGuestWithNewAPIKey() {
	userJson := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey", "Password":"stringthatispassword"}`
	req, _ := http.NewRequest("POST", "/addkey", strings.NewReader(userJson))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	s.AccountHandler.HandlePostAPIKeyRequest(c)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	cookies := w.Result().Cookies()
	s.Require().NotEmpty(cookies, "Expected sessionID cookie to be set")
	cookie := cookies[0]

	account, err := testutils.UnmarshalToType[models.Account](w)
	s.Require().NoError(err, "Failed to unmarshal response")

	assert.Equal(s.T(), "sessionID", cookie.Name, "Correct cookie name")
	assert.Equal(s.T(), *account.SessionID, cookie.Value)
	assert.Equal(s.T(), "gw2apiaccountidstring", account.AccountID)
	assert.Equal(s.T(), "gw2name", *account.GW2AccountName)
	assert.Equal(s.T(), "armourytest", *account.GW2TokenName)
}
