# Changelog

All notable changes to the **RapidGo** framework will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.8.0] - 2026-04-20

### Added
- **New `core/pubsub` package** — Cross-process publish/subscribe messaging over Redis
  - `Publisher` interface with `NewRedisPublisher(client)` — wraps Redis PUBLISH
  - `Subscriber` interface with `NewRedisSubscriber(client, opts)` — auto-reconnect with exponential backoff
  - `SubscriberOptions` for configurable `MinBackoff` (500ms), `MaxBackoff` (30s), `PingEvery` (30s), `Logger`
  - `Handler` type: `func(ctx context.Context, channel, payload string)`
  - Panic recovery in handler dispatch — panicking handlers are logged, not fatal
  - 8 tests using miniredis covering publish/receive, multi-channel, context cancel, reconnect, error, defaults, empty channels, handler panic
- New documentation: `docs/framework/infrastructure/pubsub.md` — API reference, usage examples, reconnect semantics, limitations

### Changed
- `docs/framework/infrastructure/caching.md` — bumped to v0.2.0, added cross-reference to pubsub.md in §7.1

## [2.7.3] - 2026-04-20

### Fixed
- **REDIS_DB env var silently ignored** — `core/cache` and `core/session` now honor `REDIS_DB` when building Redis clients. Previously all clients defaulted to DB 0 regardless of configuration.

### Added
- `cache.NewRedisClient(dbOverride *int)` — exported helper that builds a `*redis.Client` from env vars (`REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`) with optional DB override for multi-client registration
- New optional Redis env vars with sensible defaults: `REDIS_POOL_SIZE` (10), `REDIS_DIAL_TIMEOUT` (5s), `REDIS_READ_TIMEOUT` (3s), `REDIS_WRITE_TIMEOUT` (3s)
- 12 tests for `NewRedisClient` covering DB selection, override precedence, invalid input, pool/timeout config

### Changed
- `core/cache/cache.go` — `newRedisClient()` now delegates to `NewRedisClient(nil)`
- `core/session/factory.go` — Redis branch uses `cache.NewRedisClient(nil)` instead of inline client creation

## [2.7.2] - 2026-03-15

### Fixed
- **BUG-03**: `applyRoutesForMode()` no longer overwrites custom `FuncMap` on routers that already have templates loaded — adds `TemplatesLoaded()` guard in `core/router`

### Added
- `container.SafeSingleton()` — catches factory panics in local/dev mode and logs a warning instead of crashing. Production panics propagate as normal. (`core/container`)
- `seeders.RunAll()` now logs `slog.Warn` when the seeder registry is empty (previously silent no-op)

## [2.7.1] - 2026-03-15

### Added
- `config.MustLoadConfig[T]()` — like `LoadConfig[T]()` but panics on error. Use at startup where missing/invalid config is fatal. (`core/config`)
- 20 exported error code constants (`CodeNotFound`, `CodeBadRequest`, `CodeServiceUnavailable`, `CodeTimeout`, `CodeRateLimited`, etc.) in `core/errors`
- `errors.ServiceUnavailable()` factory — creates 503 errors with `SERVICE_UNAVAILABLE` code
- `errors.Timeout()` factory — creates 504 errors with `TIMEOUT` code
- `errors.RateLimited()` factory — creates 429 errors with `RATE_LIMITED` code
- 5 integration guides: pgx, SQL migrations, Zap, Zerolog, Viper (`docs/framework/guides/`)

### Changed
- Existing error factories now use exported constants instead of inline strings (no behavior change)

## [2.7.0] - 2026-03-15

### Added
- `UUIDBaseModel` in `database/models` — embed for UUID-based primary keys with auto-generation via `BeforeCreate` GORM hook
- `make:module [name]` CLI command — scaffolds a complete domain module with models, service, controller, and routes files in `modules/<name>/`

## [2.6.0] - 2026-03-15

### Added
- `Logger` interface in `core/logger` with `Debug`, `Info`, `Warn`, `Error`, `With` methods — enables pluggable logging backends (Zap, Zerolog, test spies)
- `SlogLogger` struct — default `Logger` implementation wrapping `*slog.Logger`
- `NewSlogLogger()` constructor
- `config.LoadConfig[T]()` — generic struct-based config loader using `env:`, `default:`, and `validate:` struct tags. Supports string, int, bool, float64, `time.Duration`, and `[]string` (comma-separated). Validates via `go-playground/validator`.
- `auth.GenerateTokenFromString(userID string)` — JWT generation for UUID/string-based primary keys, alongside existing `GenerateToken(uint)`

### Changed
- `logger.Setup()` now returns `Logger` interface instead of `*slog.Logger`. Global `slog.SetDefault()` still called — existing `slog.Info()` etc. continue to work unchanged.
- `GenerateToken(uint)` refactored internally to share `jwtConfig()` helper. Behavior unchanged.

## [2.5.0] - 2026-03-15

### Changed
- **BREAKING**: `AppError.Code int` renamed to `AppError.Status int`. `Code` is now a `string` field for machine-readable error codes (e.g., `"NOT_FOUND"`, `"BAD_REQUEST"`). Migration: replace `err.Code` (int) → `err.Status` (int). Caught at compile time.
- `ErrorResponse()` now includes `"code"` field in JSON output
- Error handler middleware uses `AppError.Status` for HTTP status code
- Startup banner extracted to `printBanner()` function with `APP_BANNER` env var override

