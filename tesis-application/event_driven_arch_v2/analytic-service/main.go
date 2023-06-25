package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"analytic-service/controllers"
	"analytic-service/logger"
	"analytic-service/models"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	rethink "gopkg.in/gorethink/gorethink.v4"
)

func main() {

	rethink_port := flag.String("rethink", "localhost:28015", "the port to listen on")
	bank_code := flag.String("bank_code", "0", "the port to listen on")
	logfile := flag.String("logfile", "analytic", "the port to listen on")
	flag.Parse()

	options := rethink.ConnectOpts{
		Address:  *rethink_port,
		Database: "ci-connector-transaction",
	}

	session, err := models.CreateSession(options)
	if err != nil {
		panic(err.Error())
	}
	defer session.Close()

	logger := logger.MyLogger{}

	err = logger.Init(*logfile + ".log")
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// create a new MQTT client
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://localhost:1884")
	opts.SetClientID("mqtt-subscriber-analytic-" + *bank_code)
	opts.SetUsername("emqx")
	opts.SetPassword("public")

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	// set up a callback function for when a message is received
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {

		fmt.Println("Incoming message ", msg.Topic())
		if msg.Topic() == "topic/incoming-analytic-bulk-transaction"+*bank_code {
			fmt.Println("Message incoming analytic", msg.Topic())
			controllers.InputBulkTransactionAnalytic(session, string(msg.Payload()), logger)
		}

		if msg.Topic() == "topic/query-information-bulk-transaction"+*bank_code {
			fmt.Println("Message incoming ", msg.Topic())
			controllers.InputBulkTransactionUpdateAnalytic(session, string(msg.Payload()), logger)
		}

		if msg.Topic() == "topic/ci-connector-execute-transaction"+*bank_code {
			fmt.Println("Message incoming ", msg.Topic())
			controllers.InputBulkTransactionIncomingAnalytic(session, string(msg.Payload()), logger)
		}

		if msg.Topic() == "topic/ci-connector-finished-transaction" {
			fmt.Println("Message incoming ", msg.Topic())
			controllers.InputBulkTransactionFinishUpdateAnalytic(session, string(msg.Payload()), logger)
		}

		if msg.Topic() == "topic/ci-connector-finished-execute-transaction" {
			fmt.Println("Message incoming ", msg.Topic())
			controllers.SuccessBulkTransactionAnalytic(session, string(msg.Payload()), logger)
		}

	})

	// connect to the MQTT broker
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// subscribe to multiple topics of interest
	topics := []string{
		"topic/incoming-analytic-bulk-transaction" + *bank_code,
		"topic/query-information-bulk-transaction" + *bank_code,
		"topic/ci-connector-execute-transaction" + *bank_code,
		"topic/ci-connector-finished-transaction",
		"topic/ci-connector-finished-execute-transaction"}

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
