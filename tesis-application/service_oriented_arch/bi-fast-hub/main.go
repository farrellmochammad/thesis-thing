package main

import (
	"bi-fast-hub/controllers"
	"bi-fast-hub/models"

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

	r.PUT("/updatetransaction", controllers.UpdateTransaction)
	r.POST("bi-fast-hub/validate-transaction", controllers.BiHubValidateTransaction)

	r.Run(":8087")
}
