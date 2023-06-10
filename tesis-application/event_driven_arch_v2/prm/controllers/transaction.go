package controllers

import (
	"encoding/json"

	"prm/logger"
	"prm/logic"
	"prm/middleware"
	"prm/models"

	"gorm.io/gorm"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func PrmProcessBulkTransaction(mqtt_client MQTT.Client, db *gorm.DB, payload string, logger logger.MyLogger) {

	var input models.BulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	logger.Log("topic/prm-process-bulk-transaction start" + input.BulkTransactionId)

	var fraudtransactions []models.FraudTransaction
	for _, transaction := range input.Transactions {
		fraudtransactions = append(fraudtransactions, GenerateIndexFraud(transaction))
	}

	bankSender, _, _ := logic.ValidateBankSenderBulk(db, input)
	bankReceiver, _, _ := logic.ValidateBankReceiverBulk(db, input)

	returnbulktransactions := models.ReturnBulkTransaction{
		BulkTransactionId: input.BulkTransactionId,
		BankSender:        bankSender.BankURL,
		BankReceiver:      bankReceiver.BankURL,
		BankSenderCode:    bankSender.BankCode,
		Transactions:      input.Transactions,
		FraudTransaction:  fraudtransactions,
	}

	logger.Log("topic/prm-process-bulk-transaction finish" + input.BulkTransactionId)

	middleware.PublishMessage(mqtt_client, "topic/query-information-bulk-transaction-receiver"+bankReceiver.BankCode, returnbulktransactions)

	// middleware.PublishMessage(mqtt_client, "topic/query-information-bulk-transaction-confirmation"+bankSender.BankCode, returnbulktransactions)

	// middleware.PublishMessage(mqtt_client, "topic/query-information-bulk-transaction"+bankSender.BankCode, returnbulktransactions)
	return

}

func GenerateIndexFraud(t models.Transaction) models.FraudTransaction {
	indexfraud, fraudcategory := logic.FraudDetection()

	fraudtransaction := models.FraudTransaction{
		TransactionID: t.TransactionHash,
		IndexFraud:    indexfraud,
		FraudCategory: fraudcategory,
	}

	return fraudtransaction
}
