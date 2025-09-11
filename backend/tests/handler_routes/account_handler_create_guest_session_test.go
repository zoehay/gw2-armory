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

type CreateGuestSessionTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
	Service    *services.Service
}

func TestCreateGuestSessionTestSuite(t *testing.T) {
	suite.Run(t, new(CreateGuestSessionTestSuite))
}

func (s *CreateGuestSessionTestSuite) SetupSuite() {
	router, repository, service, err := testutils.DBRouterSetup()
	if err != nil {
		s.T().Errorf("Error setting up router: %v", err)
	}

	s.Router = router
	s.Repository = repository
	s.Service = service
}

func (s *CreateGuestSessionTestSuite) TearDownSuite() {
	dropTables := []string{"db_accounts", "db_sessions", "db_bag_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	if err != nil {
		s.T().Errorf("Error tearing down suite: %v", err)
	}
}

func (s *CreateGuestSessionTestSuite) TestCreateGuestWithNewAPIKeyCreatesSession() {
	gin.SetMode(gin.TestMode)

	userJson := `{"APIKey":"stringapikey"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/apikeys", strings.NewReader(userJson))
	s.Router.ServeHTTP(w, req)
	assert.Equal(s.T(), 200, w.Code)

	dbAccount, err := s.Repository.AccountRepository.GetByID("gw2apiaccountidstring")
	if err != nil {
		s.T().Errorf("Error getting account from db: %v", err)
	}

	account, err := testutils.UnmarshalToType[models.Account](w)
	if err != nil {
		s.T().Errorf("Failed to unmarshal response: %v", err)
	}

	cookieSessionID := w.Result().Cookies()[0].Value

	assert.Equal(s.T(), dbAccount.SessionID, account.SessionID, "SessionID in db matches returned account")
	assert.Equal(s.T(), dbAccount.SessionID, &cookieSessionID, "SessionID in db matches returned cookie")

}
