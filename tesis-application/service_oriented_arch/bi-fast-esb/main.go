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
	r.POST("/bi-fast-validatetransaction", controllers.ValidateTransaction)
	r.POST("/prm-processtransaction", controllers.PrmProcessTransaction)
	r.PUT("/updatetransaction", controllers.UpdateTransaction)
	r.POST("/bihub-successtransaction", controllers.BiHubSuccessTransaction)
	r.POST("/bihub-failedtransaction", controllers.BiHubFailedTransaction)

	r.Run(":8084")
}
