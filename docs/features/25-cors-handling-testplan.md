# 🧪 Test Plan: CORS Handling

> **Feature**: `25` — CORS Handling
> **Total New Test Cases**: 6

---

## New Headers

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-01 | Default AllowHeaders includes X-CSRF-Token | Default CORS, GET request | `Access-Control-Allow-Headers` contains "X-CSRF-Token" |
| TC-02 | AllowCredentials header set | Default CORS, GET request | `Access-Control-Allow-Credentials: true` |
| TC-03 | ExposeHeaders set | Default CORS, GET request | `Access-Control-Expose-Headers` contains "Content-Length" and "X-Request-ID" |

## Environment Configuration

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-04 | CORS_ALLOWED_ORIGINS env overrides default | Set env to `https://example.com` | `Access-Control-Allow-Origin: https://example.com` |

## Custom Config

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-05 | Custom config disables credentials | `AllowCredentials: false` | No `Access-Control-Allow-Credentials` header |
| TC-06 | Preflight with new headers | OPTIONS request with defaults | 204 with credentials + expose headers present |

---

## Existing Tests (Must Still Pass)

- TestCORS_DefaultHeaders
- TestCORS_PreflightOptions
- TestCORS_CustomConfig

## Acceptance Criteria

1. All 6 new tests pass
2. All 3 existing CORS tests still pass
3. Full regression (`go test ./... -count=1`) — 0 failures
4. `go vet ./...` — clean
5. No new dependencies in `go.mod`
