package servicemocks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

// Item counts derived from mock test data files.
const (
	sharedInventoryCount = 4   // account_inventory_test_data.txt
	bankInventoryCount   = 150 // account_bank_test_data.txt
	romanMeowsCount      = 32  // character_test_data.txt: 15 bag + 17 equipment
	lauraLesdottirCount  = 35  // character_test_data.txt: 30 bag + 5 equipment
	allCharactersCount   = romanMeowsCount + lauraLesdottirCount
	totalInventoryCount  = sharedInventoryCount + allCharactersCount
)

type BagItemAccountServiceTestSuite struct {
	suite.Suite
	Repository *repositories.Repository
	Service    *services.Service
}

func TestBagItemAccountServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BagItemAccountServiceTestSuite))
}

func (s *BagItemAccountServiceTestSuite) SetupSuite() {
	_, repository, service, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Repository = repository
	s.Service = service
}

func (s *BagItemAccountServiceTestSuite) SetupTest() {
	err := s.Repository.AccountRepository.DB.Exec("TRUNCATE TABLE db_bag_item_infusions, db_bag_item_upgrades, db_bag_items").Error
	s.Require().NoError(err, "Error truncating bag item tables")
}

func (s *BagItemAccountServiceTestSuite) TearDownSuite() {
	err := testutils.TearDownTruncateTables(s.Repository, []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"})
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *BagItemAccountServiceTestSuite) TestStoreAllBagItems() {
	err := s.Service.BagItemService.FetchAndStoreAllBagItems("accountid", "apikeystring")
	assert.NoError(s.T(), err)

	items, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), totalInventoryCount, len(items))
}

func (s *BagItemAccountServiceTestSuite) TestStoreAllBagItemsIsIdempotent() {
	err := s.Service.BagItemService.FetchAndStoreAllBagItems("accountid", "apikeystring")
	assert.NoError(s.T(), err)
	first, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err)

	err = s.Service.BagItemService.FetchAndStoreAllBagItems("accountid", "apikeystring")
	assert.NoError(s.T(), err)
	second, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), len(first), len(second), "second sync should replace inventory, not duplicate it")
}

func (s *BagItemAccountServiceTestSuite) TestStoreSharedInventory() {
	err := s.Service.BagItemService.FetchAndStoreSharedInventory("accountid", "apikeystring")
	assert.NoError(s.T(), err)

	items, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), sharedInventoryCount, len(items))
}

func (s *BagItemAccountServiceTestSuite) TestStoreAllCharacters() {
	err := s.Service.BagItemService.FetchAndStoreAllCharacters("accountid", "apikeystring")
	assert.NoError(s.T(), err)

	items, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), allCharactersCount, len(items))
}

func (s *BagItemAccountServiceTestSuite) TestClearSharedInventory() {
	err := s.Service.BagItemService.FetchAndStoreAllBagItems("accountid", "apikeystring")
	assert.NoError(s.T(), err)

	err = s.Service.BagItemService.ClearSharedInventory("accountid")
	assert.NoError(s.T(), err)

	items, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), allCharactersCount, len(items))
}

func (s *BagItemAccountServiceTestSuite) TestClearCharacterInventory() {
	err := s.Service.BagItemService.FetchAndStoreAllBagItems("accountid", "apikeystring")
	assert.NoError(s.T(), err)

	err = s.Service.BagItemService.ClearCharacterInventory("accountid", "Laura Lesdottir")
	assert.NoError(s.T(), err)

	items, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), totalInventoryCount-lauraLesdottirCount, len(items))
}

func (s *BagItemAccountServiceTestSuite) TestStoreBankInventory() {
	err := s.Service.BagItemService.FetchAndStoreBankInventory("accountid", "apikeystring")
	assert.NoError(s.T(), err)

	items, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), bankInventoryCount, len(items))
}

func (s *BagItemAccountServiceTestSuite) TestClearBankInventory() {
	err := s.Service.BagItemService.FetchAndStoreAllBagItems("accountid", "apikeystring")
	assert.NoError(s.T(), err)

	err = s.Service.BagItemService.FetchAndStoreBankInventory("accountid", "apikeystring")
	assert.NoError(s.T(), err)

	err = s.Service.BagItemService.ClearBankInventory("accountid")
	assert.NoError(s.T(), err)

	items, err := s.Repository.BagItemRepository.GetDetailBagItemByAccountID("accountid")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), totalInventoryCount, len(items))
}
