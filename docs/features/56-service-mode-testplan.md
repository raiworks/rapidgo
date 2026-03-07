# 🧪 Test Plan: Service Mode

> **Feature**: `56` — Service Mode Architecture
> **Architecture**: [`56-service-mode-architecture.md`](56-service-mode-architecture.md)
> **Coverage Target**: All mode parsing, provider selection, route registration, and multi-port serving

---

## Pre-Conditions

- [ ] All 30 existing test packages pass before implementation begins
- [ ] No stale references or compilation errors in codebase

---

## Test Cases

### TC-01: ParseMode — Valid Inputs

| Field | Value |
|---|---|
| **ID** | TC-01 |
| **Title** | ParseMode accepts all valid mode strings |
| **Type** | Unit |
| **File** | `core/service/mode_test.go` |
| **Steps** | Call `ParseMode()` with: `"all"`, `"web"`, `"api"`, `"ws"`, `"api,ws"`, `"web,api"`, `"web,api,ws"`, `" api , ws "` (whitespace) |
| **Expected** | Returns correct Mode bitmask, nil error for each |
| **Verify** | `ModeAll` has all three flags. `"api,ws"` has `ModeAPI|ModeWS`. Whitespace is trimmed. |

---

### TC-02: ParseMode — Invalid Inputs

| Field | Value |
|---|---|
| **ID** | TC-02 |
| **Title** | ParseMode rejects invalid mode strings |
| **Type** | Unit |
| **File** | `core/service/mode_test.go` |
| **Steps** | Call `ParseMode()` with: `""`, `"invalid"`, `"api,invalid"`, `"worker"`, `","` |
| **Expected** | Returns error for each. Error message includes valid options. |
| **Verify** | Error is non-nil. Mode value is 0. |

---

### TC-03: Mode Bitmask Operations

| Field | Value |
|---|---|
| **ID** | TC-03 |
| **Title** | Mode.Has(), String(), Services(), PortEnvKey() work correctly |
| **Type** | Unit |
| **File** | `core/service/mode_test.go` |
| **Steps** | 1. `ModeAPI.Has(ModeAPI)` → true. 2. `ModeAPI.Has(ModeWeb)` → false. 3. `ModeAll.Has(ModeWeb)` → true. 4. `ModeAll.String()` → `"all"`. 5. `(ModeAPI\|ModeWS).String()` → `"api,ws"`. 6. `(ModeAPI\|ModeWS).Services()` → `[ModeAPI, ModeWS]`. 7. `ModeAPI.PortEnvKey()` → `"API_PORT"`. 8. `ModeWeb.PortEnvKey()` → `"WEB_PORT"`. |
| **Expected** | All assertions pass |

---

### TC-04: Default Mode Backward Compatibility

| Field | Value |
|---|---|
| **ID** | TC-04 |
| **Title** | No RAPIDGO_MODE or --mode flag defaults to "all" |
| **Type** | Integration |
| **File** | `core/cli/cli_test.go` |
| **Steps** | 1. Unset `RAPIDGO_MODE`. 2. Call serve without `--mode`. 3. Verify all providers load. 4. Verify all routes registered (web + API + WS). |
| **Expected** | Behavior identical to pre-feature. All 6 providers boot. Both `RegisterWeb()` and `RegisterAPI()` called. |

---

### TC-05: CLI --mode Overrides Env Var

| Field | Value |
|---|---|
| **ID** | TC-05 |
| **Title** | --mode flag takes precedence over RAPIDGO_MODE env var |
| **Type** | Unit |
| **File** | `core/cli/cli_test.go` |
| **Steps** | 1. Set `RAPIDGO_MODE=all`. 2. Run serve with `--mode=api`. 3. Verify API-only mode is active. |
| **Expected** | Mode resolves to `ModeAPI`, not `ModeAll`. |

---

### TC-06: API-Only Mode — Provider Selection

| Field | Value |
|---|---|
| **ID** | TC-06 |
| **Title** | mode=api loads correct providers and skips Session |
| **Type** | Integration |
| **File** | `app/providers/providers_test.go` |
| **Steps** | 1. Call `NewApp(service.ModeAPI)`. 2. Verify container has `"db"` binding. 3. Verify container does NOT have `"session"` binding. 4. Verify container has `"router"` binding. |
| **Expected** | DB registered, Session NOT registered, Router registered. |

---

### TC-07: Web-Only Mode — Session Loaded

| Field | Value |
|---|---|
| **ID** | TC-07 |
| **Title** | mode=web loads Session provider |
| **Type** | Integration |
| **File** | `app/providers/providers_test.go` |
| **Steps** | 1. Call `NewApp(service.ModeWeb)`. 2. Verify container has `"session"` binding. 3. Verify container has `"db"` binding (session depends on it). |
| **Expected** | Both DB and Session registered. |

