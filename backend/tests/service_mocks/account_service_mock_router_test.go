package servicemocks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type AccountServiceTestSuite struct {
	suite.Suite
	Repository *repositories.Repository
	Service    *services.Service
}

func TestAccountServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceTestSuite))
}

func (s *AccountServiceTestSuite) SetupSuite() {
	_, repository, service, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Repository = repository
	s.Service = service
}

func (s *AccountServiceTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *AccountServiceTestSuite) TestFetchAccount() {
	account, err := s.Service.AccountService.FetchAccount("apiKey")
	assert.NoError(s.T(), err, "Failed to fetch account")
	assert.Equal(s.T(), "gw2apiaccountidstring", account.AccountID)
	assert.NotNil(s.T(), account.GW2AccountName)
	assert.Equal(s.T(), "gw2name", *account.GW2AccountName)
}

func (s *AccountServiceTestSuite) TestFetchToken() {
	token, err := s.Service.AccountService.FetchToken("apiKey")
	assert.NoError(s.T(), err, "Failed to fetch token")
	assert.Equal(s.T(), "armourytest", *token.Name, "Token info returns correct name")
}
