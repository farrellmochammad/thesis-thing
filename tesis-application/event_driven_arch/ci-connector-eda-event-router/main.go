package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"ci-connector-eda-event-router/middleware"
	"ci-connector-eda-event-router/models"

	redis "github.com/go-redis/redis/v8"
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

	// Create a new context with a cancel function to stop the subscription
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new channel to receive OS signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Create a new Redis pubsub subscription
	subscriber := redisClient.Subscribe(ctx, strings.Split(*channels, ",")...)

	// Create a new channel to receive pubsub messages
	ch := subscriber.Channel()

	// Start a goroutine to handle pubsub messages
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				if msg.Channel == "get-information" {
					var transaction models.Transaction
					err := json.Unmarshal([]byte(msg.Payload), &transaction)
					if err != nil {
						// handle error
					}

					fmt.Println("Got message channel from get-information, send to /processtransaction")
					middleware.JkdPost("http://localhost:8014/processtransaction", transaction)
					break
				} else if msg.Channel == "ci-connector-update-transaction" {
					var transaction models.Transaction
					err := json.Unmarshal([]byte(msg.Payload), &transaction)
					if err != nil {
						// handle error
					}

					fmt.Println("Got message channel from ci-connector-update-transaction, send to /updatetransaction")
					middleware.JkdPut("http://localhost:8014/updatetransaction", transaction)
					break
				}
			}
		}
	}()

	// Wait for a signal to stop the subscription
	<-signals

	// Unsubscribe from the Redis pubsub subscription
	if err := subscriber.Unsubscribe(ctx); err != nil {
		log.Printf("Error unsubscribing from Redis pubsub: %v", err)
	}

	// Close the Redis client connection
	if err := redisClient.Close(); err != nil {
		log.Printf("Error closing Redis client: %v", err)
	}

	// subscriber := redisClient.Subscribe(ctx, strings.Split(*channels, ",")...)
	// defer subscriber.Close()
	// for {
	// 	msg, err := subscriber.ReceiveMessage(ctx)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	if msg.Channel == "get-information" {
	// 		var transaction models.Transaction
	// 		err = json.Unmarshal([]byte(msg.Payload), &transaction)
	// 		if err != nil {
	// 			// handle error
	// 		}

	// 		middleware.JkdPost("http://localhost:8084/processtransaction", transaction)
	// 	} else if msg.Channel == "ci-connector-update-transaction" {
	// 		var transaction models.Transaction
	// 		err = json.Unmarshal([]byte(msg.Payload), &transaction)
	// 		if err != nil {
	// 			// handle error
	// 		}

	// 		middleware.JkdPut("http://localhost:8084/updatetransaction", transaction)
	// 	}

	// }

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
