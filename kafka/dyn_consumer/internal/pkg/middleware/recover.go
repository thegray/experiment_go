package middleware

import (
	"fmt"
	"runtime/debug"

	"experiment_go/kafka/dyn_consumer/internal/pkg/logger"

	echo "github.com/labstack/echo/v4"
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
					logger.ErrorCtx(ctx, fmt.Sprintf("[HTTP:Recover] panic %s", err.Error()),
						"stacktrace", string(debug.Stack()))
					c.Error(err)
				}
			}()

			return next(c)
		}
	}
}