### Added
- `AppError.WithCode(code string)` builder for custom machine-readable error codes
- `AppError.HTTPStatus()` deprecated helper (returns `Status` field)
- Default `Code` values on all 7 error factories: `NOT_FOUND`, `BAD_REQUEST`, `INTERNAL_ERROR`, `UNAUTHORIZED`, `FORBIDDEN`, `CONFLICT`, `UNPROCESSABLE`
- Rate limit key helpers: `KeyByIP()`, `KeyByUserID()`, `KeyByHeader()`
- `ParseRate()` for validating rate limit format strings
- `container.TryMake()` and generic `container.TryMake[T]()` — safe service resolution that returns errors instead of panicking
- `config.IsLocal()` helper — returns true for `APP_ENV=local` or `APP_ENV=development`
- `migrate:fresh` command — drops all tables and re-runs migrations (requires `--force` in production)
- `db:wipe` command — truncates all tables except `migrations` (requires `--force` in production)
- `db:seed --list` flag — lists available seeders when `SetSeederList()` is configured
- `cli.SetSeederList()` hook for registering seeder names
- `make:seeder` scaffold command — generates seeder files in `database/seeders/`

## [2.4.0] - 2026-03-15

### Fixed

- **Critical**: `serveSingle()` now calls `applyRoutesForMode()` — routes registered via `SetRoutes()` were silently ignored in single-port mode (the default for most projects)
- `serveMulti()` now copies global middleware from the container's router to per-service routers — provider-registered middleware (error handler, request ID, etc.) was previously lost when creating separate routers per service

### Added

- `Router.NoRoute()` — register custom 404 handlers without accessing the underlying Gin engine
- `Router.GlobalHandlers()` — returns the global middleware handlers registered on the router
- `health.Routes()` now accepts an optional `version` parameter — `/health` response includes `"version"` field when provided
- `parseDuration()` and `resolveServerTimeouts()` helpers for env-configurable server timeouts
- `core/cli/serve_test.go` — tests for timeout parsing and resolution
- Tests for `NoRoute`, `GlobalHandlers`, and health version features

### Changed

- `config.Load()` moved into `NewApp()` — all CLI commands that bootstrap the app (`migrate`, `db:seed`, `migrate:rollback`, `migrate:status`) now automatically load `.env` values
- Server timeouts are now configurable via `SERVER_READ_TIMEOUT`, `SERVER_WRITE_TIMEOUT`, `SERVER_IDLE_TIMEOUT`, `SERVER_SHUTDOWN_TIMEOUT` env vars (defaults unchanged: 15s, 15s, 60s, 30s)
- Removed redundant `config.Load()` calls from `work` and `schedule:run` commands
- Version constant bumped to `2.4.0`

## [2.3.0] - 2026-03-14

### Fixed

- **Critical**: `SessionMiddleware` now sets the session cookie **before** `c.Next()` so the `Set-Cookie` header is included even when handlers write the response body (e.g. `c.HTML()`). Previously the cookie was set after the body was written, causing it to be silently dropped — breaking CSRF protection, flash messages, and all session-dependent features. ([Bug #01](../../docs/bugs/01-session-cookie-bug.md))

### Added

- `session.Manager.SetCookie()` — sets the session cookie on the response without persisting data to the store, allowing cookie and store operations to be called independently
- `TestSessionMiddleware_CookieSetBeforeBody` — regression test ensuring the session cookie is present even when the handler writes an HTML response body

### Changed

- `session.Manager.Save()` now delegates to `SetCookie()` internally (no behavior change for direct callers)

## [2.2.0] - 2026-03-13

### Changed

- **Repository Rename**: Renamed GitHub repository from `RAiWorks/RapidGo` to `raiworks/rapidgo` ([#org-rename])
- **Organization Rename**: GitHub organization changed from `RAiWorks` to `raiworks`
- **Go Module Path**: Updated module path from `github.com/RAiWorks/RapidGo/v2` to `github.com/raiworks/rapidgo/v2`
- **Git Remote URL**: Updated origin from `https://github.com/RAiWorks/RapidGo.git` to `https://github.com/raiworks/rapidgo.git`
- Updated all Go import paths across 31 source files to use new lowercase module path
- Updated all documentation references (100+ docs) from old org/repo names to new lowercase names
- Updated starter project references from `RapidGo-Starter` / `RapidGo-starter` to `rapidgo-starter`
- Updated CLI `new` command and tests to reference new `rapidgo-starter` archive naming

### Files Affected

#### Go Source (31 files)
- `cmd/main.go`
- `core/app/app.go`, `core/app/app_test.go`
- `core/cli/cli_test.go`, `core/cli/hooks.go`, `core/cli/hooks_test.go`
- `core/cli/make_scaffold.go`, `core/cli/migrate.go`, `core/cli/migrate_rollback.go`
- `core/cli/migrate_status.go`, `core/cli/new.go`, `core/cli/new_test.go`
- `core/cli/root.go`, `core/cli/schedule_run.go`, `core/cli/seed.go`
- `core/cli/serve.go`, `core/cli/work.go`
- `core/errors/errors.go`
- `core/health/health.go`, `core/health/health_test.go`
- `core/logger/logger.go`
- `core/middleware/auth.go`, `core/middleware/error_handler.go`
- `core/middleware/middleware_test.go`, `core/middleware/session.go`
- `core/plugin/plugin.go`, `core/plugin/plugin_test.go`
- `core/router/router.go`
- `database/connection.go`
- `testing/testutil/testutil.go`

#### Configuration (1 file)
- `go.mod`

#### Documentation (100+ files)
- `README.md`
- All files under `docs/features/`
- All files under `docs/framework/`
- Project docs: `docs/project-context.md`, `docs/rapidgo-backlog.md`, `docs/v2-*.md`, and others
