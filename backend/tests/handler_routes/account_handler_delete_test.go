package handlerroutes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type DeleteAccountTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Cookie     *http.Cookie
}

func TestDeleteAccountTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteAccountTestSuite))
}

func (s *DeleteAccountTestSuite) SetupSuite() {
	router, repository, _, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Router = router
	s.Repository = repository
}

func (s *DeleteAccountTestSuite) SetupTest() {
	err := s.Repository.AccountRepository.DB.Exec("TRUNCATE TABLE db_bag_item_infusions, db_bag_item_upgrades, db_bag_items, db_sessions, db_accounts CASCADE").Error
	s.Require().NoError(err, "Error truncating tables")

	userJson := `{"APIKey":"stringthatisapikey"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)
	s.Require().Equal(http.StatusOK, w.Code, "SetupTest: POST /apikeys must succeed")

	cookies := w.Result().Cookies()
	s.Require().NotEmpty(cookies, "SetupTest: expected sessionID cookie from POST /apikeys")
	s.Cookie = cookies[0]
}

func (s *DeleteAccountTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *DeleteAccountTestSuite) TestDeleteAccount() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/account/delete", strings.NewReader(`{"APIKey":"stringthatisapikey"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(s.Cookie)
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	_, err := s.Repository.AccountRepository.GetByID("gw2apiaccountidstring")
	assert.Error(s.T(), err, "Account should not exist after deletion")
}

func (s *DeleteAccountTestSuite) TestDeleteWithMalformedBody() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/account/delete", strings.NewReader("not valid json"))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(s.Cookie)
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}
