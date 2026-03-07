# Feature #54 — Prometheus Metrics: Discussion

> **Status**: 🟡 IN PROGRESS
> **Depends On**: #07 (Router)
> **Branch**: `docs/54-prometheus-metrics`

---

## What Problem Does This Solve?

Production applications need observability. Operators must know how many requests are being served, how fast responses are, which endpoints are slow, and how many errors occur. Without built-in metrics, every team re-invents instrumentation from scratch.

Currently, RapidGo has no metrics export. There is no way to monitor request throughput, latency distributions, or error rates without adding custom code.

## What Does This Feature Add?

A `core/metrics` package that exposes Prometheus-compatible metrics via a Gin middleware and a `/metrics` HTTP endpoint. The middleware automatically records request count, duration, and response size for every HTTP request. The endpoint serves metrics in the Prometheus text exposition format.

### Key Design Decisions

1. **Gin middleware for collection** — A single middleware function intercepts every request, records the method, path, and status code, and measures duration. This follows the established middleware pattern used by all other framework middleware.

2. **Standard Prometheus client** — Uses `prometheus/client_golang`, the official Go Prometheus library. This is the de facto standard and provides histogram, counter, and gauge types, plus a built-in HTTP handler for the `/metrics` endpoint.

3. **Three core metrics** — Following the RED method (Rate, Errors, Duration):
   - `http_requests_total` (Counter) — total requests by method, path, status
   - `http_request_duration_seconds` (Histogram) — latency distribution by method, path, status
   - `http_response_size_bytes` (Histogram) — response body size distribution by method, path, status

4. **Opt-in via env var** — Metrics are enabled when `METRICS_ENABLED=true`. When disabled, no middleware is registered, no `/metrics` endpoint is exposed, and there is zero overhead.

5. **Path normalization** — Uses the Gin route template (e.g., `/users/:id`) rather than the raw path (`/users/42`) to avoid high-cardinality label explosion.

## What's Out of Scope?

- **Custom application metrics** — Developers can register their own metrics using `prometheus/client_golang` directly. This feature provides only the HTTP instrumentation layer.
- **Push gateway support** — Prometheus pulls metrics; no push gateway integration.
- **Database or cache metrics** — Only HTTP request metrics. GORM/Redis metrics can be added later.
- **Grafana dashboards** — Infrastructure concern, not a framework feature.
- **Authentication on `/metrics`** — The endpoint is unprotected. Operators should restrict access at the reverse proxy or network level.

## How Will Developers Use It?

```env
# .env
METRICS_ENABLED=true
METRICS_PATH=/metrics    # optional, defaults to /metrics
```

When enabled, the framework automatically:
1. Registers the metrics middleware globally
2. Exposes a `/metrics` endpoint serving Prometheus text format

No code changes required. Just set the env var and point Prometheus at the `/metrics` endpoint.

For custom metrics, developers import `prometheus/client_golang` directly and register with the default registry — the `/metrics` endpoint will include them automatically.
