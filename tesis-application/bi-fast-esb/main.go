package main

import (
	"bi-fast-esb/controllers"
	"bi-fast-esb/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db, err := models.SetupModels()
	if err != nil {
		panic("Failed to connect to database!")
	}
	models.DB = db

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.POST("/processtransaction", controllers.ProcessTransaction)
	r.PUT("/updatetransaction", controllers.UpdateTransaction)

	r.GET("/transactions", controllers.FindTransaction)

	r.GET("/accounts", controllers.GetAccounts)
	r.GET("/accounts/:accountNumber", controllers.GetAccountByAccountNumber)
	r.POST("/account", controllers.InsertAccount)
	r.PUT("/acconts/:accountNumber", controllers.UpdateAccount)
	r.DELETE("/accounts/:accountNumber", controllers.DeleteAccount)

	r.Run(":8084")
}
