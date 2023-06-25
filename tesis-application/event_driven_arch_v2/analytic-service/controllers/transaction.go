package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"analytic-service/logger"
	"analytic-service/models"

	"github.com/gin-gonic/gin"

	r "gopkg.in/gorethink/gorethink.v4"
)

func InputBulkTransactionAnalytic(session *r.Session, payload string, logger logger.MyLogger) {

	var input models.BulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	logger.Log("topic/incoming-analytic-bulk-transaction" + input.BulkTransactionId)

	query_information_bulk_transaction := models.QueryInformationBulkTransaction{
		ID:              input.BulkTransactionId,
		BulkTransaction: input,
		CreatedAt:       time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err := r.DB("ci-connector-transaction").Table("query_information_bulk_transaction_eda").Insert(query_information_bulk_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

}

func InputBulkTransactionUpdateAnalytic(session *r.Session, payload string, logger logger.MyLogger) {

	var input models.BulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	logger.Log("topic/query-information-bulk-transaction" + input.BulkTransactionId)

	_, err := r.DB("ci-connector-transaction").Table("query_information_bulk_transaction_eda").Get(input.BulkTransactionId).Update(map[string]interface{}{
		"UpdatedAt": time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

}

func InputBulkTransactionIncomingAnalytic(session *r.Session, payload string, logger logger.MyLogger) {

	var input models.ReturnBulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	logger.Log("topic/ci-connector-execute-transaction" + input.BulkTransactionId)

	send_information_bulk_transaction := models.SendInformationBulkTransaction{
		ID:               input.BulkTransactionId,
		FraudTransaction: input.FraudTransaction,
		CreatedAt:        time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err := r.DB("ci-connector-transaction").Table("send_information_bulk_transaction_eda").Insert(send_information_bulk_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}
}

func InputBulkTransactionFinishUpdateAnalytic(session *r.Session, payload string, logger logger.MyLogger) {
	var input models.Transaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	logger.Log("topic/ci-connector-finished-transaction" + input.TransactionHash)

	retrieve_transaction := models.RetrieveTransaction{
		ID:          input.TransactionHash,
		Transaction: input,
		Status:      "Success",
		CreatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err := r.DB("ci-connector-transaction").Table("retrieve_transactions_eda").Insert(retrieve_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	// _, err := r.DB("ci-connector-transaction").Table("send_information_bulk_transaction_eda").Get(input.BulkTransactionId).Update(map[string]interface{}{
	// 	"UpdatedAt": time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	// 	"Status":    "Success",
	// }).RunWrite(session)
	// if err != nil {
	// 	panic(err.Error())
	// }
}

func InputTransactionIncomingAnalytic(c *gin.Context) {
	session := c.MustGet("rdb").(*r.Session)

	var input models.ProcessTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	send_information_transaction := models.SendInformationTransaction{
		ID:          input.Transaction.TransactionHash,
		Transaction: input.Transaction,
		CreatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err := r.DB("ci-connector-transaction").Table("send_information_transaction").Insert(send_information_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": input.Transaction})
}

func SuccessTransactionAnalytic(c *gin.Context) {
	// session := c.MustGet("rdb").(*r.Session)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(input)

	// _, err := r.DB("ci-connector-transaction").Table("send_information_transaction").Get(input.TransactionHash).Update(map[string]interface{}{
	// 	"UpdatedAt": time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	// 	"Status":    "Success",
	// }).RunWrite(session)
	// if err != nil {
	// 	panic(err.Error())
	// }

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func SuccessBulkTransactionAnalytic(session *r.Session, payload string, logger logger.MyLogger) {

	var input models.ReturnBulkTransaction
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		panic(err)
	}

	_, err := r.DB("ci-connector-transaction").Table("send_information_bulk_transaction_eda").Get(input.BulkTransactionId).Update(map[string]interface{}{
		"UpdatedAt": time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

}

func FailedTransactionAnalytic(c *gin.Context) {
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
		CreatedAt:   time.Now().UTC().Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	_, err := r.DB("ci-connector-transaction").Table("retrieve_transactions").Insert(retrieve_transaction).RunWrite(session)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}
