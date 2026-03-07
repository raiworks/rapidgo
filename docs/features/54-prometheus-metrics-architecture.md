# Feature #54 — Prometheus Metrics: Architecture

> **Status**: 🟡 IN PROGRESS
> **Depends On**: #07 (Router)
> **Branch**: `docs/54-prometheus-metrics`

---

## New Package: `core/metrics`

### File: `core/metrics/metrics.go`

#### Structs

```go
// Metrics holds Prometheus collectors for HTTP instrumentation.
type Metrics struct {
    requests *prometheus.CounterVec
    duration *prometheus.HistogramVec
    size     *prometheus.HistogramVec
}
```

#### Functions

```go
// New creates a Metrics instance and registers collectors with the default
// Prometheus registry.
func New() *Metrics

// Middleware returns a Gin middleware that records request count, duration,
// and response size for every HTTP request.
func (m *Metrics) Middleware() gin.HandlerFunc

// Handler returns an http.Handler that serves Prometheus metrics in the
// text exposition format. Wraps promhttp.Handler().
func Handler() gin.HandlerFunc
```

#### Unexported (for test isolation)

```go
// newMetrics creates a Metrics instance with the given registerer.
// Used by tests to provide an isolated prometheus.NewRegistry().
func newMetrics(reg prometheus.Registerer) *Metrics
```

`New()` calls `newMetrics(prometheus.DefaultRegisterer)`.

### Metric Definitions

| Metric | Type | Labels | Description |
|---|---|---|---|
| `http_requests_total` | CounterVec | method, path, status | Total HTTP requests |
| `http_request_duration_seconds` | HistogramVec | method, path, status | Request duration distribution |
| `http_response_size_bytes` | HistogramVec | method, path, status | Response body size distribution |

### Labels

| Label | Source | Example |
|---|---|---|
| `method` | `c.Request.Method` | `GET`, `POST` |
| `path` | `c.FullPath()` | `/users/:id`, `/health` |
| `status` | `c.Writer.Status()` | `200`, `404`, `500` |

**Path normalization**: `c.FullPath()` returns the route template registered with Gin (e.g., `/users/:id`), not the actual request path (e.g., `/users/42`). This prevents high-cardinality label explosion. If `c.FullPath()` is empty (unmatched route), falls back to `"unmatched"`.

### Histogram Buckets

- **Duration**: Prometheus default buckets (0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10 seconds)
- **Size**: `prometheus.ExponentialBuckets(100, 10, 7)` → 100, 1K, 10K, 100K, 1M, 10M, 100M bytes

### Middleware Flow

```
Request → Middleware (record start time)
    → c.Next() (handler runs)
    → Middleware (post-handler: observe duration, increment counter, observe size)
```

The middleware:
1. Records `time.Now()` before `c.Next()`
2. After `c.Next()`, reads `c.Writer.Status()`, `c.Writer.Size()`, `c.FullPath()`
3. Increments `http_requests_total` counter
4. Observes `http_request_duration_seconds` histogram
5. Observes `http_response_size_bytes` histogram

### Handler

```go
func Handler() gin.HandlerFunc {
    h := promhttp.Handler()
    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}
```

Wraps `promhttp.Handler()` (from `prometheus/client_golang`) into a `gin.HandlerFunc`. Serves the default Prometheus registry in text exposition format.

## Modified File: `app/providers/router_provider.go`

### Changes in `Boot()`

```go
// After health check registration, add:
if config.EnvBool("METRICS_ENABLED", false) {
    m := metrics.New()
    r.Use(m.Middleware())  // applies globally to all routes
    r.Get(config.Env("METRICS_PATH", "/metrics"), metrics.Handler())
}
```

All metrics logic lives in `RouterProvider.Boot()` — no changes to `MiddlewareProvider`.
`r.Use()` applies the middleware globally to every request, matching the health check pattern of registering framework routes in the router provider.

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `METRICS_ENABLED` | `false` | Enable Prometheus metrics collection and endpoint |
| `METRICS_PATH` | `/metrics` | URL path for the metrics endpoint |

## New Dependency

| Package | Version | License |
|---|---|---|
| `github.com/prometheus/client_golang` | latest | Apache-2.0 |

## No Changes

- `app/providers/middleware_provider.go` — unchanged
- `core/middleware/registry.go` — unchanged
- `core/router/router.go` — unchanged (r.Get and r.Use called from provider)
- `core/server/server.go` — unchanged
- `core/health/health.go` — unchanged
- No migrations
- No models
