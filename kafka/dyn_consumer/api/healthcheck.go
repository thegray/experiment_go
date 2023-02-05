package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthCheck(e *echo.Echo) {
	e.GET("/api/health_check", func(c echo.Context) error {
		res := `"code":"success","msg":"Healthy"`
		return c.JSON(http.StatusOK, res)
	})
}
