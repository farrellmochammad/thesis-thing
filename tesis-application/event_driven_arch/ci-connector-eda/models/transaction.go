package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	TransactionHash     string  `json:"transaction_hash"`
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

type FraudTransaction struct {
	TransactionID string  `json:"transactionid"`
	IndexFraud    float64 `json:"indexfraud"`
	FraudCategory bool    `json:"fraudcategory"`
}

type QueryInformationTransaction struct {
	ID          string      `gorethink:"id,omitempty"`
	Transaction Transaction `json:"transaction"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	DeletedAt   time.Time   `json:"deleted_at"`
}

type SendInformationTransaction struct {
	ID          string      `gorethink:"id,omitempty"`
	Transaction Transaction `json:"transaction"`
	Status      string      `json:"status"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	DeletedAt   time.Time   `json:"deleted_at"`
}
