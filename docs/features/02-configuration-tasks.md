# 📋 Tasks: Configuration System

> **Feature**: `02` — Configuration System
> **Architecture**: [`02-configuration-architecture.md`](02-configuration-architecture.md)
> **Status**: 🟢 Ready for Build
> **Tasks**: 18 tasks across 5 phases

---

## Phase A — Dependency Setup (2 tasks)

- [ ] **A-01**: Run `go get github.com/joho/godotenv` to add first third-party dependency
- [ ] **A-02**: Verify `go.mod` has `require` block with godotenv entry and `go.sum` is created

## Phase B — Core Config Package (5 tasks)

- [ ] **B-01**: Create `core/config/config.go` with `Load()` function
- [ ] **B-02**: Create `core/config/env.go` with `Env()`, `EnvInt()`, `EnvBool()` helpers
- [ ] **B-03**: Create `core/config/environment.go` with `AppEnv()`, `IsProduction()`, `IsDevelopment()`, `IsTesting()`, `IsDebug()`
- [ ] **B-04**: Verify all exported functions have correct signatures matching architecture doc
- [ ] **B-05**: Run `go vet ./core/config/...` — must pass with zero warnings

## Phase C — Main Entry Point Update (3 tasks)

- [ ] **C-01**: Update `cmd/main.go` to import and call `config.Load()` at start
- [ ] **C-02**: Update `cmd/main.go` banner to display `APP_NAME`, `APP_PORT`, `APP_ENV`, `IsDebug()`
- [ ] **C-03**: Run `go build -o bin/rgo.exe ./cmd/` — must compile successfully

## Phase D — Tests (5 tasks)

- [ ] **D-01**: Create `core/config/config_test.go` with tests for `Load()` (with and without `.env`)
- [ ] **D-02**: Add tests for `Env()` — key present, key absent (fallback), empty string
- [ ] **D-03**: Add tests for `EnvInt()` — valid int, invalid string, empty (fallback)
- [ ] **D-04**: Add tests for `EnvBool()` — "true", "1", "false", "0", empty (fallback)
- [ ] **D-05**: Add tests for `AppEnv()`, `IsProduction()`, `IsDevelopment()`, `IsTesting()`, `IsDebug()`

## Phase E — Integration Validation (3 tasks)

- [ ] **E-01**: Run `go test ./core/config/... -v` — all tests must pass
- [ ] **E-02**: Run `go run ./cmd/` — verify banner shows `.env` values (APP_NAME=RGo, APP_PORT=8080, etc.)
- [ ] **E-03**: Run `go vet ./...` — full project must pass with zero warnings
