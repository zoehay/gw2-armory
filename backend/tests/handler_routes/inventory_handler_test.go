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
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type InventoryHandlerTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Service    *services.Service
	Cookie     *http.Cookie
}

func TestInventoryHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryHandlerTestSuite))
}

func (s *InventoryHandlerTestSuite) SetupSuite() {
	router, repository, service, err := testutils.DBRouterSetup()
	if err != nil {
		s.T().Errorf("Error setting up router: %v", err)
	}

	s.Router = router
	s.Repository = repository
	s.Service = service

	s.T().Log("Create account with POST /apikeys")
	userJson := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey", "Password":"stringthatispassword"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	s.Router.ServeHTTP(w, req)
	assert.Equal(s.T(), http.StatusOK, w.Code)
	cookie := w.Result().Cookies()

	s.Cookie = cookie[0]
}

func (s *InventoryHandlerTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	if err != nil {
		s.T().Errorf("Error tearing down suite: %v", err)
	}
}

func (s *InventoryHandlerTestSuite) TestGetCharacterInventory() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/account/accountinventory", nil)
	req.AddCookie(s.Cookie)
	s.Router.ServeHTTP(w, req)

	_, err := testutils.UnmarshalToType[*models.AccountInventory](w)
	if err != nil {
		s.T().Errorf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(s.T(), 200, w.Code)
}
