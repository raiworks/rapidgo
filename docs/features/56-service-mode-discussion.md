# 💬 Discussion: Service Mode Architecture

> **Feature**: `56` — Service Mode Architecture
> **Status**: � COMPLETE
> **Branch**: `feature/56-service-mode`
> **Depends On**: #55 (Framework Rename — ✅ Complete)
> **Date Started**: 2026-03-07
> **Date Completed**: 2026-03-07

---

## Summary

Enable the framework to run in different service modes — monolith (all services in one binary), split microservices (each service independently), or combined subsets (e.g., API + WebSocket together). Today the framework boots everything on one port with no way to customize. The underlying Go/Gin/GORM stack fully supports service splitting; the blockers are in the bootstrapping layer, not the core components.

---

## Functional Requirements

- As a developer, I want to run the full app as a monolith so that I can develop and deploy simply during early stages
- As a developer, I want to run only the API service so that I can scale it independently from the web frontend
- As a developer, I want to run only the web SSR service so that I can deploy it on different infrastructure
- As a developer, I want to run only the WebSocket service so that I can isolate real-time traffic
- As a developer, I want to run combined subsets (e.g., API + WebSocket) so that I can co-locate related services
- As a developer, I want to skip providers I don't need (e.g., no DB for a stateless proxy) so that startup is fast and dependencies are minimal
- As a developer, I want different services on different ports so that I can route traffic independently
- As a DevOps engineer, I want to control service mode via environment variable or CLI flag so that deployment configuration is infrastructure-driven

## Current State / Reference

### What Exists

The app starts via a hardcoded bootstrap sequence:

```
cmd/main.go → cli.Execute() → serve command → NewApp()
    ↓
NewApp() registers ALL 6 providers (no choice):
    1. ConfigProvider      → loads .env
    2. LoggerProvider      → sets up slog
    3. DatabaseProvider    → registers "db" singleton (LAZY — connects on first resolve)
    4. SessionProvider     → registers "session" (LAZY — REQUIRES "db" on first resolve)
    5. MiddlewareProvider  → registers global middleware maps
    6. RouterProvider      → creates ONE router, registers ALL routes
    ↓
application.Boot() → calls Boot() on all 6 in order
    ↓
serve command → extracts "router" → server.ListenAndServe(:8080)
```

**Key files**:
- `core/cli/root.go` — `NewApp()` with hardcoded provider list
- `app/providers/session_provider.go` — panics without DB (`MustMake[*gorm.DB]`) — note: lazy singleton, panic occurs on first resolve, not at registration
- `app/providers/router_provider.go` — always calls both `RegisterWeb()` and `RegisterAPI()`
- `core/cli/serve.go` — single server on single port
- `core/middleware/registry.go` — package-level global maps
- `core/router/named.go` — package-level global map

### Provider Dependency Graph

```
ConfigProvider ──────────────────────────────────────┐
    │                                                 │
    ├──→ LoggerProvider (reads LOG_* from config)     │
    │                                                 │
    ├──→ DatabaseProvider (reads DB_* from config)    │
    │         │                                       │
    │         └──→ SessionProvider (LAZY hard dep on DB — panics on first Make("session"), not at boot)
    │                                                 │
    ├──→ MiddlewareProvider (independent)             │
    │                                                 │
    └──→ RouterProvider                               │
              ├── RegisterWeb()  ← always called      │
              ├── RegisterAPI()  ← always called      │
              └── health.Routes() ← if DB exists      │
```

### What Works Well

- The container/provider pattern is the right foundation — already decoupled
- Gin engines are independent route trees — multiple engines are fully supported
- GORM connections are plain structs — no global state
- The server package supports arbitrary `http.Handler` — not locked to one engine
- WebSocket is just an HTTP upgrade on a route — can go on any engine

### What Needs Improvement

- **Cannot skip providers** — all 6 always register, even if unused
- **SessionProvider panics without DB** — hard dependency makes stateless services impossible
- **Routes always load together** — no way to run API-only or Web-only
- **Single port** — everything on `:8080`, no multi-port capability
- **Global state** — middleware registry and named routes use package-level maps that would collide in multi-service process
- **No worker mode** — no CLI command for background job processing

## Proposed Approach

### Service Mode Concept

Add a service mode system controlled by environment variable (`RAPIDGO_MODE`) or CLI flag (`--mode`):

```
RAPIDGO_MODE=all        →  Web + API + WebSocket (monolith — current behavior)
RAPIDGO_MODE=web        →  Web SSR only (templates, static files)
RAPIDGO_MODE=api        →  API only (JSON endpoints)
RAPIDGO_MODE=ws         →  WebSocket only
RAPIDGO_MODE=worker     →  Background job processor only (no HTTP)
RAPIDGO_MODE=api,ws     →  API + WebSocket combined
RAPIDGO_MODE=web,api    →  Web + API (no WebSocket)
```

