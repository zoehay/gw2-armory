package handlerroutes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type ItemHandlerTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Service    *services.Service
}

func TestItemHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ItemHandlerTestSuite))
}

func (s *ItemHandlerTestSuite) SetupSuite() {
	router, repository, service, err := testutils.DBRouterSetup()
	if err != nil {
		s.T().Errorf("Error setting up router: %v", err)
	}

	s.Router = router
	s.Repository = repository
	s.Service = service

	// err = s.Service.ItemService.GetAndStoreAllItems()
	// if err != nil {
	// 	s.T().Errorf("Error getting and storing items: %v", err)
	// }

}

func (s *ItemHandlerTestSuite) TearDownSuite() {
	dropTables := []string{"db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	if err != nil {
		s.T().Errorf("Error tearing down suite: %v", err)
	}
}

func (s *ItemHandlerTestSuite) TestGetItemById() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/27952", nil)
	s.Router.ServeHTTP(w, req)
	assert.Equal(s.T(), 200, w.Code)
}
