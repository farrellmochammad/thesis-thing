package main

import (
	"encoding/csv"
	"fmt"
	"os"
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

type QueryInformationBulkTransactionWithDuration struct {
	QueryInformationBulkTransaction
	Duration time.Duration `json:"duration"`
}

func main() {
	QueryInformationTesting("query_information_75_tps_100_soa", "query_information_bulk_transaction")
	QueryTransactionTesting("query_transaction_75_tps_100_soa", "send_information_bulk_transaction")

}

func QueryInformationTesting(filename string, table string) {
	// Create the output CSV file
	file, err := os.Create(filename + ".csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// Write the header row
	header := []string{"ID", "CreatedAt", "UpdatedAt", "Duration"}
	err = writer.Write(header)
	if err != nil {
		panic(err)
	}

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
	cursor, err := r.Table(table).Run(session)
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

	totalDuration := time.Duration(0)

	// Print the results
	var queryinformationbulktransactionsWithDuration []QueryInformationBulkTransactionWithDuration
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

		qWithDuration := QueryInformationBulkTransactionWithDuration{
			QueryInformationBulkTransaction: q,
			Duration:                        duration,
		}
		// fmt.Printf("Time difference between CreatedAt and UpdatedAt: %v\n", duration)
		// fmt.Println(createdTime)
		// fmt.Println(updatedTime)

		queryinformationbulktransactionsWithDuration = append(queryinformationbulktransactionsWithDuration, qWithDuration)
		// ...
		totalDuration += duration

		row := []string{
			q.ID,
			q.CreatedAt,
			q.UpdatedAt,
			fmt.Sprintf("%v", duration),
		}
		err = writer.Write(row)
		if err != nil {
			panic(err)
		}
	}

	// Print the total duration in seconds
	totalSeconds := totalDuration.Seconds()
	fmt.Printf("Total duration in seconds: %.2f\n", totalSeconds)

	// Flush the writer to write any buffered data to the file
	writer.Flush()

	// Check for any errors during the flush
	if err := writer.Error(); err != nil {
		panic(err)
	}

	fmt.Println("CSV file generated successfully.")
}

func QueryTransactionTesting(filename string, table string) {
	// Create the output CSV file
	file, err := os.Create(filename + ".csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// Write the header row
	header := []string{"ID", "CreatedAt", "UpdatedAt", "Duration"}
	err = writer.Write(header)
	if err != nil {
		panic(err)
	}

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
	cursor, err := r.Table(table).Run(session)
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

	totalDuration := time.Duration(0)

	// Print the results
	var queryinformationbulktransactionsWithDuration []QueryInformationBulkTransactionWithDuration
	for _, q := range queryinformationbulktransactions {
		// do something with q
		// Calculate time difference between CreatedAt and UpdatedAt
		layout := "2006-01-02T15:04:05.999999Z07:00"

		fmt.Println("ID : ", q.ID)
		createdTime, err := time.Parse(layout, q.CreatedAt)
		if err != nil {
			panic(err)
		}

		updatedTime, err := time.Parse(layout, q.UpdatedAt)
		if err != nil {
			panic(err)
		}

		duration := updatedTime.Sub(createdTime)

		qWithDuration := QueryInformationBulkTransactionWithDuration{
			QueryInformationBulkTransaction: q,
			Duration:                        duration,
		}
		// fmt.Printf("Time difference between CreatedAt and UpdatedAt: %v\n", duration)
		// fmt.Println(createdTime)
		// fmt.Println(updatedTime)

		queryinformationbulktransactionsWithDuration = append(queryinformationbulktransactionsWithDuration, qWithDuration)
		// ...
		totalDuration += duration

		row := []string{
			q.ID,
			q.CreatedAt,
			q.UpdatedAt,
			fmt.Sprintf("%v", duration),
		}
		err = writer.Write(row)
		if err != nil {
			panic(err)
		}
	}

	// Print the total duration in seconds
	totalSeconds := totalDuration.Seconds()
	fmt.Printf("Total duration in seconds: %.2f\n", totalSeconds)

	// Flush the writer to write any buffered data to the file
	writer.Flush()

	// Check for any errors during the flush
	if err := writer.Error(); err != nil {
		panic(err)
	}

	fmt.Println("CSV file generated successfully.")
}
