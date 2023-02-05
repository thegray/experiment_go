package model

import "experiment_go/kafka/dyn_consumer/internal/pkg/kafka"

type KafkaWatcher struct {
	ID          int
	IsRunning   bool
	CreatedAt   uint64
	KafkaConfig kafka.ConsumerConfig
}
