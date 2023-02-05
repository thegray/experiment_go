package errors

const (
	systemCommonError       = "system:common:error"
	systemNotImplementedYet = "system:common:not-implemented-yet"
)

var (
	ErrCommon = func(err error) ServiceError {
		return serviceError{
			code:    systemCommonError,
			message: "general server error",
			inner:   err,
		}
	}

	ErrNotImplementedYet = serviceError{
		code:    systemNotImplementedYet,
		message: "not implemented yet",
	}
)
