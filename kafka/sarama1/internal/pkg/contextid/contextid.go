package contextid

import (
	"context"

	"github.com/google/uuid"
)

type key string

const (
	ctxID key = "context_id"
)

func New(ctx context.Context) context.Context {
	id := uuid.New().String()
	return NewWithValue(ctx, id)
}

func NewWithValue(ctx context.Context, value string) context.Context {
	if value == "" {
		value = uuid.New().String()
	}
	return context.WithValue(ctx, ctxID, value)
}

func Value(ctx context.Context) string {
	id, ok := ctx.Value(ctxID).(string)
	if !ok {
		return ""
	}

	return id
}
