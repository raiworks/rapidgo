# 🧪 Test Plan: Error Handling

> **Feature**: `04` — Error Handling
> **Tasks**: [`04-error-handling-tasks.md`](04-error-handling-tasks.md)
> **Date**: 2026-03-05

---

## Acceptance Criteria

- [ ] `AppError` implements `error` interface — `Error()` returns the user-facing message
- [ ] `AppError` implements `Unwrap()` — works with `errors.As()` and `errors.Is()`
- [ ] All 7 constructors return `*AppError` with correct HTTP status codes
- [ ] `ErrorResponse()` returns safe map in production (no internal details)
- [ ] `ErrorResponse()` includes internal error details only when `APP_DEBUG=true`
- [ ] All tests pass with `go test ./core/errors/...`
- [ ] `go vet ./...` reports no issues

---

## Test Cases

### TC-01: AppError implements error interface

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Create `AppError{Code: 404, Message: "not found"}` → 2. Call `.Error()` → 3. Assign to `var _ error` |
| **Expected Result** | `.Error()` returns `"not found"`, compiles as `error` type |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-02: Unwrap returns wrapped error

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Create `Internal(io.ErrUnexpectedEOF)` → 2. Call `.Unwrap()` → 3. Use `errors.Is(appErr, io.ErrUnexpectedEOF)` |
| **Expected Result** | `Unwrap()` returns the original error, `errors.Is` returns `true` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-03: Unwrap returns nil when no wrapped error

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | None |
| **Steps** | 1. Create `NotFound("user not found")` → 2. Call `.Unwrap()` |
| **Expected Result** | Returns `nil` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-04: errors.As extracts AppError from error chain

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Create `NotFound("x")` → 2. Assign to `var err error` → 3. Use `errors.As(err, &appErr)` → 4. Check `appErr.Code` |
| **Expected Result** | `errors.As` returns `true`, `appErr.Code == 404` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-05: NotFound constructor — code 404

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Call `NotFound("user not found")` → 2. Check Code, Message, Err |
| **Expected Result** | `Code=404`, `Message="user not found"`, `Err=nil` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-06: BadRequest constructor — code 400

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Call `BadRequest("invalid input")` → 2. Check Code, Message |
| **Expected Result** | `Code=400`, `Message="invalid input"` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-07: Internal constructor — code 500, wraps error

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Create `sentinel := errors.New("db down")` → 2. Call `Internal(sentinel)` → 3. Check Code, Message, Err |
| **Expected Result** | `Code=500`, `Message="internal server error"`, `Err==sentinel` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-08: Unauthorized constructor — code 401

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Call `Unauthorized("login required")` → 2. Check Code |
| **Expected Result** | `Code=401`, `Message="login required"` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-09: Forbidden constructor — code 403

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Call `Forbidden("access denied")` → 2. Check Code |
| **Expected Result** | `Code=403`, `Message="access denied"` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-10: Conflict constructor — code 409

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Call `Conflict("already exists")` → 2. Check Code |
| **Expected Result** | `Code=409`, `Message="already exists"` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-11: Unprocessable constructor — code 422

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Call `Unprocessable("validation failed")` → 2. Check Code |
| **Expected Result** | `Code=422`, `Message="validation failed"` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-12: ErrorResponse in debug mode — includes internal error

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `APP_DEBUG=true` set in environment |
| **Steps** | 1. Set `APP_DEBUG=true` → 2. Create `Internal(errors.New("db timeout"))` → 3. Call `.ErrorResponse()` |
| **Expected Result** | Map contains `"success": false`, `"error": "internal server error"`, `"internal": "db timeout"` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-13: ErrorResponse in production — no internal details

| Property | Value |
|---|---|
| **Category** | Security |
| **Precondition** | `APP_DEBUG=false` set in environment |
| **Steps** | 1. Set `APP_DEBUG=false` → 2. Create `Internal(errors.New("db timeout"))` → 3. Call `.ErrorResponse()` |
| **Expected Result** | Map contains `"success": false`, `"error": "internal server error"`, NO `"internal"` key |
| **Status** | ⬜ Not Run |
| **Notes** | Security-critical: internal details must never leak in production |

### TC-14: ErrorResponse with nil Err in debug mode — no internal key

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | `APP_DEBUG=true` set in environment |
| **Steps** | 1. Set `APP_DEBUG=true` → 2. Create `NotFound("x")` (Err is nil) → 3. Call `.ErrorResponse()` |
| **Expected Result** | Map contains `"success": false`, `"error": "x"`, NO `"internal"` key (because Err is nil) |
| **Status** | ⬜ Not Run |
| **Notes** | Even in debug mode, nil Err means no internal key |

---

## Edge Cases

| # | Scenario | Expected Behavior |
|---|---|---|
| 1 | `AppError` with nil `Err` — `Unwrap()` | Returns `nil` gracefully |
| 2 | `ErrorResponse()` when `Err` is nil + debug mode | No `"internal"` key in map |
| 3 | `Internal(nil)` — creating 500 with nil error | `Code=500`, `Message="internal server error"`, `Err=nil`, `Unwrap()` returns nil |

## Security Tests

| # | Test | Expected |
|---|---|---|
| 1 | `ErrorResponse()` with `APP_DEBUG=false` and non-nil Err | No `"internal"` key — internal error details never exposed |
| 2 | `Error()` method on `Internal(err)` | Returns `"internal server error"` — never the internal error string |

---

## Test Summary

| Category | Total | Pass | Fail | Skip |
|---|---|---|---|---|
| Happy Path | 9 | — | — | — |
| Edge Cases | 3 | — | — | — |
| Security | 2 | — | — | — |
| **Total** | **14** | — | — | — |

**Result**: ⬜ NOT RUN
