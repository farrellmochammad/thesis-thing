package models

import "gorm.io/gorm"

type Bank struct {
	gorm.Model
	BankURL  string `gorm:"uniqueIndex"`
	BankCode string
	BankName string
}

type BankAccount struct {
	gorm.Model
	AccountNumber string `gorm:"uniqueIndex"`
	BankCode      string
	OwnerName     string
	Balance       float64
}
