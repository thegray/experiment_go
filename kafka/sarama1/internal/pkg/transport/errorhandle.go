package transport

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"experiment_go/kafka/sarama1/internal/pkg/contextid"
	general_error "experiment_go/kafka/sarama1/internal/pkg/errors"
	"experiment_go/kafka/sarama1/internal/pkg/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ErrHandler struct {
	E *echo.Echo
}

func getVldErrorMsg(validationErr validator.FieldError) string {
	var reason string
	switch validationErr.Tag() {
	case "required":
		reason = "this field is required"
	case "numeric":
		reason = "this field should only contains numeric value"
	case "alpha":
		reason = "this field should only contains alphabet value"
	case "alphanum":
		reason = "this field should only contains alphanumeric value"
	case "email":
		reason = "this field should be a valid email address"
	case "url":
		reason = "this field should be a valid URL"
	case "max":
		reason = fmt.Sprintf("this field should not be longer than %s character(s)", validationErr.Param())
	case "min":
		reason = fmt.Sprintf("this field should not be shorter than %s character(s)", validationErr.Param())
	case "oneof":
		reason = fmt.Sprintf("this field should be one of %s", validationErr.Param())
	}
	return reason
}

type resp struct {
	ErrorObj errorCodeMessage `json:"error"`
}

type errorCodeMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (ce ErrHandler) Handler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
	)

	ctx := middleware.GetContext(c)
	errObj := errorCodeMessage{}
	switch e := err.(type) {
	case *echo.HTTPError:
		code = e.Code
		errObj.Code = general_error.ErrInternalServerError.Error()
		errObj.Message = fmt.Sprintf("%v", e.Message)

		log.Printf("[ErrorHandler] %v %v", e.Message, contextid.Value(ctx))
		// logger.Info(fmt.Sprintf("[ErrorHandler] %v", e.Message),
		// "context_id", contextid.Value(ctx),
		// )
	case validator.ValidationErrors:
		var errMsg []string
		for _, v := range e {
			errMsg = append(errMsg, fmt.Sprintf("invalid value on %s, %s", v.Field(), getVldErrorMsg(v)))
		}
		errObj.Code = general_error.ErrInvalidRequest.Error()
		errObj.Message = strings.Join(errMsg, ",")
		code = http.StatusBadRequest

		log.Printf("[ErrorHandler] %v %v", strings.Join(errMsg, ","), contextid.Value(ctx))
		// logger.Info(fmt.Sprintf("[ErrorHandler] %v", strings.Join(errMsg, ",")),
		// 	"context_id", contextid.Value(ctx),
		// )
	default:
		errObj.Code = general_error.ErrInternalServerError.Error()
		errObj.Message = e.Error()

		log.Printf("[ErrorHandler] %v %v", e.Error(), contextid.Value(ctx))
		// logger.Info(fmt.Sprintf("[ErrorHandler] %v", e.Error()),
		// 	"context_id", contextid.Value(ctx),
		// )
	}

	r := resp{ErrorObj: errObj}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == "HEAD" {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, r)
		}
		if err != nil {
			ce.E.Logger.Error(err)
		}
	}
}
