package main

import (
	"context"
	"encoding/json"
	"flag"
	"strings"

	"bi-fast-eda-event-router/middleware"
	"bi-fast-eda-event-router/models"

	redis "github.com/go-redis/redis/v8"
)

func main() {
	redis_connstring := flag.String("redis", "localhost:6379", "the port to listen on")
	channels := flag.String("channels", "bi-fast-outgoing-transaction,bi-fast-success-transaction,bi-fast-failed-transaction", "comma-separated list of channels to subscribe to")
	ctx := context.Background()

	flag.Parse()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     *redis_connstring,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	})
	defer redisClient.Close()

	subscriber := redisClient.Subscribe(ctx, strings.Split(*channels, ",")...)
	defer subscriber.Close()

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		if msg.Channel == "bi-fast-outgoing-transaction" {
			var transaction models.FraudTransaction
			err = json.Unmarshal([]byte(msg.Payload), &transaction)
			if err != nil {
				// handle error
			}

			middleware.JkdPost("http://localhost:8083/validatetransaction", transaction)
		} else if msg.Channel == "bi-fast-success-transaction" {
			var transaction models.ResultTransaction
			err = json.Unmarshal([]byte(msg.Payload), &transaction)
			if err != nil {
				// handle error
			}

			middleware.JkdPost(transaction.BankReceiverUrl+"/successtransaction", transaction.Transaction)
			middleware.JkdPost(transaction.BankSenderUrl+"/successtransaction", transaction.Transaction)
		} else if msg.Channel == "bi-fast-failed-transaction" {
			var transaction models.ResultTransaction
			err = json.Unmarshal([]byte(msg.Payload), &transaction)
			if err != nil {
				// handle error
			}

			middleware.JkdPost(transaction.BankReceiverUrl+"/failedtransaction", transaction.Transaction)
			middleware.JkdPost(transaction.BankSenderUrl+"/failedtransaction", transaction.Transaction)
		}

	}

}
