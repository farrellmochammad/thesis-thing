package main

import (
	"bi-fast-hub/controllers"
	"bi-fast-hub/logger"
	"bi-fast-hub/models"
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

	r.PUT("/updatetransaction", controllers.UpdateTransaction)
	r.POST("bi-fast-hub/validate-transaction", controllers.BiHubValidateTransaction)

	r.POST("/bi-fast-hub/validate-bulk-transaction", controllers.BiHubValidateBulkTransaction)

	r.PUT("/bi-fast-hub/update-bulk-transaction", controllers.BiHubUpdateBulkTransaction)

	// Run the Gin server in a separate goroutine
	go func() {
		if err := r.Run(":8087"); err != nil {
			panic(err)
		}
	}()

	// Perform other tasks or operations concurrently

	// Keep the main goroutine alive
	select {}

}