---

### TC-08: Route Registration Per Mode

| Field | Value |
|---|---|
| **ID** | TC-08 |
| **Title** | Only mode-relevant routes are registered |
| **Type** | Integration |
| **File** | `app/providers/providers_test.go` or `core/cli/cli_test.go` |
| **Steps** | 1. Boot with `ModeAPI` → check router has `/api/*` routes, no `/` web route. 2. Boot with `ModeWeb` → check router has `/` web route, no `/api/*` routes. 3. Boot with `ModeAll` → check router has both. |
| **Expected** | Routes match the active mode exactly. |

---

### TC-09: Mode-Aware Middleware

| Field | Value |
|---|---|
| **ID** | TC-09 |
| **Title** | CSRF middleware only registered in web mode |
| **Type** | Unit |
| **File** | `core/middleware/middleware_test.go` |
| **Steps** | 1. Reset middleware registry. 2. Boot `MiddlewareProvider{Mode: ModeAPI}`. 3. Check `"csrf"` alias does NOT exist. 4. Boot `MiddlewareProvider{Mode: ModeWeb}`. 5. Check `"csrf"` alias exists. |
| **Expected** | CSRF absent in API mode, present in Web mode. |

---

### TC-10: Multi-Port — Separate Servers

| Field | Value |
|---|---|
| **ID** | TC-10 |
| **Title** | Combined mode with different ports starts multiple servers |
| **Type** | Integration |
| **File** | `core/server/server_test.go` |
| **Steps** | 1. Create two `ServiceConfig` (api on :0, ws on :0 — random ports). 2. Start `ListenAndServeMulti` in goroutine. 3. HTTP GET both ports → get responses. 4. Send SIGINT → both shut down. |
| **Expected** | Both servers respond. Both shut down gracefully. |
| **Notes** | Use `:0` for random port to avoid conflicts in CI. `ListenAndServeMulti` uses `signal.NotifyContext` internally — existing server tests avoid calling `ListenAndServe` directly for this reason and test `http.Server` start/stop directly. Follow the same pattern here, or accept a `context.Context` parameter for testability. |

---

### TC-11: Multi-Port — Same Port Falls Back to Single

| Field | Value |
|---|---|
| **ID** | TC-11 |
| **Title** | Combined mode with same port uses single server |
| **Type** | Unit |
| **File** | `core/cli/cli_test.go` |
| **Steps** | 1. Set `API_PORT=8080`, `WS_PORT=8080`. 2. Parse mode `api,ws`. 3. Call `allSamePort()`. |
| **Expected** | `allSamePort()` returns `true`. Single server path is used. |

---

### TC-12: Invalid Mode Fails Fast

| Field | Value |
|---|---|
| **ID** | TC-12 |
| **Title** | Invalid --mode value returns error, does not start server |
| **Type** | Unit |
| **File** | `core/cli/cli_test.go` |
| **Steps** | 1. Run serve with `--mode=invalid`. 2. Check error returned. |
| **Expected** | Error contains "invalid service mode". Server does not start. |

---

### TC-13: Full Regression

| Field | Value |
|---|---|
| **ID** | TC-13 |
| **Title** | All existing tests pass with no regressions |
| **Type** | Regression |
| **File** | All packages |
| **Steps** | `go test ./... -count=1` |
| **Expected** | All 30+ packages pass. Zero failures. |

---

## Test Summary

| TC | Title | Type | Status |
|---|---|---|---|
| TC-01 | ParseMode valid inputs | Unit | ⬜ |
| TC-02 | ParseMode invalid inputs | Unit | ⬜ |
| TC-03 | Mode bitmask operations | Unit | ⬜ |
| TC-04 | Default mode backward compat | Integration | ⬜ |
| TC-05 | CLI --mode overrides env var | Unit | ⬜ |
| TC-06 | API-only provider selection | Integration | ⬜ |
| TC-07 | Web-only Session loaded | Integration | ⬜ |
| TC-08 | Route registration per mode | Integration | ⬜ |
| TC-09 | Mode-aware middleware (CSRF) | Unit | ⬜ |
| TC-10 | Multi-port separate servers | Integration | ⬜ |
| TC-11 | Same port falls back to single | Unit | ⬜ |
| TC-12 | Invalid mode fails fast | Unit | ⬜ |
| TC-13 | Full regression | Regression | ⬜ |

---

## Known Limitations

| # | Limitation | Mitigation |
|---|---|---|
| 1 | Global middleware/named-route maps may collide in multi-port | Deferred to Feature #57 (Remove Global State) |
| 2 | Worker mode not included | Deferred to Feature #42 (Queue Workers) |
| 3 | Build tags for binary optimization not included | Deferred — runtime-only for this feature |
