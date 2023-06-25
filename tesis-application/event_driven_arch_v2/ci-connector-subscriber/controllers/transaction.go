package controllers

import (
	"encoding/json"

	"ci-connector-subscriber/logger"
	"ci-connector-subscriber/middleware"
	"ci-connector-subscriber/models"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func ValidateBulkTransaction(mqtt_client_hub MQTT.Client, mqtt_client MQTT.Client, payload string, bankcode string, logger logger.MyLogger) {

	var input models.ReturnBulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	logger.Log("topic/query-information-bulk-transaction retrieve : " + bankcode + " - " + input.BulkTransactionId)

	middleware.PublishMessage(mqtt_client_hub, "topic/bi-fast-hub-execute-transaction", input)
	middleware.PublishMessage(mqtt_client, "topic/ci-connector-execute-transaction"+bankcode, input)

	return
}

func BulkTransactionFinished(mqtt_client MQTT.Client, payload string, logger logger.MyLogger) {

	var input models.ReturnBulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	logger.Log("topic/bi-fast-hub-execute-bulk-transaction-finish " + " - " + input.BulkTransactionId)

	middleware.PublishMessage(mqtt_client, "topic/ci-connector-finished-execute-transaction", input)

	return
}

func SendQueryInformationConfirmation(mqtt_client MQTT.Client, payload string, bankcode string, logger logger.MyLogger) {

	var input models.ReturnBulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	logger.Log("topic/query-information-bulk-transaction-confirmation retrieve : " + bankcode + " - " + input.BulkTransactionId)

	middleware.PublishMessage(mqtt_client, "topic/query-information-bulk-transaction"+bankcode, input)

	return
}

func SendQueryInformationCiConnectorReceiver(mqtt_client MQTT.Client, payload string, logger logger.MyLogger) {

	var input models.ReturnBulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	logger.Log("topic/query-information-bulk-transaction-receiver retrieve" + " - " + input.BulkTransactionId)

	middleware.PublishMessage(mqtt_client, "topic/query-information-bulk-transaction-confirmation"+input.BankSenderCode, input)

	middleware.PublishMessage(mqtt_client, "topic/query-information-bulk-transaction"+input.BankSenderCode, input)
	return
}

func BulkTransactionFinishedSuccess(mqtt_client MQTT.Client, payload string, logger logger.MyLogger) {

	var input models.Transaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	logger.Log("topic/bi-fast-hub-execute-transaction-finish-success" + " - " + input.TransactionHash)

	middleware.PublishMessage(mqtt_client, "topic/ci-connector-finished-transaction", input)

	return
}

func BulkTransactionFinishedFail(mqtt_client MQTT.Client, payload string, logger logger.MyLogger) {

	var input models.ReturnBulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	logger.Log("topic/bi-fast-hub-execute-transaction-finish-failed" + " - " + input.BulkTransactionId)

	middleware.PublishMessage(mqtt_client, "topic/ci-connector-finished-transaction", input)

	return
}
