package handlerroutes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/api/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type BagItemHandlerTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Service    *services.Service
	Cookie     *http.Cookie
}

func TestBagItemHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(BagItemHandlerTestSuite))
}

func (s *BagItemHandlerTestSuite) SetupSuite() {
	router, repository, service, err := testutils.DBRouterSetup()
	if err != nil {
		s.T().Errorf("Error setting up router: %v", err)
	}

	s.Router = router
	s.Repository = repository
	s.Service = service

	s.T().Log("Create account with POST /apikeys")
	userJson := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey", "Password":"stringthatispassword"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	s.Router.ServeHTTP(w, req)
	assert.Equal(s.T(), http.StatusOK, w.Code)
	cookie := w.Result().Cookies()

	s.Cookie = cookie[0]
}

func (s *BagItemHandlerTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	if err != nil {
		s.T().Errorf("Error tearing down suite: %v", err)
	}
}

func (s *BagItemHandlerTestSuite) TestGetByAccount() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/account/inventory", nil)
	req.AddCookie(s.Cookie)
	s.Router.ServeHTTP(w, req)

	responseBagItems, err := testutils.UnmarshalToType[[]models.BagItem](w)
	if err != nil {
		s.T().Errorf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(s.T(), 200, w.Code)

	bagItemsResponseOK := BagItemsResponseOK(responseBagItems)
	assert.Equal(s.T(), true, bagItemsResponseOK, "BagItem response OK")

	allSameCharacterName := BagItemsAllSameCharacterName(responseBagItems)
	assert.Equal(s.T(), false, allSameCharacterName, "BagItems should belong to multiple different characters")
}

func (s *BagItemHandlerTestSuite) TestGetByCharacterName() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/account/characters/Roman%20Meows/inventory", nil)
	req.AddCookie(s.Cookie)
	s.Router.ServeHTTP(w, req)

	responseBagItems, err := testutils.UnmarshalToType[[]models.BagItem](w)
	if err != nil {
		s.T().Errorf("Failed to unmarshal response: %v", err)
	}

	bagItemsResponseOK := BagItemsResponseOK(responseBagItems)
	assert.Equal(s.T(), true, bagItemsResponseOK)

	allSameCharacterName := BagItemsAllSameCharacterName(responseBagItems)
	assert.Equal(s.T(), true, allSameCharacterName)
	assert.Equal(s.T(), 200, w.Code)
}

func BagItemsResponseOK(bagItems *[]models.BagItem) bool {
	if len(*bagItems) == 0 {
		return false
	} else {
		return true
	}
}

func BagItemsAllSameCharacterName(bagItems *[]models.BagItem) bool {
	characterName := (*bagItems)[0].CharacterName
	for _, bagItem := range *bagItems {
		if bagItem.CharacterName == characterName {
			continue
		} else {
			return false
		}
	}
	return true
}
