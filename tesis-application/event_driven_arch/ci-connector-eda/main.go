package main

import (
	"ci-connector-eda/controllers"
	"ci-connector-eda/models"
	"context"
	"flag"

	redis "github.com/go-redis/redis/v8"
	rethink "gopkg.in/gorethink/gorethink.v4"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.String("port", ":8083", "the port to listen on")
	rethink_port := flag.String("rethink", "localhost:28015", "the port to listen on")
	redis_connstring := flag.String("redis", "localhost:6379", "the port to listen on")
	ctx := context.Background()

	flag.Parse()

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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     *redis_connstring,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	})
	defer redisClient.Close()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("rdb", session)
		c.Set("redis", redisClient)
		c.Set("context", ctx)
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
