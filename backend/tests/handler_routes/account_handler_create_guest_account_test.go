package handlerroutes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type CreateGuestAccountTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Service    *services.Service
}

func TestCreateGuestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(CreateGuestAccountTestSuite))
}

func (s *CreateGuestAccountTestSuite) SetupSuite() {
	router, repository, service, err := testutils.DBRouterSetup()
	if err != nil {
		s.T().Errorf("Error setting up router: %v", err)
	}
	s.Router = router
	s.Repository = repository
	s.Service = service
}

func (s *CreateGuestAccountTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	if err != nil {
		s.T().Errorf("Error tearing down suite: %v", err)
	}
}

func (s *CreateGuestAccountTestSuite) TestCreateGuestWithNewAPIKey() {
	gin.SetMode(gin.TestMode)

	userJson := `{"APIKey":"stringthatisapikey"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), 200, w.Code)
}
