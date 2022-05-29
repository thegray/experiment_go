package main

import (
	"experiment_go/kafka/sarama_consumer1/api"
	"experiment_go/kafka/sarama_consumer1/internal/conf"
	"experiment_go/kafka/sarama_consumer1/internal/pkg/transport"
	"experiment_go/kafka/sarama_consumer1/internal/repo"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port  = flag.String("port", "", "the port")
	topic = flag.String("topic", "", "topic name")
	group = flag.String("group", "", "group name")
)

func main() {
	flag.Parse()

	err := conf.InitServiceConfig("workerapp")
	if err != nil {
		panic(err)
	}

	appCfg := conf.GetGlobalConfig()

	if *port != "" {
		appCfg.ServerConfig.Port = *port
	}

	if *topic != "" {
		s := make([]string, 1)
		s[0] = *topic
		appCfg.ConsumerConfig.Topics = s
	}

	if *group != "" {
		appCfg.ConsumerConfig.Group = *group
	}

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
