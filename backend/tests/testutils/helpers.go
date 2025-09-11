package testutils

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zoehay/gw2-armory/backend/internal/api/routes"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
)

func PrintObject(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(s))
	return string(s)
}

func UnmarshalResponse(w *httptest.ResponseRecorder) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func UnmarshalToType[T any](w *httptest.ResponseRecorder) (*T, error) {
	var obj T
	err := json.Unmarshal(w.Body.Bytes(), &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func DBRouterSetup() (*gin.Engine, *repositories.Repository, *services.Service, error) {
	envPath := filepath.Join("../..", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		return nil, nil, nil, err
	}

	dsn := os.Getenv("TEST_DB_DSN")
	router, repository, service, err := routes.SetupRouter(dsn, true)
	if err != nil {
		return nil, nil, nil, err
	}

	return router, repository, service, nil
}

func TearDownTruncateTables(repository *repositories.Repository, tables []string) error {
	for _, tableString := range tables {
		err := repository.AccountRepository.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %v cascade;", tableString)).Error
		if err != nil {
			return err
		}
	}

	db, err := repository.AccountRepository.DB.DB()
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
