package main

import (
	"analytic-service/controllers"
	"analytic-service/logger"
	"analytic-service/models"

	"flag"
	"log"

	"github.com/joho/godotenv"
	rethink "gopkg.in/gorethink/gorethink.v4"

	"github.com/gin-gonic/gin"
)

func main() {

	port := flag.String("port", ":8088", "the port to listen on")
	rethink_port := flag.String("rethink", "localhost:28015", "the port to listen on")
	logfile := flag.String("logfile", "analytic", "the port to listen on")
	flag.Parse()

	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

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
	defer logger.Close()

	r.Use(func(c *gin.Context) {
		c.Set("rdb", session)
		c.Set("logger", logger)
		c.Next()
	})

	r.POST("/input-transaction-analytic", controllers.InputTransactionAnalytic)
	r.POST("/input-transaction-incoming-analytic", controllers.InputTransactionIncomingAnalytic)
	r.POST("/input-bulk-transaction-update-analytic", controllers.InputBulkTransactionUpdateAnalytic)
	r.POST("/input-bulk-transaction-analytic", controllers.InputBulkTransactionAnalytic)
	r.POST("/input-bulk-transaction-incoming-analytic", controllers.InputBulkTransactionIncomingAnalytic)
	r.POST("/success-transaction-analytic", controllers.SuccessTransactionAnalytic)
	r.POST("/success-bulk-transaction-analytic", controllers.SuccessBulkTransactionAnalytic)
	r.POST("/failed-transaction-analytic", controllers.FailedTransactionAnalytic)
	r.POST("/retrieve-transaction-analytic", controllers.RetrieveTransactionAnalytic)

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
