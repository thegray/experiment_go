package api

import (
	"experiment_go/kafka/sarama1/internal/data/dto"
	"experiment_go/kafka/sarama1/internal/pkg/contextid"
	"experiment_go/kafka/sarama1/internal/pkg/middleware"
	"experiment_go/kafka/sarama1/internal/service"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Integrator struct {
	svc service.Service
}

func Demo(e *echo.Echo) {
	ig := Integrator{
		svc: service.New(),
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

	res, _ := ig.svc.ProcessDemoMessage(ctx, req.Msg)
	return c.JSON(res.Code, res)
}
