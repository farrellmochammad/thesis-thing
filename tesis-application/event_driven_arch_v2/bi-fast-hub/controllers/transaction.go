package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bi-fast-hub/logic"
	"bi-fast-hub/middleware"
	"bi-fast-hub/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func UpdateTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isSucess := logic.UpdateBalance(db, input)
	bankReceiver, _, _ := logic.ValidateBankReceiver(db, input)
	bankSender, _, _ := logic.ValidateBankSender(db, input)

	sentTransaction := models.SentTransaction{
		Transaction:  input,
		BankSender:   bankSender.BankURL,
		BankReceiver: bankReceiver.BankURL,
	}

	if isSucess {
		c.JSON(http.StatusAccepted, gin.H{"Status": "OK"})

		middleware.JkdPost("http://localhost:8084/bihub-successtransaction", sentTransaction)
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Not permitted"})

		middleware.JkdPost("http://localhost:8084/bihub-failedtransaction", sentTransaction)
		return
	}

}

func BiHubValidateTransaction(mqtt_client MQTT.Client, db *gorm.DB, payload string) {

	var input models.BulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	isValidateSenderBank := logic.ValidateBulkBankSender(db, input)
	if !isValidateSenderBank {
		fmt.Println("Receiver bank doesn't exist")
		return
	}

	isValidateAmount := logic.ValidateBulkAmount(db, input)
	if !isValidateAmount {
		fmt.Println("Amount not enough")
		return
	}

	if isValidateAmount {
		middleware.PublishMessage(mqtt_client, "topic/prm-process-bulk-transaction", input)
		return
	}

	// if isValidateAmount {
	// 	middleware.JkdPost("http://localhost:8084/bi-fast-esb/prm-processtransaction", input)
	// 	return
	// }

	// middleware.JkdPost("http://localhost:8084/bi-fast-esb/failed-processtransaction", input)

}

func BiHubValidateBulkTransaction(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var input models.BulkTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isValidateSenderBank := logic.ValidateBulkBankSender(db, input)
	if !isValidateSenderBank {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Receiver bank doesn't exist"})
		return
	}

	isValidateAmount := logic.ValidateBulkAmount(db, input)
	if !isValidateAmount {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Amount not enough"})
		return
	}

	if isValidateAmount {
		middleware.JkdPost("http://localhost:8084/bi-fast-esb/prm-processbulktransaction", input)
		return
	}

	middleware.JkdPost("http://localhost:8084/bi-fast-esb/failed-processtransaction", input)

}

func BiHubUpdateBulkTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input models.ReturnBulkTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, t := range input.Transactions {
		isSucess := logic.UpdateBalance(db, t)
		bankReceiver, _, _ := logic.ValidateBankReceiver(db, t)
		bankSender, _, _ := logic.ValidateBankSender(db, t)

		sentTransaction := models.SentTransaction{
			Transaction:  t,
			BankSender:   bankSender.BankURL,
			BankReceiver: bankReceiver.BankURL,
		}

		if isSucess {

			middleware.JkdPost("http://localhost:8084/bihub-successtransaction", sentTransaction)
		} else {

			middleware.JkdPost("http://localhost:8084/bihub-failedtransaction", sentTransaction)
		}

	}

	c.JSON(http.StatusAccepted, gin.H{"Status": "OK"})

	middleware.JkdPost("http://localhost:8084/bi-fast-esb/success-qt-processbulktransaction", input)

	return

}