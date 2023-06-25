package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bi-fast-hub/logger"
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

func BiHubValidateTransaction(mqtt_client MQTT.Client, db *gorm.DB, payload string, logger logger.MyLogger) {

	var input models.BulkTransaction
	logger.Log("topic/bi-fast-hub-incoming-bulk-transaction start " + input.BulkTransactionId)
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

	logger.Log("topic/bi-fast-hub-incoming-bulk-transaction finish" + input.BulkTransactionId)

	if isValidateAmount {
		middleware.PublishMessage(mqtt_client, "topic/prm-process-bulk-transaction", input)
		return
	}

}

func BiHubValidateBulkTransaction(mqtt_client MQTT.Client, db *gorm.DB, payload string) {

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

}

func BiHubUpdateBulkTransaction(mqtt_client MQTT.Client, db *gorm.DB, payload string, logger logger.MyLogger) {

	var input models.ReturnBulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	logger.Log("topic/bi-fast-hub-execute-transaction start" + input.BulkTransactionId)
	bankCode := "0"
	for _, t := range input.Transactions {
		isSucess := logic.UpdateBalance(db, t)
		bankReceiver, _, _ := logic.ValidateBankReceiver(db, t)
		bankSender, _, _ := logic.ValidateBankSender(db, t)

		t.TransactionHash = input.BulkTransactionId

		bankCode = bankSender.BankCode

		// sentTransaction := models.SentTransaction{
		// 	Transaction:  t,
		// 	BankSender:   bankSender.BankURL,
		// 	BankReceiver: bankReceiver.BankURL,
		// }

		if isSucess {

			middleware.PublishMessage(mqtt_client, "topic/bi-fast-hub-execute-transaction-finish-success"+bankSender.BankCode, t)
			middleware.PublishMessage(mqtt_client, "topic/bi-fast-hub-execute-transaction-finish-success"+bankReceiver.BankCode, t)
		} else {
			middleware.PublishMessage(mqtt_client, "topic/bi-fast-hub-execute-transaction-finish-failed"+bankSender.BankCode, t)
			middleware.PublishMessage(mqtt_client, "topic/bi-fast-hub-execute-transaction-finish-failed"+bankSender.BankCode, t)
		}

	}

	logger.Log("topic/bi-fast-hub-execute-transaction finish" + input.BulkTransactionId)

	middleware.PublishMessage(mqtt_client, "topic/bi-fast-hub-execute-bulk-transaction-finish"+bankCode, input)

	return

}
