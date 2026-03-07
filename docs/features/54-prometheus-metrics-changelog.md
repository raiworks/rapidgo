# Feature #54 — Prometheus Metrics: Changelog

> **Status**: 🟡 IN PROGRESS
> **Depends On**: #07 (Router)
> **Branch**: `docs/54-prometheus-metrics`

---

## Files Changed

| File | Action | Description |
|---|---|---|
| `core/metrics/metrics.go` | **NEW** | `Metrics` struct, `New`, `newMetrics`, `Middleware`, `Handler` |
| `core/metrics/metrics_test.go` | **NEW** | 12 test cases (T01–T12) |
| `app/providers/router_provider.go` | MODIFIED | Create metrics, apply globally via `r.Use()`, register `/metrics` endpoint |

## New Environment Variables

| Variable | Default | Description |
|---|---|---|
| `METRICS_ENABLED` | `false` | Enable Prometheus metrics |
| `METRICS_PATH` | `/metrics` | Metrics endpoint path |

## New Dependency

| Package | License | Purpose |
|---|---|---|
| `github.com/prometheus/client_golang` | Apache-2.0 | Prometheus Go client (counters, histograms, exposition handler) |

## Metrics Exposed

| Name | Type | Labels | Description |
|---|---|---|---|
| `http_requests_total` | Counter | method, path, status | Total HTTP request count |
| `http_request_duration_seconds` | Histogram | method, path, status | Request duration distribution |
| `http_response_size_bytes` | Histogram | method, path, status | Response size distribution |

## Migrations

None.

## Breaking Changes

None. Feature is opt-in via `METRICS_ENABLED=true`.
