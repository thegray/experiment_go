package dto

import "experiment_go/kafka/dyn_consumer/internal/pkg/kafka"

type CreateWatcherRequest struct {
}

type CreateWatcherResponse struct {
	ID          int                  `json:"id"`
	IsRunning   bool                 `json:"is_running"`
	CreatedAt   string               `json:"created_at"`
	KafkaConfig kafka.ConsumerConfig `json:"kafka_config"`
}
