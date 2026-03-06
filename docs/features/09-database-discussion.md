# 💬 Discussion: Database Connection

> **Feature**: `09` — Database Connection
> **Status**: 🟢 COMPLETE
> **Branch**: `docs/09-database`
> **Depends On**: #02 (Configuration System ✅), #05 (Service Container ✅)
> **Date Started**: 2026-03-06
> **Date Completed**: 2026-03-06

---

## Summary

Implement the database connection layer for the RGo framework. This feature provides a multi-driver connection factory supporting PostgreSQL, MySQL, and SQLite via GORM, configurable connection pool settings via environment variables, and a `DatabaseProvider` that registers `*gorm.DB` as a lazy singleton in the service container. The connection is established on first resolution, not at boot time.

---

## Functional Requirements

- As a **framework developer**, I want a `DBConfig` struct that reads database settings from environment variables so that connection configuration is centralized and consistent
- As a **framework developer**, I want a `Connect()` function that returns a configured `*gorm.DB` instance so that any part of the framework can establish a database connection
- As a **framework developer**, I want a `ConnectWithConfig(cfg)` function that accepts an explicit config so that tests and advanced use cases can bypass environment variables
- As a **framework developer**, I want a `DSN()` method on `DBConfig` that formats the correct connection string per driver so that DSN logic is testable in isolation
- As a **framework user**, I want multi-driver support (PostgreSQL, MySQL, SQLite) so that I can choose the database that fits my project
- As a **framework user**, I want connection pool settings (`DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME`, `DB_CONN_MAX_IDLE_TIME`) configurable via `.env` so that I can tune performance for production
- As a **framework user**, I want the database registered as a singleton in the service container so that all parts of my application share one connection pool
- As a **framework user**, I want clear error messages when the database driver is unsupported or the connection fails

## Current State / Reference

### What Exists
- **Config (#02 ✅)**: `config.Env()`, `config.EnvInt()`, `config.EnvBool()` — read env vars with fallbacks
- **Service Container (#05 ✅)**: `Singleton()`, `Make()`, `MustMake[T]()` — lazy singleton registration
- **Service Providers (#06 ✅)**: `Provider` interface (`Register` + `Boot`) and application lifecycle
- **`.env`**: Already has `DB_DRIVER`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_SSL_MODE`
- **`database/connection.go`**: Stub file — `package database` declaration only
- **`database/` subdirectories**: `migrations/`, `models/`, `querybuilder/`, `seeders/` — all empty (`.gitkeep`)
- **Provider order in `cmd/main.go`**: Config(1) → Logger(2) → Middleware(3) → Router(4)

### Blueprint Reference
The blueprint (Database Layer section, lines 326–401) shows:
1. GORM as the recommended ORM (`gorm.io/gorm`)
2. Three supported drivers: PostgreSQL, MySQL, SQLite
3. `Connect()` function — reads `os.Getenv()` directly, builds DSN per driver, opens GORM connection
4. Connection pool settings — `SetMaxOpenConns(25)`, `SetMaxIdleConns(10)`, `SetConnMaxLifetime(5 min)`, `SetConnMaxIdleTime(3 min)`
5. Note: "Tune pool values via `.env` for production (`DB_MAX_OPEN_CONNS`, etc.)"

The blueprint (Service Container section, lines 893–907) shows:
1. `DatabaseProvider` with `Singleton("db", ...)` in `Register()`
2. Lazy connection — factory runs on first `Make("db")` call
3. Panics on connection failure (boot-time concern)
4. `Boot()` placeholder for auto-migrations

### Blueprint Adaptations
| Blueprint | Our Implementation | Reason |
|---|---|---|
| `os.Getenv("DB_PASS")` | `config.Env("DB_PASSWORD", "")` | Uses our config system (#02); `.env` uses `DB_PASSWORD` |
| `os.Getenv("DB_HOST")` | `config.Env("DB_HOST", "localhost")` | Uses framework config with fallback defaults |
| Hardcoded pool values (25, 10, 5m, 3m) | `config.EnvInt()` with same defaults | Blueprint says "Tune via `.env`" — we enable that |
| `gorm.io/driver/sqlite` (CGO) | `github.com/glebarez/sqlite` (pure Go) | Works on all platforms without CGO/GCC. Same GORM API |
| Hardcoded `sslmode=disable` | `config.Env("DB_SSL_MODE", "disable")` | `.env` already has `DB_SSL_MODE` — use it |

## Proposed Approach

### DBConfig Struct + NewDBConfig()

Centralize all database configuration into a struct, populated from env vars via `config.Env()` / `config.EnvInt()`. This separates config reading from connection logic, making both independently testable.

```go
type DBConfig struct {
    Driver          string
    Host            string
    Port            string
    Name            string
    User            string
    Password        string
    SSLMode         string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}
```

### DSN() Method

Format the data source name based on the driver. Returns the DSN string. The format varies per driver:
- **postgres**: `host=H user=U password=P dbname=D port=P sslmode=S`
- **mysql**: `user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local`
- **sqlite**: just the database name/path (e.g., `:memory:` or `app.db`)

### Connect() / ConnectWithConfig()

Two entry points:
- `Connect()` — reads env via `NewDBConfig()`, calls `ConnectWithConfig()`
- `ConnectWithConfig(cfg)` — accepts explicit config (for testing, advanced use)

Both return `(*gorm.DB, error)`. The connection flow:
1. Resolve dialector (postgres/mysql/sqlite) from config
2. Open GORM connection
3. Get underlying `sql.DB`
4. Apply pool settings
5. Return `*gorm.DB`

### DatabaseProvider

Follows the existing provider pattern:
- `Register()` — registers `Singleton("db", factory)` where factory calls `database.Connect()`
- `Boot()` — no-op (migrations are a future feature)
- Panics if `Connect()` returns error — a misconfigured database is a fatal boot issue

### Provider Ordering

Database provider must come after ConfigProvider (needs `.env` loaded) and before RouterProvider (routes may need DB). Insert as provider #3, shifting Middleware to #4 and Router to #5:

```
Config(1) → Logger(2) → Database(3) → Middleware(4) → Router(5)
```

## What Is NOT In Scope

These are all deferred to future features:

| Topic | Future Feature |
|---|---|
| GORM models / BaseModel | #16 (Models & ORM) |
| Database migrations | #16 (Models & ORM) |
| Database seeders | #16 (Models & ORM) |
| Query builder helpers | #16 (Models & ORM) |
| Route model binding | #16 (Models & ORM) |
| Database transactions helper | #16 (Models & ORM) |
| Multiple database connections | Future enhancement |
| Health check / ping endpoint | Future enhancement |
| `make:model` CLI command | #10 (CLI Foundation) + #16 |

## Decisions Made

| Decision | Rationale |
|---|---|
| Pure Go SQLite (`github.com/glebarez/sqlite`) | Works on Windows/Linux/macOS without CGO. Same GORM Dialector interface. |
| Lazy connection via `Singleton` | Blueprint pattern. App can start even if DB is temporarily unavailable. Connection established on first `Make("db")`. |
| `DB_PASSWORD` not `DB_PASS` | Matches our existing `.env` convention. More explicit. |
| `ConnectWithConfig()` public API | Enables testing with explicit config. Enables advanced use (multiple DBs, custom settings). |
| Provider panics on connection failure | Consistent with blueprint. A missing database at resolution time is a programming/infrastructure error, not a runtime condition. |
| Pool defaults match blueprint | 25 max open, 10 max idle, 5m lifetime, 3m idle time. Configurable via env vars. |
