package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"prm-eda/logic"
	"prm-eda/models"

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

	subscriber := redisClient.Subscribe(ctx, "bi-fast-incoming-transaction")
	defer subscriber.Close()

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println("Message : ", msg)

		var transaction models.Transaction
		err = json.Unmarshal([]byte(msg.Payload), &transaction)
		if err != nil {
			panic(err)
		}

		indexfraud, fraudcategory := logic.FraudDetection()

		fraudtransaction := models.FraudTransaction{
			TransactionID: transaction.Hash,
			IndexFraud:    indexfraud,
			FraudCategory: fraudcategory,
		}

		jsonBytes, err := json.Marshal(fraudtransaction)
		if err != nil {
			panic(err)
		}

		publisher := redisClient.Publish(ctx, "bi-fast-outgoing-transaction", jsonBytes)
		result, err := publisher.Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("Published message to channel: ", result)
	}

}
