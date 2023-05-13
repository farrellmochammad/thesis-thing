package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"ci-connector-subscriber/controllers"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {

	// create a new MQTT client
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("mqtt-ci-connector-subscriber")
	opts.SetUsername("emqx")
	opts.SetPassword("public")

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	// set up a callback function for when a message is received
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {

		if msg.Topic() == "topic/query-information-bulk-transaction" {
			fmt.Println("Message incoming ", msg.Topic())
			controllers.ValidateBulkTransaction(client, string(msg.Payload()))
		}

		if msg.Topic() == "topic/bi-fast-hub-execute-bulk-transaction-finish" {
			fmt.Println("Message incoming ", msg.Topic())
			controllers.BulkTransactionFinished(client, string(msg.Payload()))
		}
	})

	// connect to the MQTT broker
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// subscribe to multiple topics of interest
	topics := []string{"topic/incoming-analytic-bulk-transaction", "topic/query-information-bulk-transaction", "topic/bi-fast-hub-execute-bulk-transaction-finish"}

	for _, topic := range topics {
		if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Printf("Subscribed to topic '%s'\n", topic)
	}

	// listen for messages indefinitely
	sigChan := make(chan os.Signal, 1)
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
