package servicemocks_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type CharacterServiceTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Service    *services.Service
}

func TestCharacterServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CharacterServiceTestSuite))
}

func (s *CharacterServiceTestSuite) SetupSuite() {
	router, repository, service, err := testutils.DBRouterSetup()
	if err != nil {
		s.T().Errorf("Error setting up router: %v", err)
	}

	s.Router = router
	s.Repository = repository
	s.Service = service
}

func (s *CharacterServiceTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	if err != nil {
		s.T().Errorf("Error tearing down suite: %v", err)
	}
}

func (s *CharacterServiceTestSuite) TestGetAndStoreAllCharacters() {
	err := s.Service.BagItemService.GetAndStoreAllCharacters("accountid", "apikeystring")
	assert.NoError(s.T(), err, "Failed to get and store items")
}

func (s *CharacterServiceTestSuite) TestGetBagItemsByCharacterName() {
	_, err := s.Service.BagItemService.BagItemRepository.GetByCharacterName("accountid", "Roman Meows")
	assert.NoError(s.T(), err, "Failed to get item by id")
}
