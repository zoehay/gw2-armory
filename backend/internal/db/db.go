package db

import (
	"errors"
	"log"

	dbmodels "github.com/zoehay/gw2-armory/backend/internal/db/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresInit(dsn string) (*gorm.DB, error) {

	// Add ping db
	// time.Sleep(30 * time.Second)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error initializing database connection", err)
	}

	log.Print("Run db migrate")
	err = db.AutoMigrate(&dbmodels.DBItem{}, &dbmodels.DBBagItem{}, &dbmodels.DBAccount{}, &dbmodels.DBSession{})
	if err != nil {
		return nil, err
	}

	return db, nil

}

func SeedItems(itemRepository repositories.ItemRepository, itemService services.ItemService) error {
	_, err := itemRepository.GetFirst()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Print("Seeding database")
		err = itemService.GetAndStoreAllItems()
		if err != nil {
			return err
		}
	} else {
		log.Print("Database already seeded")
	}

	return nil
}
