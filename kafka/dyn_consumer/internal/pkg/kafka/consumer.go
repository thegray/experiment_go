package kafka

// - id consumer
// - config
// - data (array), append only
// - consumer obj

import (
	"context"
	"fmt"
	"log"

	"experiment_go/kafka/dyn_consumer/internal/pkg/scram"

	"github.com/Shopify/sarama"
)

type ConsumerGroupHandler interface {
	sarama.ConsumerGroupHandler
	Ready()
}

type saramaConsumer struct {
	consumerGroup   sarama.ConsumerGroup
	topics          []string
	consumerHandler ConsumerGroupHandler
	stopFunc        func()
}

func BuildSaramaConsumer(cfg ConsumerConfig, consumerHandler ConsumerGroupHandler, topics []string) (*saramaConsumer, error) {
	consumerConfig, errSaramaCfg := saramaConsumerConfig(cfg)
	if errSaramaCfg != nil {
		return nil, errSaramaCfg
	}

	if cfg.GetGroup() == "" {
		return nil, fmt.Errorf("[BuildSaramaConsumer] 'group' must not be empty")
	}

	consumerGroup, err := sarama.NewConsumerGroup(cfg.BrokersHost, cfg.GetGroup(), consumerConfig)
	if err != nil {
		return nil, err
	}

	return &saramaConsumer{
		consumerGroup:   consumerGroup,
		topics:          topics,
		consumerHandler: consumerHandler,
	}, nil
}

func saramaConsumerConfig(cfg ConsumerConfig) (*sarama.Config, error) {
	config := sarama.NewConfig()

	if kafkaVersion, err := sarama.ParseKafkaVersion(cfg.Version); err != nil {
		return nil, err
	} else {
		config.Version = kafkaVersion
	}

	if cfg.User != "" && cfg.Pass != "" && cfg.IsSpace {
		config.Net.SASL.User = cfg.User
		config.Net.SASL.Password = cfg.Pass
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		config.Net.SASL.Handshake = true
		config.Net.SASL.Enable = true
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &scram.XDGSCRAMClient{HashGeneratorFcn: scram.ScramSHA512}
		}
	}

	switch cfg.Strategy {
	case sarama.StickyBalanceStrategyName:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case sarama.RoundRobinBalanceStrategyName:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case sarama.RangeBalanceStrategyName:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		return nil, fmt.Errorf("unrecognized consumer group partition assignor: %s", cfg.Strategy)
	}

	if cfg.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	return config, nil
}

func (sc *saramaConsumer) Start() {
	log.Println("Starting comsumer...")

	ctx, cancel := context.WithCancel(context.Background())

	sc.stopFunc = func() {
		cancel()
	}

	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := sc.consumerGroup.Consume(ctx, sc.topics, sc.consumerHandler); err != nil {
			log.Println("consume failed",
				"error", err.Error(),
			)
		}

		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return
		}
		sc.consumerHandler.Ready()
	}
}

func (sc *saramaConsumer) Stop() {
	log.Println("Stopping consumer...")

	sc.stopFunc()

	if err := sc.consumerGroup.Close(); err != nil {
		log.Println("error closing kafka consumer",
			"error", err.Error(),
		)
	}
}
