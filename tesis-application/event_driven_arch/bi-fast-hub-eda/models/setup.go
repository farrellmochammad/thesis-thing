package models

import (
	"fmt"
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

	fmt.Println("Setup Database")
	// Check if the User table exists and insert seed data if it doesn't
	// if result := DB.Migrator().HasTable(&Bank{}); result {
	// 	var bank Bank
	// 	fmt.Println("Testing")
	// 	if err := DB.First(&bank).Error; errors.Is(err, gorm.ErrRecordNotFound) {
	// 		DB.Create(&Bank{BankURL: "http://localhost:8083", BankCode: "1", BankName: "Mandiri"})
	// 		DB.Create(&Bank{BankURL: "http://localhost:8085", BankCode: "2", BankName: "BRI"})
	// 	}
	// }

	fmt.Println("Database connected successfully.")

	return DB, err
}
