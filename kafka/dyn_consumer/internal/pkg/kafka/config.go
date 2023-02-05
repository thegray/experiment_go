package kafka

import (
	"experiment_go/kafka/dyn_consumer/internal/pkg/errors"
	"fmt"
)

type ConsumerConfig struct {
	Strategy string `json:"strategy"`
	Oldest   bool   `json:"oldest"`
	Env      string `json:"env"`
	Group    string `json:"group"`
	Topic    string `json:"topic"`

	BrokersHost []string `json:"brokers"`
	User        string   `json:"user"`
	Pass        string   `json:"pass"`
	Version     string   `json:"version"`
	IsSpace     bool     `json:"is_space"`
}

func (cfg ConsumerConfig) Validate() errors.ServiceError {
	if len(cfg.BrokersHost) == 0 {
		return errors.ErrInvalidFieldFormat("Specify at least 1 host", nil)
	}

	return nil
}

func (consumer ConsumerConfig) GetGroup() string {
	return fmt.Sprintf("%s_%s", consumer.Env, consumer.Group)
}

// unused
func (consumer ConsumerConfig) GetTopic(topic string) string {
	return fmt.Sprintf("%s.%s", consumer.Env, topic)
}
