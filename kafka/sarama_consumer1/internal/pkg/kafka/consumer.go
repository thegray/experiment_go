package kafka

import (
	"log"

	"experiment_go/kafka/sarama_consumer1/internal/model"

	"github.com/Shopify/sarama"
)

type SaramaConsumer struct {
}

func NewConsumer(saramaCfg model.SaramaConfig, consumerCfg model.ConsumerConfig) *SaramaConsumer {
	return &SaramaConsumer{}
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

func initClient(saramaCfg model.SaramaConfig, group string, clientCfg *sarama.Config) sarama.ConsumerGroup {
	client, err := sarama.NewConsumerGroup(saramaCfg.BrokersHost, group, clientCfg)
	if err != nil {
		log.Panicf("[KAFKA] Error creating consumer group client: %v", err)
	}
	return client
}
