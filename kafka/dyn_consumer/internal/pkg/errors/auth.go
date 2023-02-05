package errors

const (
	systemJWTDecodeSecretError  = "system:jwt.go:decode-secret"
	systemAuthInvalidTokenError = "system:auth:invalid-token"
	systemAuthForbiddenError    = "system:auth:forbidden"
)

var (
	ErrJWTDecodeSecret = func(msg string, err error) ServiceError {
		return serviceError{
			code:    systemJWTDecodeSecretError,
			message: msg,
			inner:   err,
		}
	}
	ErrAuthInvalidToken = func(msg string, err error) ServiceError {
		return serviceError{
			code:    systemAuthInvalidTokenError,
			message: msg,
			inner:   err,
		}
	}
	ErrAuthForbidden = func(msg string, err error) ServiceError {
		return serviceError{
			code:    systemAuthForbiddenError,
			message: msg,
			inner:   err,
		}
	}
)
