package main

import (
	"bi-fast-esb/controllers"
	"bi-fast-esb/logger"
	"bi-fast-esb/models"
	"flag"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	logfile := flag.String("logfile", "analytic", "the port to listen on")
	flag.Parse()

	db, err := models.SetupModels()
	if err != nil {
		panic("Failed to connect to database!")
	}
	models.DB = db

	logger := logger.MyLogger{}

	err = logger.Init(*logfile + ".log")
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Close()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("logger", logger)
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
	r.POST("/bi-fast-esb/query-transaction-confirmation", controllers.QueryTransactionConfirmation)
	r.POST("/bi-fast-esb/failed-processbulktransaction", controllers.FailedProcessBulkTransaction)
	r.POST("/bi-fast-esb/report-prm-processbulktransaction", controllers.ReportPrmProcessBulkTransaction)
	r.POST("/bi-fast-esb/success-qt-processbulktransaction", controllers.SuccessQtBulkTransaction)
	r.POST("/bi-fast-esb/failed-qt-processbulktransaction", controllers.FailedQtBulkTransaction)

	r.POST("/bi-fast-esb/success-transaction-confirmation", controllers.SuccessTransactionConfirmation)
	r.POST("/bi-fast-esb/fail-transaction-confirmation", controllers.FailTransactionConfirmation)

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
