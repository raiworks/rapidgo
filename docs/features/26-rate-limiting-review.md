# Feature #26 — Rate Limiting: Review

## Summary

Implemented `RateLimitMiddleware()` in `core/middleware/ratelimit.go` using `ulule/limiter/v3` with in-memory store. Registered `"ratelimit"` alias in `middleware_provider.go`.

## Files changed

| File | Change |
|------|--------|
| `core/middleware/ratelimit.go` | New — `RateLimitMiddleware()` with `RATE_LIMIT` env var (default `"60-M"`) |
| `app/providers/middleware_provider.go` | Added `"ratelimit"` alias registration |
| `core/middleware/middleware_test.go` | +6 tests (TC-32 to TC-37) |
| `go.mod` / `go.sum` | Added `github.com/ulule/limiter/v3` v3.11.2 |

## Test results

| TC | Description | Result |
|----|-------------|--------|
| TC-32 | Default rate allows requests within limit | ✅ PASS |
| TC-33 | X-RateLimit-Limit header present | ✅ PASS |
| TC-34 | X-RateLimit-Remaining decrements | ✅ PASS |
| TC-35 | Exceed limit returns 429 | ✅ PASS |
| TC-36 | Custom RATE_LIMIT env var respected | ✅ PASS |
| TC-37 | "ratelimit" alias resolves | ✅ PASS |

## Regression

- All 23 packages pass.
- `go vet` clean.

## Deviation log

_None._

## Commit

`50736c5` — merged to `main`.
