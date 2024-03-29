package controllers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"ci-connector-eda/middleware"
	"ci-connector-eda/models"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	r "gopkg.in/gorethink/gorethink.v4"
)

func CreateTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	redis := c.MustGet("redis").(*redis.Client)
	context := c.MustGet("context").(context.Context)
	session := c.MustGet("rdb").(*r.Session)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)
	}

	// Generate a hash of the random byte slice
	hash := sha256.Sum256(randBytes)

	transaction := models.Transaction{
		TransactionHash:     hex.EncodeToString(hash[:]),
		SenderName:          input.SenderName,
		SenderAccountNumber: input.SenderAccountNumber,
		SenderCity:          input.SenderCity,
		SenderBankCode:      input.SenderBankCode,
		ReceiverName:        input.ReceiverName,
		ReceiverBank:        input.ReceiverBank,
		ReceiverBankCode:    input.ReceiverBankCode,
		ReceiverAccount:     input.ReceiverAccount,
		ReceiverAddress:     input.ReceiverAddress,
		ReceiverCity:        input.ReceiverCity,
		Amount:              input.Amount,
		TransactionType:     input.TransactionType,
		Note:                input.Note,
		ReceiverType:        input.ReceiverType,
		ReceiverNationality: input.ReceiverNationality,
		SenderType:          input.SenderType,
		SenderNationality:   input.SenderNationality,
		SenderNIK:           input.SenderNIK,
		ReceiverNIK:         input.ReceiverNIK,
		ReceiverContact:     input.ReceiverContact,
	}
	result := db.Create(&transaction)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	query_information_transaction := models.QueryInformationTransaction{
		ID:          transaction.TransactionHash,
		Transaction: transaction,
		CreatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err = r.DB("ci-connector-transaction").Table("query_information_transaction").Insert(query_information_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	channel := "get-information"
	message, err := json.Marshal(transaction)
	if err != nil {
		panic(err)
	}

	resultredis, err := redis.Publish(context, channel, message).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Published message '%s' to channel '%s' and %d clients received it.\n", message, channel, resultredis)

	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

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

func SuccessTransaction(c *gin.Context) {

	session := c.MustGet("rdb").(*r.Session)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := r.DB("ci-connector-transaction").Table("send_information_transaction").Get(input.TransactionHash).Update(map[string]interface{}{
		"UpdatedAt": time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
		"Status":    "Success",
	}).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func FailedTransaction(c *gin.Context) {
	session := c.MustGet("rdb").(*r.Session)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := r.DB("ci-connector-transaction").Table("send_information_transaction").Get(input.TransactionHash).Update(map[string]interface{}{
		"UpdatedAt": time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
		"Status":    "Failed",
	}).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func ValidateTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	session := c.MustGet("rdb").(*r.Session)
	redis := c.MustGet("redis").(*redis.Client)
	context := c.MustGet("context").(context.Context)

	var input models.FraudTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var transaction models.Transaction
	result := db.Where("transaction_hash = ?", input.TransactionID).First(&transaction)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	_, err := r.DB("ci-connector-transaction").Table("query_information_transaction").Get(input.TransactionID).Update(map[string]interface{}{
		"UpdatedAt": time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	send_information_transaction := models.SendInformationTransaction{
		ID:          input.TransactionID,
		Transaction: transaction,
		CreatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err = r.DB("ci-connector-transaction").Table("send_information_transaction").Insert(send_information_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	channel := "ci-connector-update-transaction"
	message, err := json.Marshal(transaction)
	if err != nil {
		panic(err)
	}

	resultredis, err := redis.Publish(context, channel, message).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Published message '%s' to channel '%s' and %d clients received it.\n", message, channel, resultredis)

}

func StatusTransaction(c *gin.Context) {
	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"Message": "Ok"})
	return
}

func CheckOutcomeTransaction(c *gin.Context) {
	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.JkdPost("http://localhost:8084/processtransaction", input)
}

func CheckIncomeTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account models.BankAccount
	result := db.Where("account_number = ?", input.ReceiverAccount).First(&account)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// handle record not found error
			c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found"})
		} else {
			// handle other errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	if account.Balance-input.Amount < 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Not allowed"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}
