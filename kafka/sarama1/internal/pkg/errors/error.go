package errors

import "errors"

var (
	ErrInvalidRequest      = errors.New("invalid-request")
	ErrInternalServerError = errors.New("internal-server-error")
	ErrNotFound            = errors.New("not-found")
)

type ServiceError struct {
	InnerError error
	ErrCode    string
	ErrMsg     string
}

func (err ServiceError) Error() string {
	return err.ErrMsg
}

func (err ServiceError) Code() string {
	return err.ErrCode
}
