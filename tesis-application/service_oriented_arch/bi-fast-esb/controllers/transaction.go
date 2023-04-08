package controllers

import (
	"errors"
	"net/http"

	"bi-fast-esb/logic"
	"bi-fast-esb/middleware"
	"bi-fast-esb/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FindTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var transactions []models.Transaction
	result := db.Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

func PrmProcessTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input models.FraudTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var transaction models.Transaction
	result := db.Where("transaction_hash = ?", input.TransactionID).First(&transaction)

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

	resultTransaction := models.ResultTransaction{
		FraudTransaction: input,
		Transaction:      transaction,
	}

	middleware.JkdPost(bankSender.BankURL+"/validatetransaction", resultTransaction)
}

func ProcessTransaction(c *gin.Context) {

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.JkdPost("http://localhost:8086/prm-processtransaction", input)

}

func ExecuteTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bank, isValidateReceiverBank, _ := logic.ValidateBankReceiver(db, input)
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

	if isValidateAmount {
		middleware.JkdPost(bank.BankURL+"/validatetransaction", input)
	}

}

func ValidateTransaction(c *gin.Context) {

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}

func UpdateTransaction(c *gin.Context) {
	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.JkdPut("http://localhost:8087/updatetransaction", input)

	// var input models.Transaction
	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// isSucess := logic.UpdateBalance(db, input)

	// if isSucess {
	// 	c.JSON(http.StatusAccepted, gin.H{"Status": "OK"})
	// 	bankReceiver, _, _ := logic.ValidateBankReceiver(db, input)
	// 	bankSender, _, _ := logic.ValidateBankSender(db, input)

	// 	middleware.JkdPost(bankReceiver.BankURL+"/successtransaction", input)
	// 	middleware.JkdPost(bankSender.BankURL+"/successtransaction", input)
	// 	return
	// } else {
	// 	c.JSON(http.StatusBadRequest, gin.H{"Status": "Not permitted"})
	// 	bankReceiver, _, _ := logic.ValidateBankReceiver(db, input)
	// 	bankSender, _, _ := logic.ValidateBankReceiver(db, input)

	// 	middleware.JkdPost(bankReceiver.BankURL+"/failedtransaction", input)
	// 	middleware.JkdPost(bankSender.BankURL+"/failedtransaction", input)
	// 	return
	// }

}

func BiHubSuccessTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bankReceiver, _, _ := logic.ValidateBankReceiver(db, input)
	bankSender, _, _ := logic.ValidateBankSender(db, input)

	middleware.JkdPost(bankReceiver.BankURL+"/retrievetransaction", input)
	middleware.JkdPost(bankSender.BankURL+"/successtransaction", input)

}

func BiHubFailedTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bankReceiver, _, _ := logic.ValidateBankReceiver(db, input)
	bankSender, _, _ := logic.ValidateBankSender(db, input)

	middleware.JkdPost(bankReceiver.BankURL+"/retrievetransaction", input)
	middleware.JkdPost(bankSender.BankURL+"/failedtransaction", input)
}
