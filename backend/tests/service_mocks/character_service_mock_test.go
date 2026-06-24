package servicemocks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type CharacterServiceTestSuite struct {
	suite.Suite
	Repository *repositories.Repository
	Service    *services.Service
}

func TestCharacterServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CharacterServiceTestSuite))
}

func (s *CharacterServiceTestSuite) SetupSuite() {
	_, repository, service, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Repository = repository
	s.Service = service
}

func (s *CharacterServiceTestSuite) SetupTest() {
	err := s.Repository.AccountRepository.DB.Exec("TRUNCATE TABLE db_bag_item_infusions, db_bag_item_upgrades, db_bag_items").Error
	s.Require().NoError(err, "Error truncating bag item tables")
}

func (s *CharacterServiceTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *CharacterServiceTestSuite) TestFetchAndStoreAllCharacters() {
	err := s.Service.BagItemService.FetchAndStoreAllCharacters("accountid", "apikeystring")
	assert.NoError(s.T(), err, "Failed to fetch and store characters")

	items, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), allCharactersCount, len(items))
}

func (s *CharacterServiceTestSuite) TestGetBagItemsByCharacterName() {
	err := s.Service.BagItemService.FetchAndStoreAllCharacters("accountid", "apikeystring")
	s.Require().NoError(err, "Failed to seed character data")

	items, err := s.Repository.BagItemRepository.GetDetailBagItemByCharacterName("accountid", "Roman Meows")
	assert.NoError(s.T(), err, "Failed to get items by character name")
	assert.Equal(s.T(), romanMeowsCount, len(items), "Expected correct item count for Roman Meows")

	for _, item := range items {
		assert.NotNil(s.T(), item.CharacterName)
		assert.Equal(s.T(), "Roman Meows", *item.CharacterName)
	}
}
