package repo

import (
	"log"

	"github.com/Shopify/sarama"
)

type ConsumerHandlerImpl struct {
	Ready chan bool
}

func NewConsumerHandler() *ConsumerHandlerImpl {
	return &ConsumerHandlerImpl{
		Ready: make(chan bool),
	}
}

func (ch *ConsumerHandlerImpl) Setup(sarama.ConsumerGroupSession) error {
	close(ch.Ready)
	return nil
}

func (ch *ConsumerHandlerImpl) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (ch *ConsumerHandlerImpl) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}
	return nil
}
