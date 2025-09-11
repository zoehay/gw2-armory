package tests

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	database "github.com/zoehay/gw2-armory/backend/internal/db"
	"gorm.io/gorm"
)

type DBTestSuite struct {
	suite.Suite
	dsn string
	db  *gorm.DB
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

func (s *DBTestSuite) SetupSuite() {
	envPath := filepath.Join("..", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	s.dsn = os.Getenv("TEST_DB_DSN")
}

func (s *DBTestSuite) TearDownSuite() {
	db, err := s.db.DB()
	if err != nil {
		s.T().Fatal(err)
	}
	db.Close()
}

func (s *DBTestSuite) TestPostgresInit() {
	db, err := database.PostgresInit(s.dsn)
	s.db = db
	assert.NoError(s.T(), err, "Failed to connect to database")
}
