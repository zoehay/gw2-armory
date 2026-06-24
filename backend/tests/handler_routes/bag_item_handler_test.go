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
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type BagItemHandlerTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Cookie     *http.Cookie
}

func TestBagItemHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(BagItemHandlerTestSuite))
}

func (s *BagItemHandlerTestSuite) SetupSuite() {
	router, repository, _, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Router = router
	s.Repository = repository

	userJson := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey", "Password":"stringthatispassword"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)
	s.Require().Equal(http.StatusOK, w.Code, "Setup: POST /apikeys must succeed")

	cookies := w.Result().Cookies()
	s.Require().NotEmpty(cookies, "Setup: expected sessionID cookie from POST /apikeys")
	s.Cookie = cookies[0]
}

func (s *BagItemHandlerTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *BagItemHandlerTestSuite) TestGetByAccount() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/account/inventory", nil)
	req.AddCookie(s.Cookie)
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	responseBagItems, err := testutils.UnmarshalToType[[]models.BagItem](w)
	s.Require().NoError(err, "Failed to unmarshal response")
	assert.NotEmpty(s.T(), *responseBagItems, "Expected non-empty bag items")

	assert.False(s.T(), bagItemsAllSameCharacterName(responseBagItems), "BagItems should belong to multiple different characters")
}

func (s *BagItemHandlerTestSuite) TestGetByCharacterName() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/account/characters/Roman%20Meows/inventory", nil)
	req.AddCookie(s.Cookie)
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	responseBagItems, err := testutils.UnmarshalToType[[]models.BagItem](w)
	s.Require().NoError(err, "Failed to unmarshal response")
	assert.NotEmpty(s.T(), *responseBagItems, "Expected non-empty bag items for character")
	assert.True(s.T(), bagItemsAllSameCharacterName(responseBagItems), "All items should belong to the same character")
}

func bagItemsAllSameCharacterName(bagItems *[]models.BagItem) bool {
	characterName := (*bagItems)[0].CharacterName
	for _, bagItem := range *bagItems {
		if bagItem.CharacterName != characterName {
			return false
		}
	}
	return true
}
