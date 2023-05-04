package main

import (
	"context"
	"flag"

	"bi-fast-eda/controllers"

	redis "github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

func main() {
	r := gin.Default()

	redis_connstring := flag.String("redis", "localhost:6379", "the port to listen on")

	flag.Parse()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     *redis_connstring,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	})
	defer redisClient.Close()

	r.Use(func(c *gin.Context) {
		c.Set("redis", redisClient)
		c.Set("context", ctx)
		c.Next()
	})

	r.POST("/processtransaction", controllers.ProcessTransaction)
	r.PUT("/updatetransaction", controllers.UpdateTransaction)

	r.Run(":8014")

}
