package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	TransactionHash     string  `json:"transaction_hash",rethinkdb:"id"`
	SenderName          string  `json:"sender_name"`
	SenderAccountNumber string  `json:"sender_account_number"`
	SenderCity          string  `json:"sender_city"`
	SenderBankCode      string  `json:"sender_bank_code"`
	ReceiverName        string  `json:"receiver_name"`
	ReceiverBank        string  `json:"receiver_bank"`
	ReceiverBankCode    string  `json:"receiver_bank_code"`
	ReceiverAccount     string  `json:"receiver_account"`
	ReceiverAddress     string  `json:"receiver_address"`
	ReceiverCity        string  `json:"receiver_city"`
	Amount              float64 `json:"amount"`
	TransactionType     string  `json:"transaction_type"`
	Note                string  `json:"note"`
	ReceiverType        string  `json:"receiver_type"`
	ReceiverNationality string  `json:"receiver_nationality"`
	SenderType          string  `json:"sender_type"`
	SenderNationality   string  `json:"sender_nationality"`
	SenderNIK           string  `json:"sender_nik"`
	ReceiverNIK         string  `json:"receiver_nik"`
	ReceiverContact     string  `json:"receiver_contact"`
}

type BulkTransaction struct {
	gorm.Model
	BulkTransactionId   string        `json:"bulk_transaction_id"`
	SenderBankCode      string        `json:"sender_bank_code"`
	ReceiverBankCode    string        `json:"receiver_bank_code"`
	SenderAccountNumber string        `json:"sender_account_number"`
	Transactions        []Transaction `json:"transactions"`
}

type ReturnBulkTransaction struct {
	gorm.Model
	BulkTransactionId string             `json:"bulk_transaction_id"`
	BankSender        string             `json:"bank_sender"`
	BankReceiver      string             `json:"bank_receiver"`
	Transactions      []Transaction      `json:"transactions"`
	FraudTransaction  []FraudTransaction `json:"fraud_transactions"`
}

type FraudTransaction struct {
	TransactionID string  `json:"transactionid"`
	IndexFraud    float64 `json:"indexfraud"`
	FraudCategory bool    `json:"fraudcategory"`
}

type ResultTransaction struct {
	FraudTransaction FraudTransaction `json:"fraud_transaction"`
	Transaction      Transaction      `json:"transaction"`
}

type QueryInformationTransaction struct {
	ID          string      `gorethink:"id,omitempty"`
	Transaction Transaction `json:"transaction"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	DeletedAt   time.Time   `json:"deleted_at"`
}

type QueryInformationBulkTransaction struct {
	ID              string          `gorethink:"id,omitempty"`
	BulkTransaction BulkTransaction `json:"bulk_transaction"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
	DeletedAt       time.Time       `json:"deleted_at"`
}

type SendInformationTransaction struct {
	ID          string      `gorethink:"id,omitempty"`
	Transaction Transaction `json:"transaction"`
	Status      string      `json:"status"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	DeletedAt   time.Time   `json:"deleted_at"`
}

type SendInformationBulkTransaction struct {
	ID               string             `gorethink:"id,omitempty"`
	FraudTransaction []FraudTransaction `json:"fraud_transactions"`
	Status           string             `json:"status"`
	CreatedAt        string             `json:"created_at"`
	UpdatedAt        string             `json:"updated_at"`
	DeletedAt        time.Time          `json:"deleted_at"`
}

type RetrieveTransaction struct {
	ID          string      `gorethink:"id,omitempty"`
	Transaction Transaction `json:"transaction"`
	Status      string      `json:"status"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	DeletedAt   time.Time   `json:"deleted_at"`
}

type ProcessTransaction struct {
	Transaction      Transaction      `json:"transaction"`
	FraudTransaction FraudTransaction `json:"fraud_transaction"`
	BankSender       string           `json:"bank_sender"`
}
