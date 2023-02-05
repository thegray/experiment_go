package auth

import (
	"experiment_go/kafka/dyn_consumer/internal/data/model"
	"experiment_go/kafka/dyn_consumer/internal/pkg/errors"
)

type TokenParser interface {
	ParseToken(token string) (model.TokenPayload, errors.ServiceError)
}
