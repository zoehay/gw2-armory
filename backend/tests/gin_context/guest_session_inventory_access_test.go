package gincontext_test

import (
	"encoding/json"
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

type GuestSessionInventoryAccessTestSuite struct {
	suite.Suite
	Router         *gin.Engine
	Repository     *repositories.Repository
	Service        *services.Service
	AccountHandler *handlers.AccountHandler
}

func TestGuestSessionInventoryAccessSuite(t *testing.T) {
	suite.Run(t, new(GuestSessionInventoryAccessTestSuite))
}

func (s *GuestSessionInventoryAccessTestSuite) SetupSuite() {
	router, repository, service, err := testutils.DBRouterSetup()
	if err != nil {
		s.T().Errorf("Error setting up router: %v", err)
	}

	s.Router = router
	s.Repository = repository
	s.Service = service
	s.AccountHandler = handlers.NewAccountHandler(&repository.AccountRepository, &repository.SessionRepository, &repository.BagItemRepository, service.AccountService, service.BagItemService)

	s.Service.ItemService.GetAndStoreAllItems()
}

func (s *GuestSessionInventoryAccessTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	if err != nil {
		s.T().Errorf("Error tearing down suite: %v", err)
	}
}

func (s *GuestSessionInventoryAccessTestSuite) TestNoCookieNoInventoryAccess() {
	s.T().Log("GET /account with no cookie")
	w0 := httptest.NewRecorder()
	req0, _ := http.NewRequest("GET", "/account/characters/Roman%20Meows/inventory", nil)
	s.Router.ServeHTTP(w0, req0)
	assert.Equal(s.T(), http.StatusForbidden, w0.Code)
}

func (s *GuestSessionInventoryAccessTestSuite) TestGuestInventoryAccess() {
	s.T().Log("Create account with POST /apikeys")
	userJson := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey", "Password":"stringthatispassword"}`
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	s.Router.ServeHTTP(w1, req1)
	assert.Equal(s.T(), http.StatusOK, w1.Code)
	cookie := w1.Result().Cookies()

	s.T().Log("GET /account with cookie")
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/account/characters/Roman%20Meows/inventory", nil)
	req2.AddCookie(cookie[0])
	s.Router.ServeHTTP(w2, req2)
	assert.Equal(s.T(), http.StatusOK, w2.Code)

	var response []models.BagItem
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	if err != nil {
		s.T().Errorf("Failed to unmarshal response: %v", err)
	}
}
