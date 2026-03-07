# ✅ Tasks: Service Mode

> **Feature**: `56` — Service Mode Architecture
> **Architecture**: [`56-service-mode-architecture.md`](56-service-mode-architecture.md)
> **Branch**: `feature/56-service-mode`
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/14 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [ ] Feature branch created from latest `main`
- [x] Dependent features are merged to `main` (#55 complete)
- [ ] Test plan doc created
- [ ] Changelog doc created (empty)

---

## Phase A — Service Mode Infrastructure

> New `core/service/` package with Mode type, parsing, and validation.

- [ ] **A.1** — Create `core/service/mode.go`
  - [ ] `Mode` type as `uint8` bitmask
  - [ ] Constants: `ModeWeb`, `ModeAPI`, `ModeWS`, `ModeAll`
  - [ ] `modeNames` map for string-to-Mode lookup
  - [ ] `ParseMode(s string) (Mode, error)` — comma-separated parsing with validation
  - [ ] `Mode.Has(flag Mode) bool` — bitmask check
  - [ ] `Mode.String() string` — human-readable output
  - [ ] `Mode.Services() []Mode` — list of individual active modes
  - [ ] `Mode.PortEnvKey() string` — env var name for mode-specific port
- [ ] **A.2** — Create `core/service/mode_test.go`
  - [ ] Test `ParseMode` with valid inputs: `"all"`, `"web"`, `"api"`, `"ws"`, `"api,ws"`, `"web,api,ws"`
  - [ ] Test `ParseMode` with invalid inputs: `""`, `"invalid"`, `"api,invalid"`, `"worker"`
  - [ ] Test `Mode.Has()` combinations
  - [ ] Test `Mode.String()` output
  - [ ] Test `Mode.Services()` returns correct list
  - [ ] Test `Mode.PortEnvKey()` returns correct env var names
- [ ] **A.3** — Create `routes/ws.go`
  - [ ] `RegisterWS(r *router.Router)` — empty placeholder with comment
- [ ] 📍 **Checkpoint A** — `go test ./core/service/... -count=1` passes, `go build ./...` compiles

---

## Phase B — Mode-Aware Bootstrap

> Modify provider registration and route loading based on active mode.

- [ ] **B.1** — Update `core/cli/root.go`
  - [ ] Change `NewApp()` signature to `NewApp(mode service.Mode) *app.App`
  - [ ] Conditional DB provider: only if `mode.Has(ModeWeb|ModeAPI|ModeWS)`
  - [ ] Conditional Session provider: only if `mode.Has(ModeWeb)`
  - [ ] Pass `mode` to `MiddlewareProvider{Mode: mode}` and `RouterProvider{Mode: mode}`
  - [ ] Add import for `core/service`
- [ ] **B.2** — Update `app/providers/router_provider.go`
  - [ ] Add `Mode service.Mode` field to `RouterProvider` struct
  - [ ] Template/static setup conditional on `Mode.Has(ModeWeb)`
  - [ ] `routes.RegisterWeb(r)` conditional on `Mode.Has(ModeWeb)`
  - [ ] `routes.RegisterAPI(r)` conditional on `Mode.Has(ModeAPI)`
  - [ ] `routes.RegisterWS(r)` conditional on `Mode.Has(ModeWS)`
  - [ ] Health check remains: if `c.Has("db")` (unchanged)
- [ ] **B.3** — Update `app/providers/middleware_provider.go`
  - [ ] Add `Mode service.Mode` field to `MiddlewareProvider` struct
  - [ ] CSRF registration conditional on `Mode.Has(ModeWeb)`
  - [ ] All other middleware registered unconditionally (recovery, requestid, cors, error_handler, auth, ratelimit)
- [ ] **B.4** — Update `core/cli/serve.go`
  - [ ] Add `--mode` / `-m` flag (`serveMode` variable)
  - [ ] Mode resolution: `--mode` flag > `RAPIDGO_MODE` env > `"all"` default
  - [ ] Call `NewApp(mode)` instead of `NewApp()`
  - [ ] Add mode to startup banner output
  - [ ] Change `Run` to `RunE` for error return from `ParseMode`
- [ ] **B.5** — Update `.env`
  - [ ] Add `RAPIDGO_MODE=all` with comment
  - [ ] Add `WEB_PORT=`, `API_PORT=`, `WS_PORT=` (commented out, with defaults explanation)
- [ ] **B.6** — Update all other `NewApp()` callers
  - [ ] `core/cli/migrate.go` — pass `service.ModeAll`
  - [ ] `core/cli/migrate_rollback.go` — pass `service.ModeAll`
  - [ ] `core/cli/migrate_status.go` — pass `service.ModeAll`
  - [ ] `core/cli/seed.go` — pass `service.ModeAll`
  - [ ] Any other commands calling `NewApp()`
- [ ] **B.7** — Update tests for modified files
  - [ ] `app/providers/providers_test.go` — pass `Mode` field in provider constructors
  - [ ] `core/cli/cli_test.go` — verify tests still pass with `NewApp(ModeAll)`
- [ ] 📍 **Checkpoint B** — `go test ./... -count=1` all pass, `./bin/rapidgo serve` (default mode) works identically to before, `./bin/rapidgo serve --mode=api` starts API-only

---

## Phase C — Multi-Port Serving

> Multiple HTTP servers on different ports from one binary.

- [ ] **C.1** — Add `ListenAndServeMulti()` to `core/server/server.go`
  - [ ] `ServiceConfig` struct with `Name` + `Config`
  - [ ] Start each server in a goroutine
  - [ ] Wait for signal or any server error
  - [ ] Graceful shutdown for ALL servers
- [ ] **C.2** — Add multi-port logic to `core/cli/serve.go`
  - [ ] `serveSingle()` — single server (backward compatible path)
  - [ ] `serveMulti()` — per-mode routers on per-mode ports
  - [ ] `allSamePort()` — detect if multi-port is needed
  - [ ] `resolvePort()` / `resolvePortForMode()` — port resolution with fallback chain
  - [ ] `applyRoutesForMode()` — register routes on a fresh router per mode
- [ ] **C.3** — Add `core/server/server_test.go` tests for `ListenAndServeMulti`
  - [ ] Test multiple servers start on different ports
  - [ ] Test graceful shutdown stops all servers
- [ ] **C.4** — Add multi-port banner output
  - [ ] Show each service name and port in startup banner
- [ ] 📍 **Checkpoint C** — `RAPIDGO_MODE=api,ws API_PORT=8081 WS_PORT=8082 ./bin/rapidgo serve` starts two servers on separate ports. `go test ./... -count=1` all pass.

---

## Ship 🚀

- [ ] All phases complete
- [ ] All checkpoints verified
- [ ] `go test ./... -count=1` — all packages pass
- [ ] Default mode (`all`) behavior identical to pre-feature
- [ ] Final commit with descriptive message
- [ ] Push to feature branch
- [ ] Merge to `main`
- [ ] Push `main`
- [ ] **Keep the feature branch**
- [ ] Create review doc → `56-service-mode-review.md`
