package controllers

import (
	"net/http"

	"ci-connector-eda/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAccounts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var accounts []models.BankAccount
	db.Find(&accounts)
	c.JSON(http.StatusOK, accounts)
}

func GetAccountByAccountNumber(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var account models.BankAccount
	result := db.Where("account_number = ?", c.Param("accountNumber")).First(&account)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	c.JSON(http.StatusOK, account)
}

func InsertAccount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var account models.BankAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Create(&account)
	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Account already exists"})
		return
	}
	c.JSON(http.StatusCreated, account)
}

func UpdateAccount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var account models.BankAccount
	result := db.Where("account_number = ?", c.Param("accountNumber")).First(&account)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&account)
	c.JSON(http.StatusOK, account)
}

func DeleteAccount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var account models.BankAccount
	result := db.Where("account_number = ?", c.Param("accountNumber")).First(&account)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	db.Delete(&account)
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted"})
}
