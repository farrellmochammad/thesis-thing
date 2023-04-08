package controllers

import (
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

	if isValidateAmount {
		middleware.JkdPost("http://localhost:8084/prm-processtransaction", fraudtransaction)
		return
	}

}
