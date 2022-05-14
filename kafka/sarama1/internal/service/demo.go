package service

import (
	"context"
	"experiment_go/kafka/sarama1/internal/data/status"
	"experiment_go/kafka/sarama1/internal/pkg/kafka"
	"log"
	"net/http"
)

type Service struct {
}

func New(sp *kafka.SaramaProducer) Service {
	return Service{}
}

func (svc Service) ProcessDemoMessage(ctx context.Context, msg string) (status.Status, error) {

	statusMessage := "msg received: " + msg

	log.Printf(statusMessage)

	return status.Status{
		Message: statusMessage,
		Code:    http.StatusOK,
	}, nil
}
