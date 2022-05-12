package middleware

import (
	"context"
	"net/http"

	"experiment_go/kafka/sarama1/internal/pkg/contextid"

	"github.com/labstack/echo/v4"
)

const (
	echoContextKey          = "system:context"
	responseHeaderContextID = "Context-ID"
	requestHeaderRequestID  = "X-Request-Id"
)

func ContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := newContext(c.Request())
			c.Set(echoContextKey, ctx)
			c.Response().Header().Set(responseHeaderContextID, contextid.Value(ctx))
			return next(c)
		}
	}
}

func GetContext(c echo.Context) context.Context {
	i := c.Get(echoContextKey)
	if i == nil {
		ctx := newContext(c.Request())
		c.Set(echoContextKey, ctx)
		return ctx
	}

	switch v := i.(type) {
	case context.Context:
		return v
	default:
		ctx := newContext(c.Request())
		c.Set(echoContextKey, ctx)
		return ctx
	}
}

func newContext(req *http.Request) context.Context {
	if requestID := req.Header.Get(requestHeaderRequestID); requestID != "" {
		return contextid.NewWithValue(context.Background(), requestID)
	} else {
		return contextid.New(context.Background())
	}
}
