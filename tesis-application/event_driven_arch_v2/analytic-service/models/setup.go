package models

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	r "gopkg.in/gorethink/gorethink.v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupModels() (*gorm.DB, error) {
	dsn := "host=localhost user=root password=root dbname=bifast_db port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = DB.AutoMigrate(&Transaction{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	err = DB.AutoMigrate(&BankAccount{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	return DB, err
}

func CreateSession(options r.ConnectOpts) (*r.Session, error) {
	session, err := r.Connect(options)
	if err != nil {
		return nil, err
	}

	// Define the table name to delete
	tableName1 := "query_information_bulk_transaction_eda" // Replace with the name of the table you want to delete
	tableName2 := "send_information_bulk_transaction_eda"
	tableName3 := "retrieve_transactions_eda"

	// Delete the table
	// Check if the table exists
	tableExists, err := r.DB("ci-connector-transaction").TableList().Contains(tableName1).Run(session)

	if err != nil {
		fmt.Println("Failed to check if table exists:", err)
		return nil, err
	}

	var exists bool
	err = tableExists.One(&exists)

	if err != nil {
		fmt.Println("Failed to get table existence:", err)
		return nil, err
	}

	// Delete the table if it exists
	if exists {
		_, err = r.DB("ci-connector-transaction").TableDrop(tableName1).RunWrite(session)

		if err != nil {
			fmt.Println("Failed to delete table:", err)
			return nil, err
		}

		fmt.Println("Table deleted successfully!")
	} else {
		fmt.Println("Table does not exist.")
	}

	// Delete the table
	// Check if the table exists
	tableExists, err = r.DB("ci-connector-transaction").TableList().Contains(tableName2).Run(session)

	if err != nil {
		fmt.Println("Failed to check if table exists:", err)
		return nil, err
	}

	err = tableExists.One(&exists)

	if err != nil {
		fmt.Println("Failed to get table existence:", err)
		return nil, err
	}

	// Delete the table if it exists
	if exists {
		_, err = r.DB("ci-connector-transaction").TableDrop(tableName2).RunWrite(session)

		if err != nil {
			fmt.Println("Failed to delete table:", err)
			return nil, err
		}

		fmt.Println("Table deleted successfully!")
	} else {
		fmt.Println("Table does not exist.")
	}

	// Delete the table
	// Check if the table exists
	tableExists, err = r.DB("ci-connector-transaction").TableList().Contains(tableName3).Run(session)

	if err != nil {
		fmt.Println("Failed to check if table exists:", err)
		return nil, err
	}

	err = tableExists.One(&exists)

	if err != nil {
		fmt.Println("Failed to get table existence:", err)
		return nil, err
	}

	// Delete the table if it exists
	if exists {
		_, err = r.DB("ci-connector-transaction").TableDrop(tableName3).RunWrite(session)

		if err != nil {
			fmt.Println("Failed to delete table:", err)
			return nil, err
		}

		fmt.Println("Table deleted successfully!")
	} else {
		fmt.Println("Table does not exist.")
	}

	// Create the table
	_, err = r.DB("ci-connector-transaction").TableCreate(tableName1).RunWrite(session)

	if err != nil {
		fmt.Println("Failed to create table:", err)
		return nil, err
	}

	_, err = r.DB("ci-connector-transaction").TableCreate(tableName2).RunWrite(session)

	if err != nil {
		fmt.Println("Failed to create table:", err)
		return nil, err
	}

	_, err = r.DB("ci-connector-transaction").TableCreate(tableName3).RunWrite(session)

	if err != nil {
		fmt.Println("Failed to create table:", err)
		return nil, err
	}

	fmt.Println("Table created successfully!")

	// if err := r.DBCreate("test").Exec(session); err != nil {
	// 	return nil, err
	// }

	// if err := r.DB("ci-connector-transaction").TableCreate("transactions").Exec(session); err != nil {
	// 	return nil, err
	// }

	// if err := r.DB("ci-connector-transaction").Table("transactions").IndexCreate("sender_account_number").Exec(session); err != nil {
	// 	return nil, err
	// }

	return session, nil
}

func SetupMqtt(broker string, port int) mqtt.Client {
	// var broker = "0.tcp.ap.ngrok.io"
	// var port = 13489
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
