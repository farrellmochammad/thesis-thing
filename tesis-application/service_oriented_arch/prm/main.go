package main

import (
	"prm/controllers"
	"prm/models"

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

	r.POST("/prm-processtransaction", controllers.PrmProcessTransaction)

	r.POST("/prm-processbulktransaction", controllers.PrmProcessBulkTransaction)

	r.Run(":8086")
}
