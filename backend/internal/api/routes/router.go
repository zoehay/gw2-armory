package routes

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zoehay/gw2armoury/backend/internal/api/handlers"
	"github.com/zoehay/gw2armoury/backend/internal/api/middleware"
	"github.com/zoehay/gw2armoury/backend/internal/db"
	"github.com/zoehay/gw2armoury/backend/internal/db/repositories"
	"github.com/zoehay/gw2armoury/backend/internal/services"
)

func LoadEnvDSN() string {
	// replace env with docker secrets
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	appMode := os.Getenv("APP_ENV")
	dsn := os.Getenv("DB_DSN")

	if appMode == "test" {
		dsn = os.Getenv("TEST_DB_DSN")
	} else if appMode == "development" {
		dsn = os.Getenv("DEV_NO_MOCK_DB_DSN")
	}
	// else if appMode == "docker-test" {
	// 	dsn = os.Getenv("DOCKER_TEST_DB_DSN")
	// }
	return dsn
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

	router := gin.Default()

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
