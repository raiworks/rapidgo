# Feature #54 — Prometheus Metrics: Test Plan

> **Status**: 🟡 IN PROGRESS
> **Depends On**: #07 (Router)
> **Branch**: `docs/54-prometheus-metrics`

---

## Test File: `core/metrics/metrics_test.go`

### Test Cases (12 total)

| ID | Test Name | What It Verifies |
|---|---|---|
| T01 | `TestNew` | `New()` returns non-nil `*Metrics` with all collectors registered |
| T02 | `TestMiddleware_IncrementsCounter` | After a request, `http_requests_total` counter is incremented |
| T03 | `TestMiddleware_RecordsDuration` | After a request, `http_request_duration_seconds` histogram has an observation |
| T04 | `TestMiddleware_RecordsSize` | After a request, `http_response_size_bytes` histogram has an observation |
| T05 | `TestMiddleware_LabelsMethod` | Counter has correct `method` label (GET, POST) |
| T06 | `TestMiddleware_LabelsPath` | Counter has correct `path` label using route template |
| T07 | `TestMiddleware_LabelsStatus` | Counter has correct `status` label (200, 404) |
| T08 | `TestMiddleware_UnmatchedPath` | Unmatched routes use `"unmatched"` as path label |
| T09 | `TestHandler` | `Handler()` returns Prometheus text format with registered metrics |
| T10 | `TestMiddleware_MultipleRequests` | Multiple requests correctly accumulate counters |
| T11 | `TestMiddleware_DifferentStatusCodes` | Different status codes produce distinct label sets |
| T12 | `TestHandler_IncludesGoMetrics` | Handler output includes default Go runtime metrics (`go_goroutines`, etc.) |

### Test Approach

- **Registry isolation**: Each test creates an isolated metrics instance via `newMetrics(prometheus.NewRegistry())`. This unexported constructor accepts a custom `prometheus.Registerer`, preventing cross-test pollution. For handler tests, `promhttp.HandlerFor(registry, ...)` serves metric output from the isolated registry only.
- **T01**: Call `New()`, assert non-nil. Verify metric families exist via a fresh registry.
- **T02–T08, T10–T11**: Create a Gin engine with `m.Middleware()` + test routes, send requests via `httptest`, then scrape handler output to assert label values and counts.
- **T09, T12**: Create a Gin engine with handler route, GET endpoint, assert response contains expected metric names.

### Acceptance Criteria

1. All 12 tests pass with `go test ./core/metrics/...`
2. Full regression passes: `go test ./...`
3. Binary builds clean: `go build -o bin/rapidgo.exe ./cmd`