### CLI Interface

```
RapidGo serve                     →  starts in RAPIDGO_MODE (default: all)
RapidGo serve --mode=api          →  API-only service on API_PORT
RapidGo serve --mode=web          →  Web-only service on WEB_PORT
RapidGo serve --mode=api,ws       →  API + WebSocket on separate ports
RapidGo worker                    →  Job processing (no HTTP)
```

### Implementation Phases

The work breaks into 5 phases, each independently shippable:

| Phase | What | Impact | Effort |
|---|---|---|---|
| **Phase 1** | Optional Providers | Unblocks everything — providers can be skipped | Small |
| **Phase 2** | Service Mode Flag | Enables API-only / Web-only / WS-only | Small |
| **Phase 3** | Multi-Port Serving | Different services on different ports from one binary | Medium |
| **Phase 4** | Remove Global State | Multiple services in same process without conflicts | Medium |
| **Phase 5** | Worker/Queue Mode | Background job processing (`RapidGo worker`) | Large |

**Phases 1-2** should be done proactively — they're small and unblock the core capability.  
**Phase 3** is valuable when actually splitting services.  
**Phase 4** is only needed for multi-service-in-same-process scenarios.  
**Phase 5** is a separate subsystem (queue/job) — do when background processing is needed.

### Technology Validation

All four core capabilities are confirmed possible:

1. **Gin subset routes** — Gin engines are independent route trees. Multiple engines with different routes: ✅
2. **Multiple engines in one process** — Each engine on different port via goroutines: ✅
3. **WebSocket on separate port** — WebSocket is just HTTP upgrade, any engine: ✅
4. **Isolated containers** — Container is a plain struct, create multiple: ✅

### Deployment Scenarios

| Scenario | Mode | Description |
|---|---|---|
| **A — Monolith** | `RAPIDGO_MODE=all` | Everything on one port. MVP, small teams, simple deploys |
| **B — Split Services** | Separate binaries each with own `RAPIDGO_MODE` | API, Web, WS as independent services. Scale independently |
| **C — Combined Subsets** | `RAPIDGO_MODE=web,api` | Flexible grouping based on traffic patterns |
| **D — Separate Binaries** | Build tags + mode | Each binary only contains needed code. Smaller images |

### What the Framework Should NOT Build

Infrastructure concerns better solved by external tools:

| Concern | Why not in framework | Better solution |
|---|---|---|
| Service discovery | Infrastructure layer | Kubernetes / Consul / DNS |
| Load balancing | Infrastructure layer | Nginx / HAProxy / K8s Service |
| Service-to-service auth | Infrastructure layer | mTLS / service mesh (Istio) |
| Distributed tracing | Observability layer | OpenTelemetry / Jaeger |
| Container orchestration | DevOps layer | Docker Compose / Kubernetes |
| API Gateway | Infrastructure layer | Kong / Traefik / Caddy |

The framework should focus on making the **application code** splittable, not reimplementing infrastructure.

## Edge Cases & Risks

