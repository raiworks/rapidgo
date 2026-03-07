# ✅ Tasks: Service Mode

> **Feature**: `56` — Service Mode Architecture
> **Architecture**: [`56-service-mode-architecture.md`](56-service-mode-architecture.md)
> **Branch**: `feature/56-service-mode`
> **Status**: � COMPLETE
> **Progress**: 14/14 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [x] Feature branch created from latest `main`
- [x] Dependent features are merged to `main` (#55 complete)
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — Service Mode Infrastructure

> New `core/service/` package with Mode type, parsing, and validation.

- [x] **A.1** — Create `core/service/mode.go`
  - [x] `Mode` type as `uint8` bitmask
  - [x] Constants: `ModeWeb`, `ModeAPI`, `ModeWS`, `ModeAll`
  - [x] `modeNames` map for string-to-Mode lookup
  - [x] `ParseMode(s string) (Mode, error)` — comma-separated parsing with validation
  - [x] `Mode.Has(flag Mode) bool` — bitmask check
  - [x] `Mode.String() string` — human-readable output
  - [x] `Mode.Services() []Mode` — list of individual active modes
  - [x] `Mode.PortEnvKey() string` — env var name for mode-specific port
- [x] **A.2** — Create `core/service/mode_test.go`
  - [x] Test `ParseMode` with valid inputs: `"all"`, `"web"`, `"api"`, `"ws"`, `"api,ws"`, `"web,api,ws"`
  - [x] Test `ParseMode` with invalid inputs: `""`, `"invalid"`, `"api,invalid"`, `"worker"`
  - [x] Test `Mode.Has()` combinations
  - [x] Test `Mode.String()` output
  - [x] Test `Mode.Services()` returns correct list
  - [x] Test `Mode.PortEnvKey()` returns correct env var names
- [x] **A.3** — Create `routes/ws.go`
  - [x] `RegisterWS(r *router.Router)` — empty placeholder with comment
- [x] 📍 **Checkpoint A** — `go test ./core/service/... -count=1` passes, `go build ./...` compiles

---

## Phase B — Mode-Aware Bootstrap

> Modify provider registration and route loading based on active mode.

- [x] **B.1** — Update `core/cli/root.go`
  - [x] Change `NewApp()` signature to `NewApp(mode service.Mode) *app.App`
  - [x] Conditional DB provider: only if `mode.Has(ModeWeb|ModeAPI|ModeWS)`
  - [x] Conditional Session provider: only if `mode.Has(ModeWeb)`
  - [x] Pass `mode` to `MiddlewareProvider{Mode: mode}` and `RouterProvider{Mode: mode}`
  - [x] Add import for `core/service`
- [x] **B.2** — Update `app/providers/router_provider.go`
  - [x] Add `Mode service.Mode` field to `RouterProvider` struct
  - [x] Template/static setup conditional on `Mode.Has(ModeWeb)`
  - [x] `routes.RegisterWeb(r)` conditional on `Mode.Has(ModeWeb)`
  - [x] `routes.RegisterAPI(r)` conditional on `Mode.Has(ModeAPI)`
  - [x] `routes.RegisterWS(r)` conditional on `Mode.Has(ModeWS)`
  - [x] Health check remains: if `c.Has("db")` (unchanged)
- [x] **B.3** — Update `app/providers/middleware_provider.go`
  - [x] Add `Mode service.Mode` field to `MiddlewareProvider` struct
  - [x] CSRF registration conditional on `Mode.Has(ModeWeb)`
  - [x] All other middleware registered unconditionally (recovery, requestid, cors, error_handler, auth, ratelimit)
- [x] **B.4** — Update `core/cli/serve.go`
  - [x] Add `--mode` / `-m` flag (`serveMode` variable)
  - [x] Mode resolution: `--mode` flag > `RAPIDGO_MODE` env > `"all"` default
  - [x] Call `NewApp(mode)` instead of `NewApp()`
  - [x] Add mode to startup banner output
  - [x] Change `Run` to `RunE` for error return from `ParseMode`
- [x] **B.5** — Update `.env`
  - [x] Add `RAPIDGO_MODE=all` with comment
  - [x] Add `WEB_PORT=`, `API_PORT=`, `WS_PORT=` (commented out, with defaults explanation)
- [x] **B.6** — Update all other `NewApp()` callers
  - [x] `core/cli/migrate.go` — pass `service.ModeAll`
  - [x] `core/cli/migrate_rollback.go` — pass `service.ModeAll`
  - [x] `core/cli/migrate_status.go` — pass `service.ModeAll`
  - [x] `core/cli/seed.go` — pass `service.ModeAll`
  - [x] Any other commands calling `NewApp()`
- [x] **B.7** — Update tests for modified files
  - [x] `app/providers/providers_test.go` — pass `Mode` field in provider constructors
  - [x] `core/cli/cli_test.go` — verify tests still pass with `NewApp(ModeAll)`
- [x] 📍 **Checkpoint B** — `go test ./... -count=1` all pass, `./bin/rapidgo serve` (default mode) works identically to before, `./bin/rapidgo serve --mode=api` starts API-only

---

## Phase C — Multi-Port Serving

> Multiple HTTP servers on different ports from one binary.

- [x] **C.1** — Add `ListenAndServeMulti()` to `core/server/server.go`
  - [x] `ServiceConfig` struct with `Name` + `Config`
  - [x] Start each server in a goroutine
  - [x] Wait for signal or any server error
  - [x] Graceful shutdown for ALL servers
- [x] **C.2** — Add multi-port logic to `core/cli/serve.go`
  - [x] `serveSingle()` — single server (backward compatible path)
  - [x] `serveMulti()` — per-mode routers on per-mode ports
  - [x] `allSamePort()` — detect if multi-port is needed
  - [x] `resolvePort()` / `resolvePortForMode()` — port resolution with fallback chain
  - [x] `applyRoutesForMode()` — register routes on a fresh router per mode
- [x] **C.3** — Add `core/server/server_test.go` tests for `ListenAndServeMulti`
  - [x] Test multiple servers start on different ports
  - [x] Test graceful shutdown stops all servers
- [x] **C.4** — Add multi-port banner output
  - [x] Show each service name and port in startup banner
- [x] 📍 **Checkpoint C** — `RAPIDGO_MODE=api,ws API_PORT=8081 WS_PORT=8082 ./bin/rapidgo serve` starts two servers on separate ports. `go test ./... -count=1` all pass.

---

## Ship 🚀

- [x] All phases complete
- [x] All checkpoints verified
- [x] `go test ./... -count=1` — all packages pass
- [x] Default mode (`all`) behavior identical to pre-feature
- [x] Final commit with descriptive message
- [x] Push to feature branch
- [x] Merge to `main`
- [x] Push `main`
- [x] **Keep the feature branch**
- [x] Create review doc → `56-service-mode-review.md`
