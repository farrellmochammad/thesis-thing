package main

import (
	"ci-connector-esb/controllers"
	"ci-connector-esb/logger"
	"ci-connector-esb/models"

	"flag"
	"log"

	"github.com/joho/godotenv"
	rethink "gopkg.in/gorethink/gorethink.v4"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.String("port", ":8083", "the port to listen on")
	analytic_url := flag.String("analyticurl", "http://localhost:8088", "the port to listen for analytic service")
	rethink_port := flag.String("rethink", "localhost:28015", "the port to listen on")
	logfile := flag.String("logfile", "analytic", "the port to listen on")
	flag.Parse()

	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger := logger.MyLogger{}

	err = logger.Init(*logfile + ".log")
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Close()

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

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("rdb", session)
		c.Set("analytic_url", *analytic_url)
		c.Set("logger", logger)
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
	r.POST("/successtransactionconfirmation", controllers.SuccessTransactionConfirmation)
	r.POST("/failtransactionconfirmation", controllers.FailTransactionConfirmation)
	r.POST("/successbulktransaction", controllers.SuccessBulkTransaction)
	r.POST("/failedtransaction", controllers.FailedTransaction)
	r.POST("/retrievetransaction", controllers.RetrieveTransaction)
	r.POST("/querybulktransactionconfirmation", controllers.QueryBulkTransactionConfirmation)

	//For Bulk Payment
	r.POST("/bulktransaction", controllers.CreateBulkTransaction)
	r.POST("/validatebulktransaction", controllers.ValidateBulkTransaction)

	// Run the Gin server in a separate goroutine
	go func() {
		if err := r.Run(*port); err != nil {
			panic(err)
		}
	}()

	// Perform other tasks or operations concurrently

	// Keep the main goroutine alive
	select {}

}
