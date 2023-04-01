package models

import "gorm.io/gorm"

type Bank struct {
	gorm.Model
	BankURL  string `gorm:"uniqueIndex"`
	BankCode string
	BankName string
}
