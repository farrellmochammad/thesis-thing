package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

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
