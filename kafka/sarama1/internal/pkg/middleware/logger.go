package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	latencyMetrics = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "mis_request_latency",
		},
		[]string{"endpoint"},
	)
	requestCountMetrics = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mis_request_count",
		},
		[]string{"endpoint", "code"},
	)
)

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			now := time.Now()
			// ctx := GetContext(c)

			err := next(c)
			d := time.Since(now)
			elapsed := float64(d / time.Millisecond)

			latencyMetrics.WithLabelValues(c.Request().URL.String()).Observe(elapsed)
			requestCountMetrics.WithLabelValues(c.Request().URL.String(), strconv.Itoa(c.Response().Status)).Inc()

			// if c.Response().Status >= 500 {
			// 	log.Printf("context_id %v", contextid.Value(ctx))
			// 	logger.Error("[HTTP] Response",
			// 		"context_id", contextid.Value(ctx),
			// 		"url_path", c.Request().URL.String(),
			// 		"method", c.Request().Method,
			// 		"response_code", c.Response().Status,
			// 		"response_time", d,
			// 		"response_time_millis", d.Milliseconds(),
			// 	)
			// } else {
			// 	log.Printf("context_id %v", contextid.Value(ctx))
			// 	logger.Info("[HTTP] Response",
			// 		"context_id", contextid.Value(ctx),
			// 		"url_path", c.Request().URL.String(),
			// 		"method", c.Request().Method,
			// 		"response_code", c.Response().Status,
			// 		"response_time", d,
			// 		"response_time_millis", d.Milliseconds(),
			// 	)
			// }

			return err
		}
	}
}

func BodyDumpHandler(c echo.Context, req []byte, resp []byte) {
	// ctx := GetContext(c)
	// log.Printf("context_id %v", contextid.Value(ctx))

	// logger.Info("[HTTP] Body Dump",
	// 	"context_id", contextid.Value(ctx),
	// 	"url_path", c.Request().URL.String(),
	// 	"method", c.Request().Method,
	// 	"request_body", string(req),
	// 	"response_code", c.Response().Status,
	// 	"response_body", string(resp),
	// )
}

func BodyDumpSkipper(c echo.Context) bool {
	switch {
	// Skip healthcheck
	case c.Request().URL.String() == "/":
		fallthrough
	case strings.HasPrefix(c.Request().URL.String(), "/?"):
		fallthrough
	// Skip prometheus metrics
	case strings.HasPrefix(c.Request().URL.String(), "/metrics"):
		return true
	}

	return false
}
