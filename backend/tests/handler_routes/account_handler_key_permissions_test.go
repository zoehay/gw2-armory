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
	"github.com/zoehay/gw2-armory/backend/internal/gw2_client/providers"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type KeyPermissionsTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
}

func TestKeyPermissionsTestSuite(t *testing.T) {
	suite.Run(t, new(KeyPermissionsTestSuite))
}

func (s *KeyPermissionsTestSuite) SetupSuite() {
	router, repository, _, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Router = router
	s.Repository = repository
}

func (s *KeyPermissionsTestSuite) SetupTest() {
	err := s.Repository.AccountRepository.DB.Exec("TRUNCATE TABLE db_bag_item_infusions, db_bag_item_upgrades, db_bag_items, db_sessions, db_accounts CASCADE").Error
	s.Require().NoError(err, "Error truncating tables")
}

func (s *KeyPermissionsTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *KeyPermissionsTestSuite) postKey(apiKey string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(`{"APIKey":"`+apiKey+`"}`))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *KeyPermissionsTestSuite) assertNothingCreated(w *httptest.ResponseRecorder) {
	assert.Empty(s.T(), w.Result().Cookies(), "No session cookie for a rejected key")

	var accounts, sessions int64
	s.Require().NoError(s.Repository.AccountRepository.DB.Table("db_accounts").Count(&accounts).Error)
	s.Require().NoError(s.Repository.AccountRepository.DB.Table("db_sessions").Count(&sessions).Error)
	assert.Zero(s.T(), accounts, "No account stored for a rejected key")
	assert.Zero(s.T(), sessions, "No session stored for a rejected key")
}

func (s *KeyPermissionsTestSuite) TestKeyMissingPermissionsIsRejected() {
	w := s.postKey(providers.MockLimitedPermissionsKey)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.JSONEq(s.T(), `{
		"error": "API key is missing required permissions",
		"missing_permissions": ["characters", "wallet"]
	}`, w.Body.String())
	s.assertNothingCreated(w)
}

func (s *KeyPermissionsTestSuite) TestInvalidKeyIsRejected() {
	w := s.postKey(providers.MockInvalidKey)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.JSONEq(s.T(), `{"error": "invalid API key"}`, w.Body.String())
	s.assertNothingCreated(w)
}
