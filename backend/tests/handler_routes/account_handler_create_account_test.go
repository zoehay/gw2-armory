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

type CreateAccountTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Service    *services.Service
}

func TestCreateAccountTestSuite(t *testing.T) {
	suite.Run(t, new(CreateAccountTestSuite))
}

func (s *CreateAccountTestSuite) SetupSuite() {
	router, repository, service, err := testutils.DBRouterSetup()
	if err != nil {
		s.T().Errorf("Error setting up router: %v", err)
	}

	s.Router = router
	s.Repository = repository
	s.Service = service
}

func (s *CreateAccountTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	if err != nil {
		s.T().Errorf("Error tearing down suite: %v", err)
	}
}

func (s *CreateAccountTestSuite) TestCreateAccount() {
	userJson := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey", "Password":"stringthatispassword"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signup", strings.NewReader(userJson))
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), 200, w.Code)

	userJson2 := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey"}`

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/signup", strings.NewReader(userJson2))
	s.Router.ServeHTTP(w2, req2)

	assert.Equal(s.T(), 500, w2.Code)

}
