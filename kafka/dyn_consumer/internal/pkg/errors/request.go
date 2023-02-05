package errors

import (
	"fmt"
)

const (
	requestInvalid               = "request:invalid"
	requestMissingMandatoryField = "request:missing-mandatory-field"
	requestInvalidFieldFormat    = "request:invalid-field-request"
	requestDuplicate             = "request:duplicate"
	requestNotFound              = "request:not-found"
)

var (
	ErrInvalidRequest = func(err error) ServiceError {
		return serviceError{
			code:    requestInvalid,
			message: "invalid request",
			inner:   err,
		}
	}

	ErrMissingMandatoryField = func(msg string, err error) ServiceError {
		return serviceError{
			code:    requestMissingMandatoryField,
			message: msg,
			inner:   err,
		}
	}

	ErrInvalidFieldFormat = func(msg string, err error) ServiceError {
		return serviceError{
			code:    requestInvalidFieldFormat,
			message: msg,
			inner:   err,
		}
	}

	ErrDuplicateRequest = func(key string) ServiceError {
		return serviceError{
			code:    requestDuplicate,
			message: fmt.Sprintf("duplicate request: %s", key),
		}
	}

	ErrNotFound = func(msg string, err error) ServiceError {
		return serviceError{
			code:    requestNotFound,
			message: msg,
			inner:   err,
		}
	}
)
