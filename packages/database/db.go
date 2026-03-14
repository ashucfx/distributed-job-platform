package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(dsn string) *gorm.DB {
	var err error
	
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Retry connecting for docker compose wait-for-it pattern
	for i := 1; i <= 5; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to db (attempt %d/5): %v. Retrying in 2s...", i, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate schemas
	if err := DB.AutoMigrate(&Job{}, &JobLog{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return DB
}
