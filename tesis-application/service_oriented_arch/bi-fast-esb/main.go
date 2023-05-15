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

	r.POST("/processbulktransaction", controllers.BulkProcessTransaction)

	r.POST("/bi-fast-validatetransaction", controllers.ValidateTransaction)
	// r.POST("/prm-processtransaction", controllers.PrmProcessTransaction)
	r.PUT("/updatetransaction", controllers.UpdateTransaction)
	r.PUT("/updatebulktransaction", controllers.UpdateBulkTransaction)
	r.POST("/bihub-successtransaction", controllers.BiHubSuccessTransaction)
	r.POST("/bihub-failedtransaction", controllers.BiHubFailedTransaction)

	//New Code After
	r.POST("/bi-fast-esb/prm-processtransaction", controllers.PrmProcessTransaction)
	r.POST("/bi-fast-esb/failed-processtransaction", controllers.FailedProcessTransaction)
	r.POST("/bi-fast-esb/report-prm-processtransaction", controllers.ReportPrmProcessTransaction)

	//Code for Bulk transaction
	r.POST("/bi-fast-esb/prm-processbulktransaction", controllers.PrmProcessBulkTransaction)
	r.POST("/bi-fast-esb/failed-processbulktransaction", controllers.FailedProcessBulkTransaction)
	r.POST("/bi-fast-esb/report-prm-processbulktransaction", controllers.ReportPrmProcessBulkTransaction)
	r.POST("/bi-fast-esb/success-qt-processbulktransaction", controllers.SuccessQtBulkTransaction)
	r.POST("/bi-fast-esb/failed-qt-processbulktransaction", controllers.FailedQtBulkTransaction)

	// Run the Gin server in a separate goroutine
	go func() {
		if err := r.Run(":8084"); err != nil {
			panic(err)
		}
	}()

	// Perform other tasks or operations concurrently

	// Keep the main goroutine alive
	select {}

}