- [ ] **Mode validation** — Invalid mode strings (e.g., `RAPIDGO_MODE=invalid`) must fail fast with a clear error, not silently start nothing
- [ ] **Provider ordering** — When providers are optional, remaining providers must still boot in correct dependency order
- [ ] **Session without DB** — If mode=api and API needs sessions (e.g., for auth), but SessionProvider was skipped, this must be a clear error at boot, not a runtime panic
- [ ] **Health check per mode** — Each service mode should expose its own health endpoint. An API-only service shouldn't report web template status
- [ ] **Shared state across modes** — When running `RAPIDGO_MODE=all`, web and API routes share the same container/DB connection. When split, they don't. Code that assumes shared state will break when split
- [ ] **Config conflicts** — `APP_PORT` vs `WEB_PORT` vs `API_PORT` — clear naming needed to avoid confusion
- [ ] **Graceful shutdown per mode** — Multi-port serving needs graceful shutdown for ALL servers, not just one
- [ ] **Testing** — Tests currently assume all providers are loaded. Mode-aware tests need thought
- [ ] **Backward compatibility** — `RAPIDGO_MODE` unset must behave exactly as today (all services, single port)

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Feature #55 — Framework Rename | Feature | ✅ Done |
| All 41 core features complete | Feature | ✅ Done |
| Service Container (#05) | Feature | ✅ Done |
| Service Providers (#06) | Feature | ✅ Done |
| Router & Routing (#07) | Feature | ✅ Done |
| CLI Foundation (#10) | Feature | ✅ Done |
| Graceful Shutdown (#37) | Feature | ✅ Done |

## Open Questions

- [x] **Q1**: Should Phases 1-5 each be separate Mastery features (with their own doc sets), or is this one feature with phased implementation? **Decision**: Phases 1-3 = Feature #56 (this feature). Phase 4 (Remove Global State) = Feature #57. Phase 5 (Worker/Queue) = Feature #42 (already in roadmap).
- [x] **Q2**: Should the `--mode` flag override `RAPIDGO_MODE` env var, or vice versa? **Decision**: CLI flag overrides env var. Standard 12-factor precedence: flag > env > default. Matches existing `--port` / `APP_PORT` pattern.
- [x] **Q3**: For multi-port mode, should each service get its own container instance or share one? **Decision**: Shared container with per-service routers. One DB pool, one config, separate Gin engines per port. Simpler and less memory.
- [x] **Q4**: Should we support `RAPIDGO_MODE=api` as a build-time optimization (via Go build tags) to exclude web template code from the binary? Or runtime-only? **Decision**: Runtime-only for #56. Build tags are a future optimization — runtime mode selection is simpler, testable, and covers all deployment scenarios.
- [x] **Q5**: Does Phase 5 (Worker/Queue) overlap with roadmap item #42 (Queue Workers / Background Jobs)? Should they be merged? **Decision**: Yes — Phase 5 IS Feature #42. No duplication.
- [x] **Q6**: Should the `routes/` directory split into `routes/web.go`, `routes/api.go`, `routes/ws.go`? It already has `web.go` and `api.go` — do we just add `ws.go`? **Decision**: Yes, add `routes/ws.go` with `RegisterWS()` following the same pattern. RouterProvider calls it conditionally based on mode.
- [x] **Q7**: When running in `web` mode only, should API middleware (like JSON content-type enforcement) be excluded from the middleware stack entirely? **Decision**: Yes. Each mode loads only its relevant middleware. Web mode skips JSON enforcement. API mode skips CSRF. MiddlewareProvider becomes mode-aware.

## Comparison with Other Frameworks

| Capability | Laravel | Spring Boot | NestJS | Current | Proposed |
|---|---|---|---|---|---|
| Monolith mode | ✅ | ✅ | ✅ | ✅ | ✅ |
| Optional providers | ✅ | ✅ (profiles) | ✅ (modules) | ❌ | ✅ Phase 1 |
| Service modes | ✅ (artisan) | ✅ (profiles) | ✅ (modules) | ❌ | ✅ Phase 2 |
| Multi-port | ❌ | ✅ | ✅ | ❌ | ✅ Phase 3 |
| Background workers | ✅ (queues) | ✅ | ✅ (bull) | ❌ | 🔮 Feature #42 |
| Microservice toolkit | ⚠️ Lumen | ✅ (cloud) | ✅ | ❌ | ✅ Phase 2-3 |

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-07 | Confirmed Go/Gin/GORM stack supports service splitting | Multiple Gin engines, independent containers, separate ports — all proven |
| 2026-03-07 | Framework should NOT build infra (service discovery, load balancing, etc.) | Better solved by Kubernetes/Consul/Nginx — framework focuses on app code splittability |
| 2026-03-07 | 5-phase implementation plan proposed | Each phase independently shippable, ordered by impact/effort ratio |
| 2026-03-07 | Depends on #55 (Rename) completing first | Avoid renaming new service mode code after the fact |
| 2026-03-07 | Feature scope: Phases 1-3 only (#56). Phase 4 → #57, Phase 5 → #42 | Keep features focused and independently shippable |
| 2026-03-07 | CLI flag overrides env var | Standard 12-factor precedence (flag > env > default) |
| 2026-03-07 | Shared container, per-service routers for multi-port | One DB pool, separate route trees, simpler architecture |
| 2026-03-07 | Runtime-only mode selection (no build tags) | Simpler, testable, build tags deferred as future optimization |
| 2026-03-07 | Phase 5 (Worker/Queue) = Feature #42 | Already in roadmap, no duplication |
| 2026-03-07 | Add `routes/ws.go` with `RegisterWS()` | Follows existing web.go / api.go pattern |
| 2026-03-07 | Mode-aware middleware loading | Each mode loads only relevant middleware stack |

---

## Discussion Complete ✅

<!-- Fill this section when ALL open questions are resolved -->

**Summary**: Enable service mode architecture (Phases 1-3) — optional providers, `RAPIDGO_MODE` env var / `--mode` CLI flag for selecting web/api/ws/all modes, and multi-port serving with per-service routers sharing one container. Default mode `all` preserves full backward compatibility. Phase 4 (global state removal) deferred to #57, Phase 5 (worker/queue) deferred to #42.
**Completed**: 2026-03-07
**Next**: Create architecture doc → `56-service-mode-architecture.md`
