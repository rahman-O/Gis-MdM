package errors

import "fmt"

// Code identifies an application error for API clients.
type Code string

const (
	CodeNotFound       Code = "NOT_FOUND"
	CodeUnauthorized   Code = "UNAUTHORIZED"
	CodeForbidden      Code = "FORBIDDEN"
	CodeValidation     Code = "VALIDATION"
	CodeInternal       Code = "INTERNAL"
	CodeNotImplemented Code = "NOT_IMPLEMENTED"
)

// AppError is a domain or use-case level error mapped to HTTP responses.
type AppError struct {
	Code    Code
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Cause }

func NotFound(msg string) *AppError {
	return &AppError{Code: CodeNotFound, Message: msg}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Code: CodeUnauthorized, Message: msg}
}

func NotImplemented(msg string) *AppError {
	return &AppError{Code: CodeNotImplemented, Message: msg}
}

func Internal(msg string, cause error) *AppError {
	return &AppError{Code: CodeInternal, Message: msg, Cause: cause}
}
