package main

import (
	"ci-connector-esb/controllers"
	"ci-connector-esb/models"
	"flag"

	rethink "gopkg.in/gorethink/gorethink.v4"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.String("port", ":8083", "the port to listen on")
	rethink_port := flag.String("rethink", "localhost:28015", "the port to listen on")
	flag.Parse()
	// fmt.Printf("Listening on port %d\n", *port)

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
		c.Next()
	})

	r.POST("/transactions", controllers.CreateTransaction)
	r.POST("/outcometransaction", controllers.CheckOutcomeTransaction)
	r.POST("/incometransaction", controllers.CheckIncomeTransaction)
	r.GET("/transactions", controllers.FindTransaction)
	r.POST("/validatetransaction", controllers.ValidateTransaction)
	r.POST("/statustransaction", controllers.StatusTransaction)

	r.POST("/successtransaction", controllers.SuccessTransaction)
	r.POST("/failedtransaction", controllers.FailedTransaction)

	r.GET("/accounts", controllers.GetAccounts)
	r.GET("/accounts/:accountNumber", controllers.GetAccountByAccountNumber)
	r.POST("/account", controllers.InsertAccount)
	r.PUT("/acconts/:accountNumber", controllers.UpdateAccount)
	r.DELETE("/accounts/:accountNumber", controllers.DeleteAccount)

	r.Run(*port)
}
