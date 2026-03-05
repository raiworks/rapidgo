# 🏗️ Architecture: Error Handling

> **Feature**: `04` — Error Handling
> **Discussion**: [`04-error-handling-discussion.md`](04-error-handling-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-05

---

## Overview

The error handling package provides structured error types for the RGo framework. It defines an `AppError` struct that carries an HTTP status code, a user-safe message, and an optional wrapped internal error. Constructor helpers create errors for common HTTP scenarios (400, 401, 403, 404, 409, 422, 500). The package is pure Go with no HTTP framework dependency — middleware integration comes in a later feature.

## File Structure

```
core/errors/
├── errors.go           # AppError type, constructors, ErrorResponse()
└── errors_test.go      # Unit tests for all exported functions
```

No existing files are modified. This is a new, standalone package.

## Component Design

### `AppError` Struct

**Responsibility**: Structured error type that carries HTTP status code, user-facing message, and wrapped internal error
**Package**: `core/errors`
**File**: `core/errors/errors.go`

```go
package errors

import (
	"github.com/RAiWorks/RGo/core/config"
)

// AppError represents a structured application error with HTTP status context.
type AppError struct {
	Code    int    // HTTP status code (e.g., 404, 500)
	Message string // User-facing safe message
	Err     error  // Internal error (logged, never exposed in production)
}
```

### Interface Compliance

`AppError` implements:
- `error` interface via `Error() string` — returns `Message`
- `Unwrap() error` — returns `Err` for `errors.As()` / `errors.Is()` support

```go
// Error returns the user-facing message. Implements the error interface.
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error for errors.Is/As support.
func (e *AppError) Unwrap() error {
	return e.Err
}
```

### Constructor Helpers

Seven constructors for common HTTP error scenarios:

```go
// NotFound creates a 404 error.
func NotFound(message string) *AppError {
	return &AppError{Code: 404, Message: message}
}

// BadRequest creates a 400 error.
func BadRequest(message string) *AppError {
	return &AppError{Code: 400, Message: message}
}

// Internal creates a 500 error wrapping an internal error.
func Internal(err error) *AppError {
	return &AppError{Code: 500, Message: "internal server error", Err: err}
}

// Unauthorized creates a 401 error.
func Unauthorized(message string) *AppError {
	return &AppError{Code: 401, Message: message}
}

// Forbidden creates a 403 error.
func Forbidden(message string) *AppError {
	return &AppError{Code: 403, Message: message}
}

// Conflict creates a 409 error.
func Conflict(message string) *AppError {
	return &AppError{Code: 409, Message: message}
}

// Unprocessable creates a 422 error.
func Unprocessable(message string) *AppError {
	return &AppError{Code: 422, Message: message}
}
```

### ErrorResponse Helper

Returns a map suitable for JSON serialization. Debug-aware: includes internal error details only when `APP_DEBUG=true`.

```go
// ErrorResponse returns a map for JSON error responses. In debug mode,
// it includes internal error details. In production, only the safe message.
func (e *AppError) ErrorResponse() map[string]any {
	resp := map[string]any{
		"success": false,
		"error":   e.Message,
	}
	if config.IsDebug() && e.Err != nil {
		resp["internal"] = e.Err.Error()
	}
	return resp
}
```

## Data Flow

```
Application code
    → calls errors.NotFound("user not found") or errors.Internal(err)
    → returns *AppError with Code, Message, Err
    → caller can:
        1. Return it as error (Error() gives Message)
        2. Use errors.As() to extract *AppError and read Code
        3. Call ErrorResponse() to get JSON-ready map
        4. Log with slog using structured fields
```

## Configuration

No new environment variables. Uses existing:

| Key | Type | Default | Used For |
|---|---|---|---|
| `APP_DEBUG` | bool | `true` | `ErrorResponse()` — show/hide internal error details |

## Security Considerations

- **Production safety**: `ErrorResponse()` NEVER includes `internal` details when `APP_DEBUG=false`
- **Error() returns Message only**: The user-safe message, never the internal error string
- **Wrapped errors**: Internal errors are available via `Unwrap()` for code inspection but must never be sent to clients in production
- **No stack traces**: Stack trace capture is not in scope — framework uses structured logging context instead

## Trade-offs & Alternatives

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| `AppError` struct with constructors | Simple, idiomatic Go, works with `errors.As/Is`, clear API | Not extensible to rich validation errors (yet) | ✅ Selected |
| Sentinel error vars (`ErrNotFound = errors.New(...)`) | Simpler, standard Go pattern | No HTTP status code, no message customization | ❌ Insufficient for HTTP framework |
| Interface-based errors (`type HTTPError interface`) | Maximum flexibility, polymorphism | Over-engineered for current needs | ❌ Premature abstraction |
| Third-party error package (pkg/errors, etc.) | Stack traces, rich context | Extra dependency, most features unused | ❌ Unnecessary dependency |

## Next

Create tasks doc → `04-error-handling-tasks.md`
