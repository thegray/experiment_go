package transport

import (
	"context"
	"net/http"
	"os"
	"time"

	"experiment_go/kafka/sarama1/internal/model"
	"experiment_go/kafka/sarama1/internal/pkg/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type server struct {
	e            *echo.Echo
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewServer(cfg model.ServerConfig) server {
	e := echo.New()

	// global middlewares
	e.Use(
		middleware.ContextMiddleware(),
		middleware.Logger(),
		echoMiddleware.BodyDumpWithConfig(echoMiddleware.BodyDumpConfig{
			Skipper: middleware.BodyDumpSkipper,
			Handler: middleware.BodyDumpHandler,
		}),
		middleware.Recover(),
		middleware.CORS(),
		middleware.Headers(),
	)

	e.HTTPErrorHandler = ErrHandler{E: e}.Handler

	e.Debug = cfg.Loglevel == "DEBUG"
	srv := server{
		e:            e,
		port:         cfg.Port,
		readTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		writeTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}

	return srv
}

func (srv server) Engine() *echo.Echo {
	return srv.e
}

func (srv server) Start() {
	s := &http.Server{
		Addr:         srv.port,
		ReadTimeout:  srv.readTimeout,
		WriteTimeout: srv.writeTimeout,
	}

	if err := srv.e.StartServer(s); err != nil {
		srv.e.Logger.Error(err)
		srv.e.Logger.Info("Shutting down the server")
		os.Exit(1)
	}
}

func (srv server) Stop() {
	ctx := context.Background()
	if err := srv.e.Shutdown(ctx); err != nil {
		srv.e.Logger.Fatal(err)
	}
}
