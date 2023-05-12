package main

import (
	"analytic-service/controllers"
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

	r.Use(func(c *gin.Context) {
		c.Set("rdb", session)
		c.Next()
	})

	r.POST("/input-transaction-analytic", controllers.InputTransactionAnalytic)
	r.POST("/input-transaction-incoming-analytic", controllers.InputTransactionIncomingAnalytic)
	r.POST("/input-bulk-transaction-analytic", controllers.InputBulkTransactionAnalytic)
	r.POST("/input-bulk-transaction-incoming-analytic", controllers.InputBulkTransactionIncomingAnalytic)
	r.POST("/success-transaction-analytic", controllers.SuccessTransactionAnalytic)
	r.POST("/success-bulk-transaction-analytic", controllers.SuccessBulkTransactionAnalytic)
	r.POST("/failed-transaction-analytic", controllers.FailedTransactionAnalytic)
	r.POST("/retrieve-transaction-analytic", controllers.RetrieveTransactionAnalytic)

	r.Run(*port)
}
