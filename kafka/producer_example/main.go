package main

import (
	"context"
	"kafka_exp/producer_example/producer"
	"log"
)

func main() {
	// create a new context
	ctx := context.Background()

	log.Println("Starting producer...")
	producer.Produce(ctx)
}
