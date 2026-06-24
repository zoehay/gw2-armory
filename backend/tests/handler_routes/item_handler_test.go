package handlerroutes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zoehay/gw2-armory/backend/internal/api/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/tests/testutils"
)

type ItemHandlerTestSuite struct {
	suite.Suite
	Router     *gin.Engine
	Repository *repositories.Repository
}

func TestItemHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ItemHandlerTestSuite))
}

func (s *ItemHandlerTestSuite) SetupSuite() {
	router, repository, _, err := testutils.DBRouterSetup()
	s.Require().NoError(err, "Error setting up router")
	s.Router = router
	s.Repository = repository
}

func (s *ItemHandlerTestSuite) TearDownSuite() {
	dropTables := []string{"db_items"}
	err := testutils.TearDownTruncateTables(s.Repository, dropTables)
	s.Require().NoError(err, "Error tearing down suite")
}

func (s *ItemHandlerTestSuite) TestGetItemByID() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/27952", nil)
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	item, err := testutils.UnmarshalToType[models.Item](w)
	s.Require().NoError(err, "Failed to unmarshal response")
	assert.Equal(s.T(), uint(27952), item.ID)
	assert.NotNil(s.T(), item.Name)
}

func (s *ItemHandlerTestSuite) TestGetItemByNonNumericID() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/notanumber", nil)
	s.Router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}
