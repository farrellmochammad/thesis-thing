package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"bi-fast-hub-eda/logic"
	"bi-fast-hub-eda/models"

	redis "github.com/go-redis/redis/v8"
)

func main() {
	redis_connstring := flag.String("redis", "localhost:6379", "the port to listen on")
	ctx := context.Background()

	flag.Parse()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     *redis_connstring,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	})
	defer redisClient.Close()

	subscriber := redisClient.Subscribe(ctx, "bi-fast-incoming-processtransaction")
	defer subscriber.Close()

	db, err := models.SetupModels()
	if err != nil {
		panic("Failed to connect to database!")
	}

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println("Message bi fast hub eda : ", msg)

		var transaction models.Transaction
		err = json.Unmarshal([]byte(msg.Payload), &transaction)
		if err != nil {
			panic(err)
		}

		isSucess := logic.UpdateBalance(db, transaction)

		if isSucess {
			bankReceiver, _, _ := logic.ValidateBankReceiver(db, transaction)
			bankSender, _, _ := logic.ValidateBankSender(db, transaction)

			resulttransaction := models.ResultTransaction{
				BankSenderUrl:   bankSender.BankURL,
				BankReceiverUrl: bankReceiver.BankURL,
				Transaction:     transaction,
			}

			jsonBytes, err := json.Marshal(resulttransaction)
			if err != nil {
				panic(err)
			}

			publisher := redisClient.Publish(ctx, "bi-fast-success-transaction", jsonBytes)
			result, err := publisher.Result()
			if err != nil {
				panic(err)
			}
			fmt.Println("Published message to channel: ", result)
			// middleware.JkdPost(bankReceiver.BankURL+"/successtransaction", input)
			// middleware.JkdPost(bankSender.BankURL+"/successtransaction", input)
		} else {
			bankReceiver, _, _ := logic.ValidateBankReceiver(db, transaction)
			bankSender, _, _ := logic.ValidateBankReceiver(db, transaction)

			resulttransaction := models.ResultTransaction{
				BankSenderUrl:   bankSender.BankURL,
				BankReceiverUrl: bankReceiver.BankURL,
				Transaction:     transaction,
			}

			jsonBytes, err := json.Marshal(resulttransaction)
			if err != nil {
				panic(err)
			}

			publisher := redisClient.Publish(ctx, "bi-fast-failed-transaction", jsonBytes)
			result, err := publisher.Result()
			if err != nil {
				panic(err)
			}
			fmt.Println("Published message to channel: ", result)
			// middleware.JkdPost(bankReceiver.BankURL+"/failedtransaction", input)
			// middleware.JkdPost(bankSender.BankURL+"/failedtransaction", input)
		}

	}

}
