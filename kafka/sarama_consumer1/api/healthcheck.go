package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthCheck(e *echo.Echo) {
	e.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "healthcheck OK")
	})
}
