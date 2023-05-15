package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/bxcodec/faker/v3"
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

func main() {
	rand.Seed(time.Now().UnixNano())

	baseContent := `POST http://localhost:8083/bulktransaction/
Content-Type: application/json
@/Users/farrell/Documents/kuliah/testing-script/transactions-%d.json`

	filePath := "target.list"

	transactionSize := 50

	var content string
	// Generate random transaction data for 2 transactions
	for a := 0; a < transactionSize; a++ {
		var bulkTxn BulkTransaction
		bulkTxn.BulkTransactionId = faker.UUIDHyphenated()
		bulkTxn.SenderAccountNumber = fmt.Sprintf("%d", rand.Intn(9)+1)
		bulkTxn.SenderBankCode = fmt.Sprintf("%d", rand.Intn(2)+1)

		bankCodeUse := "2"
		if bulkTxn.SenderBankCode == "2" {
			bankCodeUse = "1"
		}
		for i := 0; i < 100; i++ {
			var txn Transaction
			faker.FakeData(&txn)

			// Override sender and receiver account numbers and bank codes
			txn.SenderAccountNumber = bulkTxn.SenderAccountNumber
			txn.SenderBankCode = bulkTxn.SenderBankCode
			txn.ReceiverAccount = fmt.Sprintf("%d", rand.Intn(9)+1)
			txn.ReceiverBankCode = bankCodeUse

			bulkTxn.Transactions = append(bulkTxn.Transactions, txn)
		}

		// Convert to JSON
		jsonData, err := json.MarshalIndent(bulkTxn, "", "  ")
		if err != nil {
			panic(err)
		}

		// Write to file
		err = ioutil.WriteFile("transactions-"+fmt.Sprintf("%d", a+1)+".json", jsonData, 0644)
		if err != nil {
			panic(err)
		}

		content += fmt.Sprintf(baseContent, a+1) + "\n"
	}

	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}

	fmt.Println("File created successfully!")

}
