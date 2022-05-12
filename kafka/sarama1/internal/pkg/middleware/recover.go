package middleware

import (
	"fmt"
	"runtime/debug"

	"experiment_go/kafka/sarama1/internal/pkg/contextid"

	"git.sipp-now.com/spid/logger"
	"github.com/labstack/echo/v4"
)

func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					ctx := GetContext(c)

					if !ok {
						err = fmt.Errorf("%+v", r)
					}

					logger.Error(fmt.Sprintf("[HTTP:Recover] panic %s", err.Error()),
						"context_id", contextid.Value(ctx),
						"stacktrace", string(debug.Stack()),
					)

					c.Error(err)
				}
			}()

			return next(c)
		}
	}
}
