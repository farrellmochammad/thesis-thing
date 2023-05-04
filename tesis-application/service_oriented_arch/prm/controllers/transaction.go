package controllers

import (
	"errors"
	"net/http"

	"prm/logic"
	"prm/middleware"
	"prm/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PrmProcessTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, isValidateReceiverBank, _ := logic.ValidateBankReceiver(db, input)
	isValidateAccount := logic.ValidateAccount(db, input)
	isValidateAmount := logic.ValidateAmount(db, input)

	if !isValidateReceiverBank {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Receiver bank doesn't exist"})
		return
	}

	if !isValidateAccount {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Account number doesn't exist"})
		return
	}

	if !isValidateAmount {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Amount not enough"})
		return
	}

	indexfraud, fraudcategory := logic.FraudDetection()

	fraudtransaction := models.FraudTransaction{
		TransactionID: input.TransactionHash,
		IndexFraud:    indexfraud,
		FraudCategory: fraudcategory,
	}

	var transaction models.Transaction
	result := db.Where("transaction_hash = ?", fraudtransaction.TransactionID).First(&transaction)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// handle record not found error
			c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction not found"})
		} else {
			// handle other errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	bankSender, _, _ := logic.ValidateBankSender(db, transaction)

	processtransaction := models.ProcessTransaction{
		Transaction:      transaction,
		FraudTransaction: fraudtransaction,
		BankSender:       bankSender.BankURL,
	}

	if isValidateAmount {
		middleware.JkdPost("http://localhost:8084/bi-fast-esb/report-prm-processtransaction", processtransaction)
		return
	}

}
