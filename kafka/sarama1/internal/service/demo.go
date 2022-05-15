package service

import (
	"context"
	"experiment_go/kafka/sarama1/internal/data/dto"
	"experiment_go/kafka/sarama1/internal/data/status"
	"experiment_go/kafka/sarama1/internal/pkg/kafka"
	"log"
	"net/http"
	"time"
)

type Service struct {
	sarama *kafka.SaramaProducer
}

func New(sp *kafka.SaramaProducer) Service {
	return Service{
		sarama: sp,
	}
}

func (svc Service) ProcessDemoMessage(ctx context.Context, msg string) (status.Status, error) {

	statusMessage := "msg received: " + msg
	log.Printf(statusMessage)

	msgToQue := dto.QueMessage{
		Msg:  msg,
		Time: time.Now().Unix(),
	}

	if err := svc.sarama.SendMessage(msgToQue); err != nil {
		log.Println("Cannot que message!")
		return status.Status{}, err
	}

	return status.Status{
		Message: statusMessage,
		Code:    http.StatusOK,
	}, nil
}
