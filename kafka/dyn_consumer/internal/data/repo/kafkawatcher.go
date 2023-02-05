package repo

import (
	"context"
	"experiment_go/kafka/dyn_consumer/internal/data/dao"
	"experiment_go/kafka/dyn_consumer/internal/data/model"
	"experiment_go/kafka/dyn_consumer/internal/pkg/errors"
	"fmt"
	"log"
)

func KWRepo() dao.KWDAO {
	localData["tes1"] = model.KafkaWatcher{}
	localData["tes2"] = model.KafkaWatcher{}
	fmt.Println("!!!! localData:", localData)
	return &kwRepo{}
}

type kwRepo struct {
}

var localData = make(map[string]model.KafkaWatcher)

func (kwRepo *kwRepo) CreateNewWatcher(ctx context.Context) (model.KafkaWatcher, errors.Error) {
	log.Println("CreateNewWatcher not implemented yet")
	return model.KafkaWatcher{}, nil
}

func (kwRepo *kwRepo) GetWatcherByTopic(ctx context.Context, topicName string) (model.KafkaWatcher, errors.Error) {
	return model.KafkaWatcher{}, nil
}

func (kwRepo *kwRepo) DummyInsertWatchers(ctx context.Context, key string) errors.Error {
	localData[key] = model.KafkaWatcher{}
	return nil
}

func (kwRepo *kwRepo) DummyPrintWatchers(ctx context.Context) errors.Error {
	fmt.Println("!!!! localData:", localData)
	return nil
}
