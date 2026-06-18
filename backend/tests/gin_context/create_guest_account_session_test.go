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
	"github.com/zoehay/gw2-armory/backend/internal/api/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type CreateGuestAccountSessionTestSuite struct {
	suite.Suite
	Router         *gin.Engine
	Repository     *repositories.Repository
	Service        *services.Service
	AccountHandler *handlers.AccountHandler
}

func TestCreateGuestAccountSessionSuite(t *testing.T) {
	suite.Run(t, new(CreateGuestAccountSessionTestSuite))
}

func (s *CreateGuestAccountSessionTestSuite) SetupSuite() {
	router, repository, service, err := testutils.DBRouterSetup()
	if err != nil {
		s.T().Errorf("Error setting up router: %v", err)
	}

	s.Router = router
	s.Repository = repository
	s.Service = service
	s.AccountHandler = handlers.NewAccountHandler("localhost", service.AccountService, service.BagItemService)

}

func (s *CreateGuestAccountSessionTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	if err != nil {
		s.T().Errorf("Error tearing down suite: %v", err)
	}
}

func (s *CreateGuestAccountSessionTestSuite) TestCreateGuestWithNewAPIKey() {

	userJson := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey", "Password":"stringthatispassword"}`
	req, _ := http.NewRequest("POST", "/addkey", strings.NewReader(userJson))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = req
	s.AccountHandler.HandlePostAPIKeyRequest(c)

	cookie := w.Result().Cookies()[0]

	account, err := testutils.UnmarshalToType[models.Account](w)
	if err != nil {
		s.T().Errorf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(s.T(), 200, w.Code)
	assert.Equal(s.T(), "sessionID", cookie.Name, "Correct cookie name")
	assert.Equal(s.T(), *account.SessionID, cookie.Value)
	assert.Equal(s.T(), "gw2apiaccountidstring", account.AccountID)
	assert.Equal(s.T(), "gw2name", *account.GW2AccountName)
	assert.Equal(s.T(), "armourytest", *account.GW2TokenName)
}

// func (s *CreateGuestAccountSessionTestSuite) TestOldAPIKeyRefreshesSession() {}
