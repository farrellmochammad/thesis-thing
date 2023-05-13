package main

import (
	"fmt"
	"time"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type Transaction struct {
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
	BulkTransactionId   string        `json:"bulk_transaction_id"`
	SenderBankCode      string        `json:"sender_bank_code"`
	SenderAccountNumber string        `json:"sender_account_number"`
	Transactions        []Transaction `json:"transactions"`
}

type QueryInformationBulkTransaction struct {
	ID              string          `gorethink:"id,omitempty"`
	BulkTransaction BulkTransaction `json:"bulk_transaction"`
	CreatedAt       string          `gorethink:"CreatedAt"`
	UpdatedAt       string          `gorethink:"UpdatedAt"`
	DeletedAt       time.Time       `json:"deleted_at"`
}

func main() {
	// Connect to RethinkDB
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "ci-connector-transaction",
	})
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Get all data from the "persons" table
	cursor, err := r.Table("query_information_bulk_transaction").Run(session)
	if err != nil {
		panic(err)
	}
	defer cursor.Close()

	// Iterate over the results and decode them into Person structs
	var queryinformationbulktransactions []QueryInformationBulkTransaction
	err = cursor.All(&queryinformationbulktransactions)
	if err != nil {
		panic(err)
	}

	// Print the results

	for _, q := range queryinformationbulktransactions {
		// do something with q
		// Calculate time difference between CreatedAt and UpdatedAt
		layout := "2006-01-02T15:04:05.999999Z07:00"

		createdTime, err := time.Parse(layout, q.CreatedAt)
		if err != nil {
			panic(err)
		}

		updatedTime, err := time.Parse(layout, q.UpdatedAt)
		if err != nil {
			panic(err)
		}

		duration := updatedTime.Sub(createdTime)
		fmt.Printf("Time difference between CreatedAt and UpdatedAt: %v\n", duration)
		fmt.Println(createdTime)
		fmt.Println(updatedTime)
		// ...
	}
}
