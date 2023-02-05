package transport

import (
	"fmt"
	"net/http"

	"experiment_go/kafka/dyn_consumer/internal/pkg/errors"
	"experiment_go/kafka/dyn_consumer/internal/pkg/middleware"

	"experiment_go/kafka/dyn_consumer/internal/pkg/logger"
	echo "github.com/labstack/echo/v4"
)

type ErrHandler struct {
	E *echo.Echo
}

type resp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (ce ErrHandler) Handler(err error, c echo.Context) {
	var (
		statusCode = http.StatusInternalServerError
		ctx        = middleware.GetContext(c)
		r          = resp{}
	)

	switch e := err.(type) {
	case errors.ServiceError:
		statusCode = statusCodeByErrorCode(e.Code())

		r.Code = e.Code()
		r.Message = e.Message()
	case *echo.HTTPError:
		statusCode = e.Code

		r.Code = fmt.Sprint(e.Code)
		r.Message = fmt.Sprint(e.Message)

		logger.ErrorCtx(ctx, "[ErrHandler] HTTPError",
			"error", e.Error(),
		)
	default:
		r.Code = fmt.Sprint(statusCode)
		r.Message = http.StatusText(statusCode)

		logger.ErrorCtx(ctx, "[ErrHandler] ",
			"error", e.Error(),
		)
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == "HEAD" {
			err = c.NoContent(statusCode)
		} else {
			err = c.JSON(statusCode, r)
		}
		if err != nil {
			ce.E.Logger.Error(err)
		}
	}
}

func statusCodeByErrorCode(code string) int {
	switch code {
	case errors.ErrInvalidRequest(nil).Code():
		fallthrough
	case errors.ErrMissingMandatoryField("", nil).Code():
		fallthrough
	case errors.ErrDuplicateRequest("").Code():
		fallthrough
	case errors.ErrDatabaseDuplicateEntry(nil).Code():
		fallthrough
	case errors.ErrInvalidFieldFormat("", nil).Code():
		return http.StatusBadRequest
	case errors.ErrDatabaseNoRows(nil).Code():
		fallthrough
	case errors.ErrNotFound("", nil).Code():
		return http.StatusNotFound
	case errors.ErrAuthForbidden("", nil).Code():
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
