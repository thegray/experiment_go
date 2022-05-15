package api

import (
	"experiment_go/kafka/sarama1/internal/data/dto"
	"experiment_go/kafka/sarama1/internal/pkg/contextid"
	"experiment_go/kafka/sarama1/internal/pkg/kafka"
	"experiment_go/kafka/sarama1/internal/pkg/middleware"
	"experiment_go/kafka/sarama1/internal/service"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Integrator struct {
	svc service.Service
}

func Demo(e *echo.Echo, sp *kafka.SaramaProducer) {
	ig := Integrator{
		svc: service.New(sp),
	}

	route := e.Group("/v1")
	route.POST("/demo", ig.receiveDemoMessage)
}

func (ig Integrator) receiveDemoMessage(c echo.Context) error {
	ctx := middleware.GetContext(c)
	var req dto.MessageDemo
	if err := c.Bind(&req); err != nil {
		log.Println("request data invalid",
			"context_id", contextid.Value(ctx),
			"error", err.Error(),
		)

		return c.JSON(http.StatusBadRequest, nil)
	}

	res, err := ig.svc.ProcessDemoMessage(ctx, req.Msg)
	if err != nil {

		return c.JSON(http.StatusInternalServerError,
			err, // just for demo, can present this error more proper and good
		)
	}
	return c.JSON(res.Code, res)
}
