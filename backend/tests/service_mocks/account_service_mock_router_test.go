package servicemocks_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/api/handlers"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type AccountRouterServiceTestSuite struct {
	suite.Suite
	Router         *gin.Engine
	Repository     *repositories.Repository
	Service        *services.Service
	AccountHandler *handlers.AccountHandler
}

func TestAccountRouterServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AccountRouterServiceTestSuite))
}

func (s *AccountRouterServiceTestSuite) SetupSuite() {
	router, repository, service, err := testutils.DBRouterSetup()
	if err != nil {
		s.T().Errorf("Error setting up router: %v", err)
	}

	s.Router = router
	s.Repository = repository
	s.Service = service
	s.AccountHandler = handlers.NewAccountHandler(&repository.AccountRepository, &repository.SessionRepository, &repository.BagItemRepository, service.AccountService, service.BagItemService)
}

func (s *AccountRouterServiceTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	if err != nil {
		s.T().Errorf("Error tearing down suite: %v", err)
	}
}

func (s *AccountRouterServiceTestSuite) TestGetAccount() {
	account, err := s.Service.AccountService.GetAccount("apiKey")
	testutils.PrintObject(account)
	assert.NoError(s.T(), err, "Failed to get account")
}

func (s *AccountRouterServiceTestSuite) TestGetTokenInfo() {
	token, err := s.Service.AccountService.GetTokenInfo("apiKey")
	testutils.PrintObject(token)
	assert.NoError(s.T(), err, "Failed to get account")
	assert.Equal(s.T(), "armourytest", *token.Name, "Token info returns correct name")

}
