package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"analytic-service/logger"
	"analytic-service/models"

	"github.com/gin-gonic/gin"

	r "gopkg.in/gorethink/gorethink.v4"
)

func InputTransactionAnalytic(c *gin.Context) {
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

	query_information_transaction := models.QueryInformationTransaction{
		ID:          transaction.TransactionHash,
		Transaction: transaction,
		CreatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err = r.DB("ci-connector-transaction").Table("query_information_transaction").Insert(query_information_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": transaction})

}

func InputBulkTransactionAnalytic(c *gin.Context) {
	session := c.MustGet("rdb").(*r.Session)

	var input models.BulkTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger := c.MustGet("logger").(logger.MyLogger)
	logger.Log("/input-bulk-transaction-analytic/" + input.BulkTransactionId)

	query_information_bulk_transaction := models.QueryInformationBulkTransaction{
		ID:              input.BulkTransactionId,
		BulkTransaction: input,
		CreatedAt:       time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err := r.DB("ci-connector-transaction").Table("query_information_bulk_transaction").Insert(query_information_bulk_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": input})

}

func InputBulkTransactionUpdateAnalytic(c *gin.Context) {
	session := c.MustGet("rdb").(*r.Session)

	var input models.BulkTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger := c.MustGet("logger").(logger.MyLogger)
	logger.Log("/input-bulk-transaction-update-analytic/" + input.BulkTransactionId)

	_, err := r.DB("ci-connector-transaction").Table("query_information_bulk_transaction").Get(input.BulkTransactionId).Update(map[string]interface{}{
		"UpdatedAt": time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func InputBulkTransactionIncomingAnalytic(c *gin.Context) {
	session := c.MustGet("rdb").(*r.Session)

	var input models.ReturnBulkTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger := c.MustGet("logger").(logger.MyLogger)
	logger.Log("/input-bulk-transaction-incoming-analytic/" + input.BulkTransactionId)

	send_information_bulk_transaction := models.SendInformationBulkTransaction{
		ID:               input.BulkTransactionId,
		FraudTransaction: input.FraudTransaction,
		CreatedAt:        time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err := r.DB("ci-connector-transaction").Table("send_information_bulk_transaction").Insert(send_information_bulk_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func InputTransactionIncomingAnalytic(c *gin.Context) {
	session := c.MustGet("rdb").(*r.Session)
	logger := c.MustGet("logger").(logger.MyLogger)

	var input models.ProcessTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log("/success-bulk-transaction-analytic start" + input.Transaction.TransactionHash)

	send_information_transaction := models.SendInformationTransaction{
		ID:          input.Transaction.TransactionHash,
		Transaction: input.Transaction,
		CreatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err := r.DB("ci-connector-transaction").Table("send_information_bulk_transaction").Insert(send_information_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": input.Transaction})
}

func SuccessTransactionAnalytic(c *gin.Context) {
	session := c.MustGet("rdb").(*r.Session)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	retrieve_transaction := models.RetrieveTransaction{
		ID:          input.TransactionHash,
		Transaction: input,
		Status:      "Success",
		CreatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err := r.DB("ci-connector-transaction").Table("retrieve_transactions").Insert(retrieve_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func SuccessBulkTransactionAnalytic(c *gin.Context) {
	session := c.MustGet("rdb").(*r.Session)

	var input models.ReturnBulkTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger := c.MustGet("logger").(logger.MyLogger)
	logger.Log("/success-bulk-transaction-analytic/" + input.BulkTransactionId)

	_, err := r.DB("ci-connector-transaction").Table("send_information_bulk_transaction").Get(input.BulkTransactionId).Update(map[string]interface{}{
		"UpdatedAt": time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
		"Status":    "Success",
	}).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func FailedTransactionAnalytic(c *gin.Context) {
	session := c.MustGet("rdb").(*r.Session)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := r.DB("ci-connector-transaction").Table("send_information_bulk_transaction").Get(input.TransactionHash).Update(map[string]interface{}{
		"UpdatedAt": time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
		"Status":    "Failed",
	}).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func RetrieveTransactionAnalytic(c *gin.Context) {
	session := c.MustGet("rdb").(*r.Session)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	retrieve_transaction := models.RetrieveTransaction{
		ID:          input.TransactionHash,
		Transaction: input,
		Status:      "Success",
		CreatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err := r.DB("ci-connector-transaction").Table("retrieve_transactions").Insert(retrieve_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}
