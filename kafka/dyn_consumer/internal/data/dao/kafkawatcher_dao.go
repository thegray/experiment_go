//go:generate mockgen -destination=../../mocks/dao/order_dao.go git.garena.com/wallet/sipp/va/core-v2/internal/data/dao OrderDAO
package dao

import (
	"context"
	"experiment_go/kafka/dyn_consumer/internal/data/model"
	"experiment_go/kafka/dyn_consumer/internal/pkg/errors"
)

type KWDAO interface {
	CreateNewWatcher(ctx context.Context) (model.KafkaWatcher, errors.Error)
	GetWatcherByTopic(ctx context.Context, topicName string) (model.KafkaWatcher, errors.Error)

	//dummy, shall remove later
	DummyInsertWatchers(ctx context.Context, key string) errors.Error
	DummyPrintWatchers(ctx context.Context) errors.Error
}
