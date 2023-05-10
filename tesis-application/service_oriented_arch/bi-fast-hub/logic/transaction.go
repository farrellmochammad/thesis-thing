package logic

import (
	"bi-fast-hub/models"
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

func ValidateBulkBankSender(db *gorm.DB, input models.BulkTransaction) bool {
	var account models.BankAccount
	result := db.Where("account_number = ?", input.SenderAccountNumber).First(&account)

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

func ValidateBulkAmount(db *gorm.DB, input models.BulkTransaction) bool {
	totalamount := GetTotalAmount(input)

	var account models.BankAccount
	result := db.Where("account_number = ?", input.SenderAccountNumber).First(&account)

	if result.Error != nil {
		return false
	}

	if account.Balance-totalamount < 0 {
		return false
	} else {
		return true
	}
}

func GetTotalAmount(bt models.BulkTransaction) float64 {
	var totalAmount float64
	for _, transaction := range bt.Transactions {
		totalAmount += transaction.Amount
	}
	return totalAmount
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
