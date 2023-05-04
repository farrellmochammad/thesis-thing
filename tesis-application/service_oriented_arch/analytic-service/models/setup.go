package models

import (
	"log"

	r "gopkg.in/gorethink/gorethink.v4"
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

	err = DB.AutoMigrate(&BankAccount{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	return DB, err
}

func CreateSession(options r.ConnectOpts) (*r.Session, error) {
	session, err := r.Connect(options)
	if err != nil {
		return nil, err
	}

	// if err := r.DBCreate("test").Exec(session); err != nil {
	// 	return nil, err
	// }

	// if err := r.DB("ci-connector-transaction").TableCreate("transactions").Exec(session); err != nil {
	// 	return nil, err
	// }

	// if err := r.DB("ci-connector-transaction").Table("transactions").IndexCreate("sender_account_number").Exec(session); err != nil {
	// 	return nil, err
	// }

	return session, nil
}
