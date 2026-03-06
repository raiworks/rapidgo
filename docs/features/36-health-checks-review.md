# Feature #36 — Health Checks: Review

## Summary

Two health-check endpoints for Docker/Kubernetes liveness and readiness probes.

## Delivered

| Item | Detail |
|------|--------|
| Package | `core/health/health.go` |
| Liveness | `GET /health` → 200 `{"status":"ok"}` |
| Readiness | `GET /health/ready` → pings DB, 200 or 503 |
| Integration | `RouterProvider.Boot()` auto-registers when `"db"` is bound |
| Tests | 3 (liveness OK, readiness OK, readiness with closed DB) |

## Design Notes

- `Routes()` accepts `func() *gorm.DB` (lazy resolver) instead of `*gorm.DB` directly. This prevents eager DB resolution during boot — the DB is only resolved on the first `/health/ready` request.
- Health routes are guarded by `c.Has("db")` — they are not registered if no database provider is configured.

## Blueprint Compliance

Matches blueprint exactly: same paths, same JSON responses, same HTTP status codes.

## Test Results

All 30 packages pass. `go vet` clean.
