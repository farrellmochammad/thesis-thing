package main

import (
	"ci-connector-eda/controllers"
	"ci-connector-eda/models"
	"flag"
	"log"

	"github.com/joho/godotenv"
	rethink "gopkg.in/gorethink/gorethink.v4"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.String("port", ":8083", "the port to listen on")
	analytic_url := flag.String("analyticurl", "localhost:8088", "the port to listen for analytic service")
	rethink_port := flag.String("rethink", "localhost:28015", "the port to listen on")
	flag.Parse()

	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	db, err := models.SetupModels()
	if err != nil {
		panic("Failed to connect to database!")
	}
	models.DB = db

	options := rethink.ConnectOpts{
		Address:  *rethink_port,
		Database: "ci-connector-transaction",
	}

	session, err := models.CreateSession(options)
	if err != nil {
		panic(err.Error())
	}
	defer session.Close()

	mqttClient := models.SetupMqtt("localhost", 1883)
	if err != nil {
		panic(err.Error())
	}

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("rdb", session)
		c.Set("analytic_url", *analytic_url)
		c.Set("mqtt_client", mqttClient)
		c.Next()
	})

	r.POST("/transactions", controllers.CreateTransaction)
	r.POST("/outcometransaction", controllers.CheckOutcomeTransaction)
	r.POST("/incometransaction", controllers.CheckIncomeTransaction)
	r.GET("/transactions", controllers.FindTransaction)
	r.POST("/validatetransaction", controllers.ValidateTransaction)
	r.POST("/statustransaction", controllers.StatusTransaction)

	//For SOA
	r.POST("/successtransaction", controllers.SuccessTransaction)
	r.POST("/successbulktransaction", controllers.SuccessBulkTransaction)
	r.POST("/failedtransaction", controllers.FailedTransaction)
	r.POST("/retrievetransaction", controllers.RetrieveTransaction)

	//For Bulk Payment
	r.POST("/bulktransaction", controllers.CreateBulkTransaction)
	r.POST("/validatebulktransaction", controllers.ValidateBulkTransaction)

	r.Run(*port)
}
