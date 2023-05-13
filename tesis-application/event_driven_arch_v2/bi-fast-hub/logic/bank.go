package logic

import (
	"bi-fast-hub/models"
	"errors"

	"gorm.io/gorm"
)

func ValidateBankReceiver(db *gorm.DB, input models.Transaction) (models.Bank, bool, error) {
	var bank models.Bank
	result := db.Where("bank_code = ?", input.ReceiverBankCode).First(&bank)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// handle record not found error
			return models.Bank{}, false, nil
		} else {
			// handle other errors
			return models.Bank{}, false, result.Error
		}
	}

	return bank, true, nil

}

func ValidateBankSender(db *gorm.DB, input models.Transaction) (models.Bank, bool, error) {
	var bank models.Bank
	result := db.Where("bank_code = ?", input.SenderBankCode).First(&bank)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// handle record not found error
			return models.Bank{}, false, nil
		} else {
			// handle other errors
			return models.Bank{}, false, result.Error
		}
	}

	return bank, true, nil

}
