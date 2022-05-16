package transport

import (
	"context"
	"log"
	"sync"

	"experiment_go/kafka/sarama_consumer1/internal/model"
	"experiment_go/kafka/sarama_consumer1/internal/repo"

	"github.com/Shopify/sarama"
)

type stopConsume func()

type SaramaConsumer struct {
	clientGroup     sarama.ConsumerGroup
	topics          []string
	consumerHandler repo.ConsumerHandlerImpl
	wg              *sync.WaitGroup
	stopFunc        stopConsume
}

func NewConsumer(saramaCfg model.SaramaConfig, consumerCfg model.ConsumerConfig, consumerHandler repo.ConsumerHandlerImpl) *SaramaConsumer {

	clientCfg := createConfig(consumerCfg)
	if consumerCfg.Group == "" {
		log.Panicf("[KAFKA] 'group' must not empty, specify in config!")
	}
	if len(consumerCfg.Topics) == 0 {
		log.Panicf("[KAFKA] 'topics' must not empty, specify in config!")
	}
	client := initConsumerGroup(saramaCfg, consumerCfg.Group, clientCfg)

	return &SaramaConsumer{
		clientGroup:     client,
		topics:          consumerCfg.Topics,
		consumerHandler: consumerHandler,
	}
}

func createConfig(cfg model.ConsumerConfig) *sarama.Config {
	kafkaVersion, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		log.Panicf("[KAFKA] Error parsing Kafka version: %v", err)
	}
	config := sarama.NewConfig()
	config.Version = kafkaVersion

	switch cfg.Strategy {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("[KAFKA] Unrecognized consumer group partition assignor: %s",
			cfg.Strategy,
		)
	}

	if cfg.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	return config
}

func initConsumerGroup(saramaCfg model.SaramaConfig, group string, clientCfg *sarama.Config) sarama.ConsumerGroup {
	client, err := sarama.NewConsumerGroup(saramaCfg.BrokersHost, group, clientCfg)
	if err != nil {
		log.Panicf("[KAFKA] Error creating consumer group client: %v", err)
	}
	return client
}

// Start function will be run in goroutine
func (sc *SaramaConsumer) Start() {
	log.Println("startt consumerrrrrr")
	ctx, cancel := context.WithCancel(context.Background())
	sc.wg = &sync.WaitGroup{}
	sc.wg.Add(1)
	// func() {
	go func() {
		defer sc.wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := sc.clientGroup.Consume(ctx, sc.topics, &sc.consumerHandler); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Println("Consumer:", ctx.Err())
				return
			}
			sc.consumerHandler.Ready = make(chan bool)
		}
	}()

	sc.stopFunc = func() {
		cancel()
	}
	// sc.wg.Wait()
}

func (sc *SaramaConsumer) Stop() {
	log.Println("stopping consumer..")
	sc.stopFunc()
	sc.wg.Wait()

	log.Println("closing client group safely...")
	if err := sc.clientGroup.Close(); err != nil {
		log.Panicf("Error closing kafka consumer: %v", err)
	}
}
