package main

import (
	"experiment_go/kafka/dyn_consumer/internal/conf"
	"experiment_go/kafka/dyn_consumer/internal/data/model"
	"experiment_go/kafka/dyn_consumer/internal/pkg/transport"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"experiment_go/kafka/dyn_consumer/api"

	"experiment_go/kafka/dyn_consumer/internal/pkg/logger"
)

// build vars
var (
	GitCommit  = "unknown"
	AppVersion = "unknown"
)

func main() {
	buildInfo := map[string]string{
		"commit": GitCommit,
	}

	// When running app, use -configpath={path to config}
	// Example: ./app -configpath=../conf/config.dev.yaml
	cfgPath := flag.String("configpath", "./conf/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := conf.Init(*cfgPath)
	if err != nil {
		panic(fmt.Errorf("error parsing config. %w", err))
	}

	newLogger, err := logger.NewFromConfig(model.LoggerConfig(*cfg))
	if err != nil {
		panic(fmt.Errorf("failed to create logger"))
	}

	logger.SetDefaultLogger(*newLogger)

	logger.DefaultLogger = logger.DefaultLogger.
		With("buildinfo", buildInfo).
		With("country", "ID").
		With("version", AppVersion).
		With("service", "devtools")

	logger.Info(fmt.Sprintf("[AppConfig] app environment is set to [ %v ]", conf.ENV()))

	srv := transport.NewServer(cfg.Server)

	api.HealthCheck(srv.Engine())
	// api.Metrics(srv.Engine())

	api.KafkaWatcher(srv.Engine())

	stopFn := transport.TransportController(srv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	sig := <-quit
	logger.Info(fmt.Sprintf("exiting. received signal: %s", sig.String()))

	stopFn(time.Duration(30) * time.Second)
}
