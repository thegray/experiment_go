package main

import (
	"context"
	"kafka_exp/consumer_example/consumer"
	"log"
)

func main() {
	// create a new context
	ctx := context.Background()

	log.Println("Starting consumer...")
	consumer.Consume(ctx)
}
