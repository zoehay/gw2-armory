package routes

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zoehay/gw2-armory/backend/internal/api/handlers"
	"github.com/zoehay/gw2-armory/backend/internal/api/middleware"
	"github.com/zoehay/gw2-armory/backend/internal/db"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
)

func LoadEnvDSN() (string, error) {
	var dsn string
	// docker secrets
	if dbPasswordFile := os.Getenv("ARMORY_DB_PASSWORD_FILE"); dbPasswordFile != "" {
		data, err := os.ReadFile(dbPasswordFile)
		if err != nil {
			return "", fmt.Errorf("docker secret, failed to read DB_DSN_FILE: %w", err)
		}
		password := strings.TrimSpace(string(data))
		dsn = fmt.Sprintf("host=armory-db user=postgres password=%s dbname=armory port=5432", password)
	} else {
		// local dev env file
		err := godotenv.Load()
		if err != nil {
			return "", fmt.Errorf("local env, error loading .env file: %w", err)
		}
		appMode := os.Getenv("APP_ENV")
		if appMode == "test" {
			dsn = os.Getenv("TEST_DB_DSN")
		} else {
			dsn = os.Getenv("DEV_NO_MOCK_DB_DSN")
		}

	}
	return dsn, nil

}

func SetupRouter(dsn string, mocks bool) (*gin.Engine, *repositories.Repository, *services.Service, error) {
	database, err := db.PostgresInit(dsn)
	if err != nil {
		log.Fatal("Error initializing database connection", err)
	}

	repository := repositories.NewRepository(database)
	service := services.NewService(repository, mocks)

	itemHandler := handlers.NewItemHandler(&repository.ItemRepository)
	bagItemHandler := handlers.NewBagItemHandler(&repository.BagItemRepository, service.ItemService)
	accountHandler := handlers.NewAccountHandler(&repository.AccountRepository, &repository.SessionRepository, &repository.BagItemRepository, service.AccountService, service.BagItemService)

	err = db.SeedItems(repository.ItemRepository, *service.ItemService)
	if err != nil {
		log.Fatal("Error seeding database", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	err = router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	router.Use(middleware.SetCORS())

	router.GET("/items", itemHandler.GetAllItems)
	router.GET("/items/:id", itemHandler.GetItemByID)

	router.POST("/login", accountHandler.Login)
	router.POST("/signup", accountHandler.HandlePostAPIKeyRequest) //change signup handler with password verification
	router.POST("/apikeys", accountHandler.HandlePostAPIKeyRequest)

	account := router.Group("/account")
	account.Use(middleware.UseSession(&repository.AccountRepository, &repository.SessionRepository))
	{
		account.GET("/info", accountHandler.GetAccount)
		account.GET("/inventory", bagItemHandler.GetByAccount)
		account.GET("/characters/:charactername/inventory", bagItemHandler.GetByCharacter)
		account.DELETE("/delete", accountHandler.Delete)
		account.GET("/accountinventory", bagItemHandler.GetAccountInventory)
		account.POST("/searchinventory", bagItemHandler.GetFilteredAccountInventory)
	}

	return router, repository, service, nil

}
