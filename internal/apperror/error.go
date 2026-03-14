package apperror

import "fmt"

const (
	CodeInvalidInput = "INVALID_INPUT"
	CodeNotFound     = "NOT_FOUND"
	CodeInternal     = "INTERNAL_ERROR"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeConflict     = "CONFLICT"
)

type AppError struct {
	Code    string
	Message string
	Status  int
	Err     error
}

func (e *AppError) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code, message string, status int, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

func InvalidInput(message string, err error) *AppError {
	return New(CodeInvalidInput, message, 400, err)
}

func NotFound(message string, err error) *AppError {
	return New(CodeNotFound, message, 404, err)
}

func Internal(message string, err error) *AppError {
	return New(CodeInternal, message, 500, err)
}

func Unauthorized(message string, err error) *AppError {
	return New(CodeUnauthorized, message, 401, err)
}

func Forbidden(message string, err error) *AppError {
	return New(CodeForbidden, message, 403, err)
}

func Conflict(message string, err error) *AppError {
	return New(CodeConflict, message, 409, err)
}
