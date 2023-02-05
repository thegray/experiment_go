package kafkawatcher

import (
	"context"
	"experiment_go/kafka/dyn_consumer/internal/data/dao"
	"experiment_go/kafka/dyn_consumer/internal/data/model"
	"experiment_go/kafka/dyn_consumer/internal/pkg/errors"
	"experiment_go/kafka/dyn_consumer/internal/pkg/kafka"
	"log"
)

type KWService interface {
	CreateNewWatcher(ctx context.Context, kc kafka.ConsumerConfig) (model.KafkaWatcher, errors.Error)

	//dummy
	DummyGenerate(ctx context.Context, key string) errors.Error
	DummyPrintWatcher(ctx context.Context) errors.Error
}

type kwService struct {
	kwRepo dao.KWDAO
}

func NewKWService(kwDAO dao.KWDAO) KWService {
	return &kwService{
		kwRepo: kwDAO,
	}
}

func (kwService *kwService) CreateNewWatcher(ctx context.Context, kc kafka.ConsumerConfig) (model.KafkaWatcher, errors.Error) {
	log.Println("---- service level", kc)
	kwService.kwRepo.DummyInsertWatchers(ctx, kc.Topic)

	return model.KafkaWatcher{}, nil
}

func (kwService *kwService) DummyGenerate(ctx context.Context, key string) errors.Error {
	kwService.kwRepo.DummyInsertWatchers(ctx, key)
	return nil
}

func (kwService *kwService) DummyPrintWatcher(ctx context.Context) errors.Error {
	kwService.kwRepo.DummyPrintWatchers(ctx)
	return nil
}
