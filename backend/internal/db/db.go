package db

import (
	"errors"
	"log"

	"github.com/lib/pq"
	dbmodels "github.com/zoehay/gw2-armory/backend/internal/db/models"
	"github.com/zoehay/gw2-armory/backend/internal/db/repositories"
	"github.com/zoehay/gw2-armory/backend/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresInit(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("Error initializing database connection", err)
	}

	log.Print("Run db migrate")
	err = db.AutoMigrate(&dbmodels.DBItem{}, &dbmodels.DBBagItem{}, &dbmodels.DBAccount{}, &dbmodels.DBSession{})
	if err != nil {
		return nil, err
	}

	if err := migrateArraysToJoinTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

// migrateArraysToJoinTables backfills db_bag_item_infusions and db_bag_item_upgrades
// from the infusions/upgrades integer array columns, then drops those columns.
// Safe to remove once all databases have been migrated.
func migrateArraysToJoinTables(db *gorm.DB) error {
	var colCount int64
	db.Raw(`SELECT COUNT(*) FROM information_schema.columns WHERE table_name = 'db_bag_items' AND column_name = 'infusions'`).Scan(&colCount)
	if colCount == 0 {
		return nil
	}

	log.Print("Migrating infusion/upgrade arrays to join tables")

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	rows, err := sqlDB.Query(`SELECT id, infusions, upgrades FROM db_bag_items WHERE infusions IS NOT NULL OR upgrades IS NOT NULL`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id uint
		var infusions pq.Int64Array
		var upgrades pq.Int64Array
		if err := rows.Scan(&id, &infusions, &upgrades); err != nil {
			return err
		}
		for _, infID := range infusions {
			stub := dbmodels.DBItem{ID: uint(infID)}
			db.FirstOrCreate(&stub, dbmodels.DBItem{ID: uint(infID)})
			db.Exec(`INSERT INTO db_bag_item_infusions (db_bag_item_id, db_item_id) VALUES (?, ?) ON CONFLICT DO NOTHING`, id, infID)
		}
		for _, upgID := range upgrades {
			stub := dbmodels.DBItem{ID: uint(upgID)}
			db.FirstOrCreate(&stub, dbmodels.DBItem{ID: uint(upgID)})
			db.Exec(`INSERT INTO db_bag_item_upgrades (db_bag_item_id, db_item_id) VALUES (?, ?) ON CONFLICT DO NOTHING`, id, upgID)
		}
	}

	db.Exec(`ALTER TABLE db_bag_items DROP COLUMN IF EXISTS infusions`)
	db.Exec(`ALTER TABLE db_bag_items DROP COLUMN IF EXISTS upgrades`)

	log.Print("Array migration complete")
	return nil
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
