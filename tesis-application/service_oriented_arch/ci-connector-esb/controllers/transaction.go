package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"ci-connector-esb/middleware"
	"ci-connector-esb/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	analytic_url := c.MustGet("analytic_url").(string)

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

	middleware.JkdPost(analytic_url+"/input-transaction-analytic", transaction)
	middleware.JkdPost(os.Getenv("BI_FAST_ESB_URL")+"/processtransaction", transaction)

	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

func CreateBulkTransaction(c *gin.Context) {
	analytic_url := c.MustGet("analytic_url").(string)

	var input models.BulkTransaction
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

	input.BulkTransactionId = hex.EncodeToString(hash[:])

	middleware.JkdPost(analytic_url+"/input-bulk-transaction-analytic", input)
	middleware.JkdPost(os.Getenv("BI_FAST_ESB_URL")+"/processbulktransaction", input)

	c.JSON(http.StatusOK, gin.H{"data": input})
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
	analytic_url := c.MustGet("analytic_url").(string)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.JkdPost(analytic_url+"/success-transaction-analytic", input)

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func SuccessBulkTransaction(c *gin.Context) {
	analytic_url := c.MustGet("analytic_url").(string)

	var input models.ReturnBulkTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.JkdPost(analytic_url+"/success-bulk-transaction-analytic", input)

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func FailedTransaction(c *gin.Context) {

	analytic_url := c.MustGet("analytic_url").(string)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.JkdPost(analytic_url+"/failed-transaction-analytic", input)

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func RetrieveTransaction(c *gin.Context) {

	analytic_url := c.MustGet("analytic_url").(string)

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.JkdPost(analytic_url+"/retrieve-transaction-analytic", input)

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func ValidateTransaction(c *gin.Context) {
	analytic_url := c.MustGet("analytic_url").(string)

	var input models.ProcessTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.JkdPost(analytic_url+"/input-transaction-incoming-analytic", input.Transaction)
	middleware.JkdPut(os.Getenv("BI_FAST_ESB_URL")+"/updatetransaction", input.Transaction)

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
	return
}

func ValidateBulkTransaction(c *gin.Context) {
	analytic_url := c.MustGet("analytic_url").(string)

	var input models.ReturnBulkTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Marshal the struct to a JSON-encoded byte slice
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	middleware.JkdPost(analytic_url+"/input-bulk-transaction-incoming-analytic", input)
	middleware.JkdPutFile(os.Getenv("BI_FAST_ESB_URL")+"/updatebulktransaction", jsonBytes)

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
	return
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

	middleware.JkdPost(os.Getenv("BI_FAST_HUB_URL")+"/processtransaction", input)
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
