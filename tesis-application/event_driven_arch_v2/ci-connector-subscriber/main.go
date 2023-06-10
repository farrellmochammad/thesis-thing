package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"ci-connector-subscriber/controllers"
	"ci-connector-subscriber/logger"
	"ci-connector-subscriber/models"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {

	bank_code := flag.String("bank_code", "1", "the port to listen on")
	flag.Parse()

	logger := logger.MyLogger{}

	err := logger.Init("ci-connector-subsriber.log")
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Close()

	// create a new MQTT client
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("mqtt-ci-connector-subscriber-bankcode-" + *bank_code)
	opts.SetUsername("emqx")
	opts.SetPassword("public")

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	mqttClientAnalytic := models.SetupMqtt("localhost", 1884, *bank_code)

	mqttClientTransaction := models.SetupMqtt("localhost", 1883, *bank_code)

	// set up a callback function for when a message is received
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {

		if msg.Topic() == "topic/query-information-bulk-transaction"+*bank_code {
			fmt.Println("Message incoming ", msg.Topic())
			controllers.ValidateBulkTransaction(mqttClientTransaction, mqttClientAnalytic, string(msg.Payload()), *bank_code, logger)
		}

		if msg.Topic() == "topic/query-information-bulk-transaction-confirmation"+*bank_code {
			fmt.Println("Message incoming ", msg.Topic())
			controllers.SendQueryInformationConfirmation(mqttClientAnalytic, string(msg.Payload()), *bank_code, logger)
		}

		if msg.Topic() == "topic/query-information-bulk-transaction-receiver"+*bank_code {
			fmt.Println("Message incoming ", msg.Topic())
			controllers.SendQueryInformationCiConnectorReceiver(mqttClientTransaction, string(msg.Payload()), logger)
		}

		if msg.Topic() == "topic/bi-fast-hub-execute-bulk-transaction-finish"+*bank_code {
			fmt.Println("Message incoming ", msg.Topic())
			controllers.BulkTransactionFinished(mqttClientAnalytic, string(msg.Payload()), logger)
		}

	})

	// create a new MQTT client
	// opts1 := MQTT.NewClientOptions()
	// opts1.AddBroker("tcp://localhost:1885")
	// opts1.SetClientID("mqtt-ci-connector-subscriber-bankcode-" + *bank_code)
	// opts1.SetUsername("emqx")
	// opts1.SetPassword("public")

	// opts1.OnConnect = connectHandler
	// opts1.OnConnectionLost = connectLostHandler

	// // set up a callback function for when a message is received
	// opts1.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {

	// 	if msg.Topic() == "topic/bi-fast-hub-execute-bulk-transaction-finish"+*bank_code {
	// 		fmt.Println("Message incoming ", msg.Topic())
	// 		controllers.BulkTransactionFinished(mqttClientAnalytic, string(msg.Payload()))
	// 	}

	// })

	// connect to the MQTT broker
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// connect to the MQTT broker
	// client1 := MQTT.NewClient(opts1)
	// if token1 := client1.Connect(); token1.Wait() && token1.Error() != nil {
	// 	panic(token1.Error())
	// }

	// subscribe to multiple topics of interest
	topics := []string{
		"topic/query-information-bulk-transaction" + *bank_code,
		"topic/bi-fast-hub-execute-bulk-transaction-finish" + *bank_code,
		"topic/query-information-bulk-transaction-confirmation" + *bank_code,
		"topic/query-information-bulk-transaction-receiver" + *bank_code,
	}

	for _, topic := range topics {
		if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Printf("Subscribed to topic '%s'\n", topic)
	}

	// for _, topic := range topics {
	// 	if token := client1.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
	// 		panic(token.Error())
	// 	}
	// 	fmt.Printf("Subscribed to topic '%s'\n", topic)
	// }

	// listen for messages indefinitely
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt)
	for {
		select {
		case <-sigChan:
			fmt.Println("Interrupt signal received, stopping...")
			client.Disconnect(250) // wait 250 ms before disconnecting
			return
		case <-time.After(time.Second):
			// do nothing, just keep listening for messages
		}
	}
}

var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	fmt.Println("Connected")
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
