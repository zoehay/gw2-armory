package handlerroutes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	dbmodels "github.com/zoehay/gw2-armory/backend/internal/db/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type SessionMiddlewareTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Cookie     *http.Cookie
}

func TestSessionMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(SessionMiddlewareTestSuite))
}

func (s *SessionMiddlewareTestSuite) SetupSuite() {
	router, repository, _, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Router = router
	s.Repository = repository
}

func (s *SessionMiddlewareTestSuite) SetupTest() {
	err := s.Repository.AccountRepository.DB.Exec("TRUNCATE TABLE db_bag_item_infusions, db_bag_item_upgrades, db_bag_items, db_sessions, db_accounts CASCADE").Error
	s.Require().NoError(err, "Error truncating tables")

	userJson := `{"AccountName":"Name forAccount", "APIKey":"stringthatisapikey"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)
	s.Require().Equal(http.StatusOK, w.Code, "Setup: POST /apikeys must succeed")

	cookies := w.Result().Cookies()
	s.Require().NotEmpty(cookies, "Setup: expected sessionID cookie from POST /apikeys")
	s.Cookie = cookies[0]
}

func (s *SessionMiddlewareTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items", "db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *SessionMiddlewareTestSuite) getAccountInfo(cookie *http.Cookie) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/account/info", nil)
	if cookie != nil {
		req.AddCookie(cookie)
	}
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *SessionMiddlewareTestSuite) assertForbidden(w *httptest.ResponseRecorder, message string) {
	assert.Equal(s.T(), http.StatusForbidden, w.Code)
	response, err := testutils.UnmarshalResponse(w)
	s.Require().NoError(err, "Failed to unmarshal response")
	assert.Equal(s.T(), message, response["error"])
}

func (s *SessionMiddlewareTestSuite) TestValidSession() {
	w := s.getAccountInfo(s.Cookie)
	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *SessionMiddlewareTestSuite) TestMissingCookie() {
	w := s.getAccountInfo(nil)
	s.assertForbidden(w, "missing session cookie")
}

func (s *SessionMiddlewareTestSuite) TestEmptyCookie() {
	w := s.getAccountInfo(&http.Cookie{Name: "sessionID", Value: ""})
	s.assertForbidden(w, "empty session cookie")
}

func (s *SessionMiddlewareTestSuite) TestUnknownSession() {
	w := s.getAccountInfo(&http.Cookie{Name: "sessionID", Value: "notarealsession"})
	s.assertForbidden(w, "session not found")
}

func (s *SessionMiddlewareTestSuite) TestExpiredSession() {
	err := s.Repository.SessionRepository.DB.Model(&dbmodels.DBSession{}).
		Where("session_id = ?", s.Cookie.Value).
		Update("expires", time.Now().Add(-time.Hour)).Error
	s.Require().NoError(err, "Error expiring session")

	w := s.getAccountInfo(s.Cookie)
	s.assertForbidden(w, "session is expired")
}
