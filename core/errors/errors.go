package errors

import (
	"github.com/raiworks/rapidgo/v2/core/config"
)

// AppError represents a structured application error with HTTP status context.
type AppError struct {
	Status  int    // HTTP status code (e.g., 404, 500)
	Code    string // Machine-readable error code (e.g., "NOT_FOUND", "VALIDATION_FAILED")
	Message string // User-facing safe message
	Err     error  // Internal error (logged, never exposed in production)
}

// Error returns the user-facing message. Implements the error interface.
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error for errors.Is/As support.
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatus returns the HTTP status code.
// Deprecated: Access .Status directly instead.
func (e *AppError) HTTPStatus() int {
	return e.Status
}

// WithCode returns the error with a custom machine-readable code,
// overriding the default code set by the factory function.
func (e *AppError) WithCode(code string) *AppError {
	e.Code = code
	return e
}

// NotFound creates a 404 error.
func NotFound(message string) *AppError {
	return &AppError{Status: 404, Code: "NOT_FOUND", Message: message}
}

// BadRequest creates a 400 error.
func BadRequest(message string) *AppError {
	return &AppError{Status: 400, Code: "BAD_REQUEST", Message: message}
}

// Internal creates a 500 error wrapping an internal error.
func Internal(err error) *AppError {
	return &AppError{Status: 500, Code: "INTERNAL_ERROR", Message: "internal server error", Err: err}
}

// Unauthorized creates a 401 error.
func Unauthorized(message string) *AppError {
	return &AppError{Status: 401, Code: "UNAUTHORIZED", Message: message}
}

// Forbidden creates a 403 error.
func Forbidden(message string) *AppError {
	return &AppError{Status: 403, Code: "FORBIDDEN", Message: message}
}

// Conflict creates a 409 error.
func Conflict(message string) *AppError {
	return &AppError{Status: 409, Code: "CONFLICT", Message: message}
}

// Unprocessable creates a 422 error.
func Unprocessable(message string) *AppError {
	return &AppError{Status: 422, Code: "UNPROCESSABLE", Message: message}
}

// ErrorResponse returns a map for JSON error responses. In debug mode,
// it includes internal error details. In production, only the safe message.
func (e *AppError) ErrorResponse() map[string]any {
	resp := map[string]any{
		"success": false,
		"error":   e.Message,
		"code":    e.Code,
	}
	if config.IsDebug() && e.Err != nil {
		resp["internal"] = e.Err.Error()
	}
	return resp
}
