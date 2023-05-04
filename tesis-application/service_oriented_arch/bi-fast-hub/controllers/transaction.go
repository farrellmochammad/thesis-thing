package controllers

import (
	"net/http"

	"bi-fast-hub/logic"
	"bi-fast-hub/middleware"
	"bi-fast-hub/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isSucess := logic.UpdateBalance(db, input)

	if isSucess {
		c.JSON(http.StatusAccepted, gin.H{"Status": "OK"})

		middleware.JkdPost("http://localhost:8084/bihub-successtransaction", input)
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Not permitted"})

		middleware.JkdPost("http://localhost:8084/bihub-failedtransaction", input)
		return
	}

}

func BiHubValidateTransaction(c *gin.Context) {

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

	if isValidateAmount {
		middleware.JkdPost("http://localhost:8084/bi-fast-esb/prm-processtransaction", input)
		return
	}

	middleware.JkdPost("http://localhost:8084/bi-fast-esb/failed-processtransaction", input)

}
