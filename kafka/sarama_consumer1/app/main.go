package main

import (
	"experiment_go/kafka/sarama_consumer1/api"
	"experiment_go/kafka/sarama_consumer1/internal/conf"
	"experiment_go/kafka/sarama_consumer1/internal/pkg/transport"
	"experiment_go/kafka/sarama_consumer1/internal/repo"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	err := conf.InitServiceConfig("workerapp")
	if err != nil {
		panic(err)
	}

	appCfg := conf.GetGlobalConfig()

	httpServer := transport.NewServer(appCfg.ServerConfig)
	api.HealthCheck(httpServer.Engine())

	consumerHandler := repo.NewConsumerHandler()
	saramaConsumer := transport.NewConsumer(appCfg.SaramaConfig, appCfg.ConsumerConfig, *consumerHandler)

	stopFn := transport.TransportController(
		httpServer,
		saramaConsumer,
	)

	<-consumerHandler.Ready
	log.Println("Sarama consumer up and running!...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	sig := <-quit
	log.Printf("exiting. received signal: %s", sig.String())

	stopFn(time.Duration(30) * time.Second)
}
