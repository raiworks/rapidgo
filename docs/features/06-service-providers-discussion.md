# 💬 Discussion: Service Providers

> **Feature**: `06` — Service Providers
> **Status**: 🟢 COMPLETE
> **Branch**: `feature/06-service-providers`
> **Depends On**: #05 (Service Container ✅)
> **Date Started**: 2026-03-06
> **Date Completed**: 2026-03-06

---

## Summary

Implement the built-in service providers that bootstrap existing framework features (#02 Config, #03 Logger, #04 Errors) through the Provider interface, update `cmd/main.go` to use the `App` bootstrap pattern, and establish the `app/providers/` directory as the home for all providers.

The blueprint's full set of built-in providers (Database, Session, Cache, Mail, Event) depend on subsystems not yet implemented — those providers will be created in their respective feature branches (#09, #20, #29, #32, #34). Feature #06 bootstraps what exists today.

---

## Functional Requirements

- As a **framework developer**, I want a `ConfigProvider` so that `config.Load()` is called during registration and the config module is available as a container service
- As a **framework developer**, I want a `LoggerProvider` so that `logger.Setup()` is called during boot (after config is loaded) and the logger is initialized through the provider lifecycle
- As a **framework developer**, I want `cmd/main.go` updated to use `app.New()` → `app.Register()` → `app.Boot()` so that the application follows the provider bootstrap pattern from the blueprint
- As a **framework developer**, I want the `app/providers/` directory to contain real provider implementations so that it serves as the canonical location for all providers going forward

## Current State / Reference

### What Exists
- **Container** (#05 ✅): `Container`, `Provider` interface, `App` struct — all tested and working
- **Configuration** (#02 ✅): `config.Load()`, `config.Env()`, `config.EnvInt()`, `config.EnvBool()`, `config.AppEnv()`, `config.IsDebug()`, `config.IsProduction()`, `config.IsDevelopment()`, `config.IsTesting()`
- **Logging** (#03 ✅): `logger.Setup()`, `logger.Close()` — configures `slog` based on env vars
- **Error handling** (#04 ✅): `core/errors` package — `AppError` struct, constructors, `ErrorResponse()`
- **`cmd/main.go`**: Currently calls `config.Load()` and `logger.Setup()` directly — no App/provider usage
- **`app/providers/`**: Empty directory with `.gitkeep`

### What Works Well
- Provider interface and App struct are implemented and tested in Feature #05
- Config and Logger packages have clean, callable APIs (`config.Load()`, `logger.Setup()`)
- The blueprint shows providers registering services via `Singleton`/`Instance` — same pattern applies here

### What Needs Improvement
- `cmd/main.go` bypasses the provider lifecycle entirely — calls config/logger directly
- No providers exist yet — `app/providers/` is empty
- Framework features aren't registered in the container — can't be resolved via `container.Make()`

## Proposed Approach

Create two concrete providers and update `main.go`:

1. **`app/providers/config_provider.go`** — `ConfigProvider`
   - `Register()`: calls `config.Load()`, registers a config accessor as `"config"` instance
   - `Boot()`: empty — config is fully available after `Load()`

2. **`app/providers/logger_provider.go`** — `LoggerProvider`
   - `Register()`: empty — logger needs config to be loaded first
   - `Boot()`: calls `logger.Setup()` — runs after all providers are registered, config is guaranteed available

3. **`cmd/main.go`** — Updated to use App bootstrap
   - Creates `app.New()`
   - Registers `ConfigProvider` (first — everything depends on config)
   - Registers `LoggerProvider`
   - Calls `app.Boot()`
   - Banner and server init remain

**NOT in scope for this feature**:
- DatabaseProvider — Feature #09
- SessionProvider — Feature #20
- CacheProvider — Feature #32
- MailProvider — Feature #29
- EventProvider — Feature #34
- Custom provider CLI generator (`make:provider`) — Feature #41

## Edge Cases & Risks

- [x] Config must load before logger — enforced by: ConfigProvider registered first, Logger uses Boot() not Register()
- [x] `config.Load()` is idempotent — safe if called more than once
- [x] `logger.Setup()` in Boot() guarantees config is loaded — all Register() calls complete before any Boot()
- [x] Provider order matters — documented in main.go with comments

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Feature #05 — Service Container | Feature | ✅ Done |
| Feature #02 — Configuration | Feature | ✅ Done |
| Feature #03 — Logging | Feature | ✅ Done |
| `config` package | Internal | ✅ Available |
| `logger` package | Internal | ✅ Available |

## Open Questions

_All resolved during discussion._

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-06 | Only create ConfigProvider and LoggerProvider | Other providers depend on unimplemented subsystems; build them in their respective features |
| 2026-03-06 | ConfigProvider calls `config.Load()` in Register() | Config must be the first thing loaded; other providers may read config during their Register() |
| 2026-03-06 | LoggerProvider calls `logger.Setup()` in Boot() | Logger depends on config values (LOG_LEVEL, LOG_FORMAT, LOG_OUTPUT); Boot() guarantees config is loaded |
| 2026-03-06 | Update `cmd/main.go` to App bootstrap | Establishes the blueprint's bootstrap pattern; makes the framework usable as intended |
| 2026-03-06 | No ErrorProvider needed | Error package is stateless utility functions — no service to register in container |

## Discussion Complete ✅

**Summary**: Feature #06 creates ConfigProvider and LoggerProvider in `app/providers/`, updates `cmd/main.go` to use App bootstrap. Scoped to existing features only — future providers built in their own features.
**Completed**: 2026-03-06
**Next**: Create architecture doc → `06-service-providers-architecture.md`
