package api

import (
	"experiment_go/kafka/dyn_consumer/api/dto"
	"experiment_go/kafka/dyn_consumer/internal/data/repo"
	"experiment_go/kafka/dyn_consumer/internal/pkg/errors"
	"experiment_go/kafka/dyn_consumer/internal/pkg/kafka"
	"experiment_go/kafka/dyn_consumer/internal/pkg/middleware"
	"experiment_go/kafka/dyn_consumer/internal/pkg/time"
	"experiment_go/kafka/dyn_consumer/internal/service/kafkawatcher"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

type kwIntegrator struct {
	svc kafkawatcher.KWService
}

// create kafka config and start consuming from it
func KafkaWatcher(e *echo.Echo) {
	kwRepo := repo.KWRepo()
	kwService := kafkawatcher.NewKWService(kwRepo)

	integrator := kwIntegrator{
		svc: kwService,
	}

	e.POST("/api/kw/config", integrator.createWatcher()) // use kafka/config model

	// get messages from a config
	e.GET("/api/kw/config/:id", nil)

	// dummy generate watcher in local data
	e.POST("/api/kw/dummy/generate", integrator.dummyGenerate())

	// dummy print local data to log
	e.GET("/api/kw/dummy/print", integrator.dummyPrint())
}

// dummy generate watcher in local data
func (i *kwIntegrator) dummyGenerate() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := middleware.GetContext(c)

		req := struct {
			Key string `json:"key"`
		}{}
		if err := c.Bind(&req); err != nil {
			return errors.ErrInvalidRequest(err)
		}

		err := i.svc.DummyGenerate(ctx, req.Key)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, dto.ResponseEntity{
			Code:    "success",
			Message: "generate ok",
		})
	}
}

// dummy print local data to log
func (i *kwIntegrator) dummyPrint() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := middleware.GetContext(c)

		err := i.svc.DummyPrintWatcher(ctx)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, dto.ResponseEntity{
			Code:    "success",
			Message: "print data to log",
		})
	}
}

func (i *kwIntegrator) createWatcher() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := middleware.GetContext(c)

		var req kafka.ConsumerConfig
		if err := c.Bind(&req); err != nil {
			return errors.ErrInvalidRequest(err)
		}

		if err := req.Validate(); err != nil {
			return err
		}

		result, err := i.svc.CreateNewWatcher(ctx, req)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, dto.ResponseEntity{
			Code:    "success",
			Message: "create watcher success",
			Data: dto.CreateWatcherResponse{
				ID:          result.ID,
				IsRunning:   result.IsRunning,
				CreatedAt:   time.GetTimeForDisplay(result.CreatedAt),
				KafkaConfig: req,
			},
		})
	}
}
