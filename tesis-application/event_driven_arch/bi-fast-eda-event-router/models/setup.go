package models

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupModels() (*gorm.DB, error) {
	dsn := "host=localhost user=root password=root dbname=bifast_db port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = DB.AutoMigrate(&Transaction{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	return DB, err
}
