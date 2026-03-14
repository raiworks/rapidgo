package errors

import (
	"errors"
	"io"
	"os"
	"testing"
)

// --- TC-01: AppError implements error interface ---
func TestAppError_ImplementsErrorInterface(t *testing.T) {
	appErr := &AppError{Status: 404, Code: "NOT_FOUND", Message: "not found"}
	var _ error = appErr // compile-time check

	if got := appErr.Error(); got != "not found" {
		t.Errorf("Error() = %q, want %q", got, "not found")
	}
}

// --- TC-02: Unwrap returns wrapped error, errors.Is works ---
func TestUnwrap_ReturnsWrappedError(t *testing.T) {
	appErr := Internal(io.ErrUnexpectedEOF)

	if appErr.Unwrap() != io.ErrUnexpectedEOF {
		t.Errorf("Unwrap() = %v, want %v", appErr.Unwrap(), io.ErrUnexpectedEOF)
	}
	if !errors.Is(appErr, io.ErrUnexpectedEOF) {
		t.Error("errors.Is(appErr, io.ErrUnexpectedEOF) = false, want true")
	}
}

// --- TC-03: Unwrap returns nil when no wrapped error ---
func TestUnwrap_NilErr(t *testing.T) {
	appErr := NotFound("user not found")

	if appErr.Unwrap() != nil {
		t.Errorf("Unwrap() = %v, want nil", appErr.Unwrap())
	}
}

// --- TC-04: errors.As extracts AppError ---
func TestErrorsAs_ExtractsAppError(t *testing.T) {
	var err error = NotFound("page not found")

	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatal("errors.As failed to extract *AppError")
	}
	if appErr.Status != 404 {
		t.Errorf("Status = %d, want 404", appErr.Status)
	}
}

// --- TC-05 through TC-11: Constructor tests ---
func TestConstructors(t *testing.T) {
	tests := []struct {
		name    string
		create  func() *AppError
		status  int
		code    string
		message string
		hasErr  bool
	}{
		// TC-05
		{"NotFound", func() *AppError { return NotFound("user not found") }, 404, "NOT_FOUND", "user not found", false},
		// TC-06
		{"BadRequest", func() *AppError { return BadRequest("invalid input") }, 400, "BAD_REQUEST", "invalid input", false},
		// TC-07
		{"Internal", func() *AppError { return Internal(io.EOF) }, 500, "INTERNAL_ERROR", "internal server error", true},
		// TC-08
		{"Unauthorized", func() *AppError { return Unauthorized("login required") }, 401, "UNAUTHORIZED", "login required", false},
		// TC-09
		{"Forbidden", func() *AppError { return Forbidden("access denied") }, 403, "FORBIDDEN", "access denied", false},
		// TC-10
		{"Conflict", func() *AppError { return Conflict("already exists") }, 409, "CONFLICT", "already exists", false},
		// TC-11
		{"Unprocessable", func() *AppError { return Unprocessable("validation failed") }, 422, "UNPROCESSABLE", "validation failed", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr := tt.create()
			if appErr.Status != tt.status {
				t.Errorf("Status = %d, want %d", appErr.Status, tt.status)
			}
			if appErr.Code != tt.code {
				t.Errorf("Code = %q, want %q", appErr.Code, tt.code)
			}
			if appErr.Message != tt.message {
				t.Errorf("Message = %q, want %q", appErr.Message, tt.message)
			}
			if tt.hasErr && appErr.Err == nil {
				t.Error("Err = nil, want non-nil")
			}
			if !tt.hasErr && appErr.Err != nil {
				t.Errorf("Err = %v, want nil", appErr.Err)
			}
		})
	}
}

// --- TC-12: ErrorResponse in debug mode includes internal error ---
func TestErrorResponse_DebugMode(t *testing.T) {
	os.Setenv("APP_DEBUG", "true")
	defer os.Unsetenv("APP_DEBUG")

	appErr := Internal(errors.New("db timeout"))
	resp := appErr.ErrorResponse()

	if resp["success"] != false {
		t.Errorf("success = %v, want false", resp["success"])
	}
	if resp["error"] != "internal server error" {
		t.Errorf("error = %v, want %q", resp["error"], "internal server error")
	}
	if resp["internal"] != "db timeout" {
		t.Errorf("internal = %v, want %q", resp["internal"], "db timeout")
	}
}

// --- TC-13: ErrorResponse in production — no internal details ---
func TestErrorResponse_ProductionMode(t *testing.T) {
	os.Setenv("APP_DEBUG", "false")
	defer os.Unsetenv("APP_DEBUG")

	appErr := Internal(errors.New("db timeout"))
	resp := appErr.ErrorResponse()

	if resp["success"] != false {
		t.Errorf("success = %v, want false", resp["success"])
	}
	if resp["error"] != "internal server error" {
		t.Errorf("error = %v, want %q", resp["error"], "internal server error")
	}
	if _, exists := resp["internal"]; exists {
		t.Errorf("internal key should not exist in production, got %v", resp["internal"])
	}
}

// --- TC-14: ErrorResponse with nil Err in debug mode — no internal key ---
func TestErrorResponse_DebugMode_NilErr(t *testing.T) {
	os.Setenv("APP_DEBUG", "true")
	defer os.Unsetenv("APP_DEBUG")

	appErr := NotFound("not found")
	resp := appErr.ErrorResponse()

	if resp["success"] != false {
		t.Errorf("success = %v, want false", resp["success"])
	}
	if resp["error"] != "not found" {
		t.Errorf("error = %v, want %q", resp["error"], "not found")
	}
	if _, exists := resp["internal"]; exists {
		t.Error("internal key should not exist when Err is nil")
	}
}

// --- TC-15: WithCode overrides default code ---
func TestWithCode_OverridesDefault(t *testing.T) {
	appErr := BadRequest("invalid email").WithCode("INVALID_EMAIL")
	if appErr.Code != "INVALID_EMAIL" {
		t.Errorf("Code = %q, want %q", appErr.Code, "INVALID_EMAIL")
	}
	if appErr.Status != 400 {
		t.Errorf("Status = %d, want 400", appErr.Status)
	}
}

// --- TC-16: ErrorResponse includes code field ---
func TestErrorResponse_IncludesCode(t *testing.T) {
	os.Setenv("APP_DEBUG", "false")
	defer os.Unsetenv("APP_DEBUG")

	appErr := NotFound("user not found")
	resp := appErr.ErrorResponse()

	if resp["code"] != "NOT_FOUND" {
		t.Errorf("code = %v, want %q", resp["code"], "NOT_FOUND")
	}
	if resp["success"] != false {
		t.Errorf("success = %v, want false", resp["success"])
	}
	if resp["error"] != "user not found" {
		t.Errorf("error = %v, want %q", resp["error"], "user not found")
	}
}

// --- TC-17: HTTPStatus deprecated helper ---
func TestHTTPStatus_ReturnsStatus(t *testing.T) {
	appErr := Forbidden("no access")
	if appErr.HTTPStatus() != 403 {
		t.Errorf("HTTPStatus() = %d, want 403", appErr.HTTPStatus())
	}
	if appErr.HTTPStatus() != appErr.Status {
		t.Error("HTTPStatus() should equal Status")
	}
}
