package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"bi-fast-eda/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	redis "github.com/go-redis/redis/v8"
)

func FindTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var transactions []models.Transaction
	result := db.Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

func ProcessTransaction(c *gin.Context) {

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redis := c.MustGet("redis").(*redis.Client)
	context := c.MustGet("context").(context.Context)

	channel := "bi-fast-incoming-transaction"
	message, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	resultredis, err := redis.Publish(context, channel, message).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Published message '%s' to channel '%s' and %d clients received it.\n", message, channel, resultredis)

}

func UpdateTransaction(c *gin.Context) {

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redis := c.MustGet("redis").(*redis.Client)
	context := c.MustGet("context").(context.Context)

	channel := "bi-fast-incoming-processtransaction"
	message, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	resultredis, err := redis.Publish(context, channel, message).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Published message '%s' to channel '%s' and %d clients received it.\n", message, channel, resultredis)

}
