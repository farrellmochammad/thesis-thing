package logic

import (
	"bi-fast-hub-eda/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

func ValidateAccount(db *gorm.DB, input models.Transaction) bool {
	var account models.BankAccount
	result := db.Where("account_number = ?", input.ReceiverAccount).First(&account)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// handle record not found error
			return false
		} else {
			// handle other errors
			return false
		}
		return false
	}

	return true

}

func ValidateAmount(db *gorm.DB, input models.Transaction) bool {
	var account models.BankAccount
	result := db.Where("account_number = ?", input.ReceiverAccount).First(&account)

	if result.Error != nil {
		return false
	}

	if account.Balance-input.Amount < 0 {
		return false
	} else {
		return true
	}
}

func UpdateBalance(db *gorm.DB, input models.Transaction) bool {
	var bankAccount models.BankAccount
	err := db.Model(&bankAccount).Where("account_number = ?", input.ReceiverAccount).Update("balance", gorm.Expr("balance + ?", input.Amount)).Error
	if err != nil {
		log.Fatal(err)
		return false
	}

	err = db.Model(&bankAccount).Where("account_number = ?", input.SenderAccountNumber).Update("balance", gorm.Expr("balance - ?", input.Amount)).Error
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}
