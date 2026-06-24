package servicemocks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type ItemServiceTestSuite struct {
	suite.Suite
	Repository *repositories.Repository
	Service    *services.Service
}

func TestItemServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ItemServiceTestSuite))
}

func (s *ItemServiceTestSuite) SetupSuite() {
	_, repository, service, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Repository = repository
	s.Service = service

	err = s.Service.ItemService.ItemRepository.DB.Exec("TRUNCATE TABLE db_items;").Error
	s.Require().NoError(err, "Failed to clear items table before test")
}

func (s *ItemServiceTestSuite) TearDownSuite() {
	dropTables := []string{"db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")

	db, err := s.Service.ItemService.ItemRepository.DB.DB()
	s.Require().NoError(err, "Failed to get underlying DB")
	db.Close()
}

func (s *ItemServiceTestSuite) TestGetAndStoreAllItems() {
	err := s.Service.ItemService.FetchAndStoreAllItems()
	assert.NoError(s.T(), err, "Failed to fetch and store items")

	item, err := s.Service.ItemService.ItemRepository.GetById(27952)
	assert.NoError(s.T(), err, "Failed to get item by id")
	assert.Equal(s.T(), "Axiquiotl", item.Name, "Correct item name")
}
