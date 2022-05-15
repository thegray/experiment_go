package main

import (
	"experiment_go/kafka/sarama1/api"
	"experiment_go/kafka/sarama1/internal/conf"
	"experiment_go/kafka/sarama1/internal/pkg/kafka"
	"experiment_go/kafka/sarama1/internal/pkg/transport"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.sipp-now.com/spid/logger"
)

var (
	appName = flag.String("appname", "demoapp", "app name used to change config")
	addr    = flag.String("addr", ":8080", "The address to bind to")
	brokers = flag.String("brokers", "localhost:9092", "The Kafka brokers to connect to, as a comma separated list")
	verbose = flag.Bool("verbose", false, "Turn on Sarama logging")
)

func main() {
	flag.Parse()

	if *verbose {
		// sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	if *brokers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := conf.InitServiceConfig(*appName)
	if err != nil {
		panic(err)
	}

	serverCfg := conf.GetGlobalConfig().ServerConfig
	saramaCfg := conf.GetGlobalConfig().SaramaConfig

	saramaProd := kafka.NewProducer(saramaCfg)

	httpServer := transport.NewServer(serverCfg)
	api.HealthCheck(httpServer.Engine())
	api.Demo(httpServer.Engine(), saramaProd)

	stopFn := transport.TransportController(httpServer)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	sig := <-quit
	logger.Info(fmt.Sprintf("exiting. received signal: %s", sig.String()))

	stopFn(time.Duration(30) * time.Second)
	saramaProd.Close()
}
