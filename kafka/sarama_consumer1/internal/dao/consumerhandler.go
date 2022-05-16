package dao

// UNUSED
import "github.com/Shopify/sarama"

type ConsumerHandlerDAO interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}
