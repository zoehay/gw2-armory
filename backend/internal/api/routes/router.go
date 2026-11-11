package routes

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zoehay/gw2-armory/backend/internal/api/handlers"
	"github.com/zoehay/gw2-armory/backend/internal/api/middleware"
	"github.com/zoehay/gw2-armory/backend/internal/db"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
)

func SetupRouter(allowedOrigins []string, secureCookies bool, dsn string, mocks bool) (*gin.Engine, *repositories.Repository, *services.Service, error) {
	database, err := db.PostgresInit(dsn)
	if err != nil {
		log.Fatal("Error initializing database connection", err)
	}

	repository := repositories.NewRepository(database)
	service, err := services.NewService(repository, mocks)
	if err != nil {
		log.Fatal("Error setting up services", err)
	}

	itemHandler := handlers.NewItemHandler(service.ItemService)
	bagItemHandler := handlers.NewBagItemHandler(service.BagItemService)
	accountHandler := handlers.NewAccountHandler(secureCookies, service.AccountService, service.BagItemService)

	err = db.SeedItems(repository.ItemRepository, service.ItemService)
	if err != nil {
		log.Fatal("Error seeding database", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	err = router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	router.Use(middleware.SetCORS(allowedOrigins))

	router.GET("/items", itemHandler.GetAllItems)
	router.GET("/items/:id", itemHandler.GetItemByID)

	router.POST("/apikeys", accountHandler.HandlePostAPIKeyRequest)
	// outside the session middleware so an expired or unknown session can still clear its cookie
	router.POST("/account/logout", accountHandler.Logout)

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

func LoadEnv() (allowedOrigins []string, secureCookies bool) {
	allowedOriginsString, ok := os.LookupEnv("CORS_ALLOW_ORIGIN")
	if !ok {
		allowedOriginsString = "https://localhost:5173 http://localhost:5173"
		fmt.Printf("CORS_ALLOW_ORIGIN not set, default to %s", allowedOriginsString)
		os.Setenv("CORS_ALLOW_ORIGIN", allowedOriginsString)
	}
	allowedOrigins = strings.Split(allowedOriginsString, " ")

	secureCookies = cookieSecureFromEnv()
	if !secureCookies {
		fmt.Println("COOKIE_SECURE=false, session cookies will be sent over plain HTTP")
	}

	return allowedOrigins, secureCookies
}

// cookieSecureFromEnv defaults to secure cookies, so prod doesn't need to set COOKIE_SECURE.
// Only an explicit false value, for local plain-HTTP dev, turns it off
func cookieSecureFromEnv() bool {
	value, ok := os.LookupEnv("COOKIE_SECURE")
	if !ok {
		return true
	}
	secure, err := strconv.ParseBool(value)
	if err != nil {
		return true
	}
	return secure
}

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
