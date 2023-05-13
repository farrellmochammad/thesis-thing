package controllers

import (
	"encoding/json"

	"ci-connector-subscriber/middleware"
	"ci-connector-subscriber/models"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func ValidateBulkTransaction(mqtt_client MQTT.Client, payload string) {

	var input models.ReturnBulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	middleware.PublishMessage(mqtt_client, "topic/bi-fast-hub-execute-transaction", input)
	middleware.PublishMessage(mqtt_client, "topic/ci-connector-execute-transaction", input)

	return
}

func BulkTransactionFinished(mqtt_client MQTT.Client, payload string) {

	var input models.ReturnBulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	middleware.PublishMessage(mqtt_client, "topic/ci-connector-finished-transaction", input)

	return
}
