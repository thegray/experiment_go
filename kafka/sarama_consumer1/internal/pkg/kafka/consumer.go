package kafka

import (
	"context"
	"log"

	"experiment_go/kafka/sarama_consumer1/internal/model"
	"experiment_go/kafka/sarama_consumer1/internal/repo"

	"github.com/Shopify/sarama"
)

type SaramaConsumer struct {
	consumerGroup sarama.ConsumerGroup
}

func NewConsumer(ctx context.Context, saramaCfg model.SaramaConfig, consumerCfg model.ConsumerConfig, consumerHandler repo.ConsumerHandlerImpl) *SaramaConsumer {
	clientCfg := createConfig(consumerCfg)
	if consumerCfg.Group == "" {
		log.Panicf("[KAFKA] Consumer Group must not empty, specify in config!")
	}
	clientGroup := initConsumerGroup(saramaCfg, consumerCfg.Group, clientCfg)
	startConsumerGroup(ctx, clientGroup, consumerCfg.Topics, consumerHandler)

	return &SaramaConsumer{consumerGroup: clientGroup}
}

func (sc *SaramaConsumer) Close() error {
	if err := sc.consumerGroup.Close(); err != nil {
		return err
	}
	return nil
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

func startConsumerGroup(ctx context.Context, groupClient sarama.ConsumerGroup, topics []string, consumerHandler repo.ConsumerHandlerImpl) {
	// ctx, _ := context.WithCancel(context.Background())
	// wg := &sync.WaitGroup{}
	// wg.Add(1)
	// func() {
	go func() {
		// defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := groupClient.Consume(ctx, topics, &consumerHandler); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumerHandler.Ready = make(chan bool)
		}
	}()
	// wg.Wait()
}
