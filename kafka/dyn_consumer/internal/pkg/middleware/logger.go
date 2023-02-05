package middleware

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"experiment_go/kafka/dyn_consumer/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

/*
omit prometheus part for now
var (
	latencyMetrics = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "sipp_service_response_time",
			Help: "The response time of each http request that we call to other service.",
			Buckets: []float64{
				0.01, 0.05, 0.1, 0.15, 0.2, 0.25, 0.3, 0.35, 0.4, 0.45, 0.5,
				0.55, 0.6, 0.65, 0.7, 0.75, 0.8, 1.0, 1.1,
				1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9, 2.0, 2.2, 2.4,
				2.6, 2.8, 3.0, 3.3,
				3.6, 4.0, 6.0, 8.0, 10.0, 15.0, 30.0,
			},
		},
		[]string{"service", "endpoint"},
	)
	requestCountMetrics = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sipp_service_error_code_total",
			Help: "The error code from CoreV2 http request.",
		},
		[]string{"service", "endpoint", "code"},
	)
)
*/

// this logger middleware is adapted from echo BodyDump middleware
// https://github.com/labstack/echo/blob/master/middleware/body_dump.go
type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if BodyDumpSkipper(c) {
				return next(c)
			}

			now := time.Now()

			// Request
			reqBody := []byte{}
			if c.Request().Body != nil { // Read
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			// Response
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			if err = next(c); err != nil {
				c.Error(err)
			}

			d := time.Since(now)

			/* omit prometheus part for now
			elapsed := float64(d / time.Millisecond)

			pathRegExp := regexp.MustCompile(`/[0-9]+$`)
			path := pathRegExp.ReplaceAllString(c.Request().URL.Path, "")

			latencyMetrics.WithLabelValues(constant.MetricsService, path).Observe(elapsed)
			requestCountMetrics.WithLabelValues(constant.MetricsService, path, strconv.Itoa(c.Response().Status)).Inc()
			*/

			ctx := GetContext(c)

			logger.AccessCtx(ctx,
				"[HTTP] Body Dump",
				"url_path", c.Request().URL.String(),
				"method", c.Request().Method,
				"request_body", string(reqBody),
				"response_code", c.Response().Status,
				"response_body", resBody.String(),
				"response_time", d,
				"response_time_millis", d.Milliseconds(),
			)

			return
		}
	}
}

func BodyDumpSkipper(c echo.Context) bool {
	switch {
	// Skip healthcheck
	case c.Request().URL.String() == "/api/health_check":
		fallthrough
	case strings.HasPrefix(c.Request().URL.String(), "/?"):
		fallthrough
	// Skip prometheus metrics
	case strings.HasPrefix(c.Request().URL.String(), "/metrics"):
		return true
	}

	return false
}

func BodyDumpHandler(c echo.Context, req []byte, resp []byte) {
	ctx := GetContext(c)

	logger.InfoCtx(ctx, "[HTTP] Body Dump",
		"url_path", c.Request().URL.String(),
		"method", c.Request().Method,
		"request_body", string(req),
		"response_code", c.Response().Status,
		"response_body", string(resp),
	)
}
