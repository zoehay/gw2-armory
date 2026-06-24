package gincontext_test

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

type GuestSessionInventoryAccessTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Service    *services.Service
}

func TestGuestSessionInventoryAccessSuite(t *testing.T) {
	suite.Run(t, new(GuestSessionInventoryAccessTestSuite))
}

func (s *GuestSessionInventoryAccessTestSuite) SetupSuite() {
	router, repository, service, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Router = router
	s.Repository = repository
	s.Service = service

	err = s.Service.ItemService.FetchAndStoreAllItems()
	s.Require().NoError(err, "Error seeding items")
}

func (s *GuestSessionInventoryAccessTestSuite) SetupTest() {
	err := s.Repository.AccountRepository.DB.Exec("TRUNCATE TABLE db_bag_item_infusions, db_bag_item_upgrades, db_bag_items, db_sessions, db_accounts").Error
	s.Require().NoError(err, "Error truncating tables")
}

func (s *GuestSessionInventoryAccessTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *GuestSessionInventoryAccessTestSuite) TestNoCookieNoInventoryAccess() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/account/characters/Roman%20Meows/inventory", nil)
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusForbidden, w.Code)
}

func (s *GuestSessionInventoryAccessTestSuite) TestInvalidSessionCookieAccess() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/account/characters/Roman%20Meows/inventory", nil)
	req.AddCookie(&http.Cookie{Name: "sessionID", Value: "not-a-real-session-id"})
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusForbidden, w.Code)
}

func (s *GuestSessionInventoryAccessTestSuite) TestGuestInventoryAccess() {
	userJson := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey", "Password":"stringthatispassword"}`
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	req1.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w1, req1)
	s.Require().Equal(http.StatusOK, w1.Code, "POST /apikeys must succeed")

	cookies := w1.Result().Cookies()
	s.Require().NotEmpty(cookies, "Expected sessionID cookie from POST /apikeys")

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/account/characters/Roman%20Meows/inventory", nil)
	req2.AddCookie(cookies[0])
	s.Router.ServeHTTP(w2, req2)

	assert.Equal(s.T(), http.StatusOK, w2.Code)

	inventory, err := testutils.UnmarshalToType[[]models.BagItem](w2)
	s.Require().NoError(err, "Failed to unmarshal inventory response")
	assert.NotEmpty(s.T(), *inventory, "Expected non-empty inventory for Roman Meows")
}
