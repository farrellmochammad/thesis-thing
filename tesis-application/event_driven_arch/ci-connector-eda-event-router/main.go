package main

import (
	"context"
	"encoding/json"
	"flag"
	"strings"

	redis "github.com/go-redis/redis/v8"

	"ci-connector-eda-event-router/middleware"
	"ci-connector-eda-event-router/models"
)

func main() {
	redis_connstring := flag.String("redis", "localhost:6379", "the port to listen on")
	channels := flag.String("channels", "get-information,ci-connector-update-transaction", "comma-separated list of channels to subscribe to")
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

		if msg.Channel == "get-information" {
			var transaction models.Transaction
			err = json.Unmarshal([]byte(msg.Payload), &transaction)
			if err != nil {
				// handle error
			}

			middleware.JkdPost("http://localhost:8084/processtransaction", transaction)
		} else if msg.Channel == "ci-connector-update-transaction" {
			var transaction models.Transaction
			err = json.Unmarshal([]byte(msg.Payload), &transaction)
			if err != nil {
				// handle error
			}

			middleware.JkdPut("http://localhost:8084/updatetransaction", transaction)
		}

	}

	// subscriber := redisClient.Subscribe(ctx, "get-information")
	// defer subscriber.Close()
	// for {
	// 	msg, err := subscriber.ReceiveMessage(ctx)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	var transaction models.Transaction
	// 	err = json.Unmarshal([]byte(msg.Payload), &transaction)
	// 	if err != nil {
	// 		// handle error
	// 	}

	// 	middleware.JkdPost("http://localhost:8084/processtransaction", transaction)
	// }

}
