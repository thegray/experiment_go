package main

import (
	"context"
	"experiment_go/kafka/sarama_consumer1/internal/conf"
	"experiment_go/kafka/sarama_consumer1/internal/pkg/kafka"
	"experiment_go/kafka/sarama_consumer1/internal/repo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	err := conf.InitServiceConfig("workerapp")
	if err != nil {
		panic(err)
	}

	appCfg := conf.GetGlobalConfig()
	kafkaConsumer := repo.NewConsumerHandler()
	kc := kafka.NewConsumer(ctx, appCfg.SaramaConfig, appCfg.ConsumerConfig, *kafkaConsumer)

	<-kafkaConsumer.Ready
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// keepRunning := true
	// for keepRunning {
	// 	select {
	// 	case <-ctx.Done():
	// 		log.Println("terminating: context cancelled")
	// 		keepRunning = false
	// 	case <-sigterm:
	// 		log.Println("terminating: via signal")
	// 		keepRunning = false
	// 		// case <-sigusr1:
	// 		// 	toggleConsumptionFlow(client, &consumptionIsPaused)
	// 	}
	// }

	defer func() {
		log.Println("try to stop everything...")
		cancel()
		if err = kc.Close(); err != nil {
			log.Panicf("Error closing kafka consumer: %v", err)
		}
	}()

	e := echo.New()
	e.GET("/", hello)
	e.Start(appCfg.ServerConfig.Port)
	// e.Logger.Fatal(e.Start(appCfg.ServerConfig.Port))
}

func hello(c echo.Context) error {
	log.Print("hello called")
	return c.String(http.StatusOK, "hello world")
}
