package servicemocks_test

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type BagItemAccountServiceTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Service    *services.Service
}

func TestBagItemAccountServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BagItemAccountServiceTestSuite))
}

func (s *BagItemAccountServiceTestSuite) SetupSuite() {
	router, repository, service, err := testutils.DBRouterSetup()
	if err != nil {
		s.T().Errorf("Error setting up router: %v", err)
	}

	s.Router = router
	s.Repository = repository
	s.Service = service
}

func (s *BagItemAccountServiceTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	if err != nil {
		s.T().Errorf("Error tearing down suite: %v", err)
	}
}

func (s *BagItemAccountServiceTestSuite) TestGetAndStoreAllBagItems() {
	err := s.Service.BagItemService.GetAndStoreAllBagItems("accountid", "apikeystring")
	assert.NoError(s.T(), err, "Failed to get and store items")
	bagItems, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err, "Failed to get account bag items")
	numberAllBagItems := len(bagItems)
	fmt.Println(numberAllBagItems)
	assert.Equal(s.T(), 71, numberAllBagItems, "Correct number of items added")
}

func (s *BagItemAccountServiceTestSuite) TestGetAndStoreSharedInventory() {
	err := s.Service.BagItemService.GetAndStoreSharedInventory("accountid", "apikeystring")
	assert.NoError(s.T(), err, "Failed to get and store items")
}

func (s *BagItemAccountServiceTestSuite) TestGetBagItemsByAccount() {
	_, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err, "Failed to get account bag items")
}

func (s *BagItemAccountServiceTestSuite) TestClearAccountBagItems() {
	err := s.Service.BagItemService.GetAndStoreSharedInventory("accountid", "apikeystring")
	assert.NoError(s.T(), err, "Failed to get and store account inventory")
	err = s.Service.BagItemService.GetAndStoreAllCharacters("accountid", "apikeystring")
	assert.NoError(s.T(), err, "Failed to get and store character inventory")
	bagItems, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err, "Failed to get account bag items")
	numberAllBagItems := len(bagItems)
	fmt.Println(numberAllBagItems)
	err = s.Service.BagItemService.ClearSharedInventory("accountid")
	assert.NoError(s.T(), err, "Failed to clear account shared inventory")
	bagItems, err = s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err, "Failed to get account bag items")
	numberCharacterBagItems := len(bagItems)
	fmt.Println(numberCharacterBagItems)
	assert.Equal(s.T(), 4, (numberAllBagItems - numberCharacterBagItems))

}

func (s *BagItemAccountServiceTestSuite) TestClearCharacterInventory() {
	err := s.Service.BagItemService.GetAndStoreSharedInventory("accountid", "apikeystring")
	assert.NoError(s.T(), err, "Failed to get and account inventory")
	err = s.Service.BagItemService.GetAndStoreAllCharacters("accountid", "apikeystring")
	assert.NoError(s.T(), err, "Failed to get and store character inventory")
	bagItems, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err, "Failed to get account bag items")
	numberAllBagItems := len(bagItems)
	fmt.Println(numberAllBagItems)
	err = s.Service.BagItemService.ClearCharacterInventory("accountid", "Laura Lesdottir")
	assert.NoError(s.T(), err, "Failed to clear character inventory")
	err = s.Service.BagItemService.ClearCharacterInventory("accountid", "Roman Meows")
	assert.NoError(s.T(), err, "Failed to clear character inventory")
	bagItems, err = s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err, "Failed to get account bag items")
	numberBagItemsWithoutCharacter := len(bagItems)
	fmt.Println(numberBagItemsWithoutCharacter)
	assert.Equal(s.T(), 67, (numberAllBagItems - numberBagItemsWithoutCharacter))

}
