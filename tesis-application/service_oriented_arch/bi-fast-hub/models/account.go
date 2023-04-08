package models

import "gorm.io/gorm"

type BankAccount struct {
	gorm.Model
	AccountNumber string `gorm:"uniqueIndex"`
	BankCode      string
	OwnerName     string
	Balance       float64
}
