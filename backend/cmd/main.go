package main

import (
	"log"
	"os"

	"github.com/zoehay/gw2-armory/backend/internal/api/routes"
)

func main() {
	dsn, err := routes.LoadEnvDSN()
	if err != nil {
		log.Fatal("Error getting database dsn", err)
	}
	mocks := false
	appMode := os.Getenv("APP_ENV")
	if appMode == "test" || appMode == "docker-test" {
		mocks = true
	}

	router, _, _, err := routes.SetupRouter(dsn, mocks)
	if err != nil {
		log.Fatal("Error setting up router", err)
	}

	router.Run(":8000")
}
