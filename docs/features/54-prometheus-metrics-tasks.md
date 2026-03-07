# Feature #54 — Prometheus Metrics: Tasks

> **Status**: 🟡 IN PROGRESS
> **Depends On**: #07 (Router)
> **Branch**: `docs/54-prometheus-metrics`

---

## Task List

### T1 — Add `prometheus/client_golang` dependency

- [ ] Run `go get github.com/prometheus/client_golang`
- [ ] Verify Apache-2.0 license

### T2 — Create `core/metrics/metrics.go`

- [ ] Create new file `core/metrics/metrics.go`
- [ ] Define `Metrics` struct with three Prometheus collector fields
- [ ] Implement `New()` — registers `http_requests_total` (CounterVec), `http_request_duration_seconds` (HistogramVec), `http_response_size_bytes` (HistogramVec)
- [ ] Implement `Middleware() gin.HandlerFunc` — records start time, calls `c.Next()`, then observes duration/count/size
- [ ] Implement `Handler() gin.HandlerFunc` — wraps `promhttp.Handler()` for Gin
- [ ] Path normalization: use `c.FullPath()`, fall back to `"unmatched"`
- [ ] Status label: `strconv.Itoa(c.Writer.Status())`

### T3 — Register metrics in `RouterProvider`

- [ ] Modify `app/providers/router_provider.go`
- [ ] Import `core/metrics` and `core/config`
- [ ] When `METRICS_ENABLED=true`: create `metrics.New()`, apply globally via `r.Use(m.Middleware())`, register `GET` route at `METRICS_PATH` (default `/metrics`) using `metrics.Handler()`

### T4 — Write tests in `core/metrics/metrics_test.go`

- [ ] Create `core/metrics/metrics_test.go`
- [ ] 12 test cases (T01–T12) per test plan
- [ ] Use `newMetrics(prometheus.NewRegistry())` for isolated test instances

### T5 — Verify full regression

- [ ] Run `go test ./...` — all packages pass
- [ ] Run `go build -o bin/rapidgo.exe ./cmd` — binary builds clean
