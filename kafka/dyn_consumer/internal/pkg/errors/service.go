package errors

import "fmt"

type ServiceError interface {
	error
	Unwrap() error
	Code() string
	Message() string
}

type serviceError struct {
	code    string
	message string
	inner   error
}

func (err serviceError) Error() string {
	if err.inner != nil {
		return fmt.Sprintf("[%s] %s ; %s", err.code, err.message, err.inner)
	}

	return fmt.Sprintf("[%s] %s", err.code, err.message)
}

func (err serviceError) Unwrap() error {
	return err.inner
}

func (err serviceError) Code() string {
	return err.code
}

func (err serviceError) Message() string {
	return err.message
}
