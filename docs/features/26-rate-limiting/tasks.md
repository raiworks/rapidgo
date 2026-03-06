# Feature #26 — Rate Limiting: Tasks

## Prerequisites

- [x] Middleware infrastructure shipped (#08)
- [x] Alias registration pattern established (#08, #24, #25)

## Implementation tasks

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 1 | `go get github.com/ulule/limiter/v3` and sub-packages | `go.mod`, `go.sum` | ⬜ |
| 2 | Create `RateLimitMiddleware()` | `core/middleware/ratelimit.go` | ⬜ |
| 3 | Register `"ratelimit"` alias | `app/providers/middleware_provider.go` | ⬜ |
| 4 | Write tests | `core/middleware/middleware_test.go` | ⬜ |
| 5 | Full regression + `go vet` | — | ⬜ |
| 6 | Commit, merge, review doc, roadmap update | — | ⬜ |

## Acceptance criteria

- `RateLimitMiddleware()` returns a `gin.HandlerFunc`.
- Default rate is `"60-M"` when `RATE_LIMIT` env is unset.
- Custom rate is read from `RATE_LIMIT` env var.
- Requests within the limit receive `200 OK` with `X-RateLimit-*` headers.
- Requests exceeding the limit receive `429 Too Many Requests`.
- `"ratelimit"` alias resolves in middleware registry.
- All existing tests pass (regression).
