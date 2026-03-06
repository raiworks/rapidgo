# 🏗️ Architecture: Database Connection

> **Feature**: `09` — Database Connection
> **Discussion**: [`09-database-discussion.md`](09-database-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

The Database Connection feature provides a GORM-based connection factory with multi-driver support (PostgreSQL, MySQL, SQLite), configurable connection pool settings via environment variables, and a `DatabaseProvider` that registers `*gorm.DB` as a lazy singleton in the service container. The connection is established on first resolution, not at boot time.

## File Structure

```
database/
└── connection.go        # DBConfig, NewDBConfig, DSN, Connect, ConnectWithConfig

database/
└── database_test.go     # Tests for config, DSN, connection, pool settings

app/providers/
└── database_provider.go # DatabaseProvider — registers *gorm.DB singleton
```

### Files Created (1)
| File | Package | Lines (est.) |
|---|---|---|
| `app/providers/database_provider.go` | `providers` | ~25 |

### Files Modified (2)
| File | Change |
|---|---|
| `database/connection.go` | Expand from stub to full implementation: `DBConfig`, `NewDBConfig`, `DSN`, `Connect`, `ConnectWithConfig`, `newDialector` |
| `cmd/main.go` | Insert `DatabaseProvider` as provider #3 (Middleware → #4, Router → #5) |

---

## Component Design

### Database Configuration (`database/connection.go` — DBConfig)

**Responsibility**: Hold all database connection parameters in a single struct, populated from environment variables.
**Package**: `database`

```go
package database

import (
	"fmt"
	"time"

	"github.com/RAiWorks/RGo/core/config"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConfig holds database connection configuration.
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

// NewDBConfig reads database configuration from environment variables.
func NewDBConfig() DBConfig {
	return DBConfig{
		Driver:          config.Env("DB_DRIVER", ""),
		Host:            config.Env("DB_HOST", "localhost"),
		Port:            config.Env("DB_PORT", "5432"),
		Name:            config.Env("DB_NAME", "rgo_dev"),
		User:            config.Env("DB_USER", ""),
		Password:        config.Env("DB_PASSWORD", ""),
		SSLMode:         config.Env("DB_SSL_MODE", "disable"),
		MaxOpenConns:    config.EnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    config.EnvInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime: time.Duration(config.EnvInt("DB_CONN_MAX_LIFETIME", 5)) * time.Minute,
		ConnMaxIdleTime: time.Duration(config.EnvInt("DB_CONN_MAX_IDLE_TIME", 3)) * time.Minute,
	}
}
```

**Design notes**:
- Uses `config.Env()` / `config.EnvInt()` from Feature #02 — consistent with framework patterns
- `DB_DRIVER` defaults to empty string — forces explicit configuration. An unset driver produces a clear error from `Connect()`
- `DB_PASSWORD` (not `DB_PASS`) — matches our `.env` convention
- Pool settings default to blueprint values (25, 10, 5m, 3m) — configurable via env vars for production tuning
- `ConnMaxLifetime` / `ConnMaxIdleTime` read as integer minutes then convert to `time.Duration`

### DSN Builder (`database/connection.go` — DSN method)

**Responsibility**: Format the data source name string for the configured driver.

```go
// DSN returns the data source name for the configured driver.
func (cfg DBConfig) DSN() string {
	switch cfg.Driver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode,
		)
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
		)
	case "sqlite":
		return cfg.Name
	default:
		return ""
	}
}
```

**Design notes**:
- PostgreSQL DSN includes `sslmode` from config (blueprint hardcodes `disable`)
- MySQL DSN includes `charset=utf8mb4&parseTime=True&loc=Local` — standard for Go MySQL
- SQLite DSN is just the file path or `:memory:` — per GORM convention
- Unsupported drivers return empty string — the error is raised in `newDialector()`, not here

### Connection Factory (`database/connection.go` — Connect)

**Responsibility**: Establish a GORM database connection with pool configuration.

```go
// Connect establishes a database connection using environment configuration.
func Connect() (*gorm.DB, error) {
	return ConnectWithConfig(NewDBConfig())
}

// ConnectWithConfig establishes a database connection using the provided configuration.
func ConnectWithConfig(cfg DBConfig) (*gorm.DB, error) {
	dialector, err := newDialector(cfg)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}

// newDialector creates the appropriate GORM dialector for the configured driver.
func newDialector(cfg DBConfig) (gorm.Dialector, error) {
	switch cfg.Driver {
	case "postgres":
		return postgres.Open(cfg.DSN()), nil
	case "mysql":
		return mysql.Open(cfg.DSN()), nil
	case "sqlite":
		return sqlite.Open(cfg.DSN()), nil
	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER: %s", cfg.Driver)
	}
}
```

**Design notes**:
- `Connect()` is the primary public API — reads env vars, creates config, connects
- `ConnectWithConfig()` is the testable entry point — accepts explicit config, no env dependency
- `newDialector()` is internal — validates driver and creates the appropriate GORM dialector
- Error wrapping with `%w` enables callers to inspect underlying errors
- Pool settings are applied to the underlying `sql.DB` after connection
- Uses `github.com/glebarez/sqlite` (pure Go) — works on all platforms without CGO

### DatabaseProvider (`app/providers/database_provider.go`)

**Responsibility**: Register `*gorm.DB` as a lazy singleton in the service container.
**Package**: `providers`

```go
package providers

import (
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/database"
)

// DatabaseProvider registers the database connection in the service container.
type DatabaseProvider struct{}

// Register binds a *gorm.DB singleton. The connection is established lazily
// on first container.Make("db") call, not at registration time.
func (p *DatabaseProvider) Register(c *container.Container) {
	c.Singleton("db", func(c *container.Container) interface{} {
		db, err := database.Connect()
		if err != nil {
			panic("database connection failed: " + err.Error())
		}
		return db
	})
}

// Boot is a no-op. Migrations and seeding are future features.
func (p *DatabaseProvider) Boot(c *container.Container) {}
```

**Design notes**:
- `Singleton` factory is lazy — runs on first `Make("db")`, not at `Register()` time
- Panics on connection failure — consistent with blueprint pattern (misconfigured DB is fatal)
- `Boot()` is a no-op — auto-migrations will be added in Feature #16 (Models & ORM)
- The provider does not import GORM directly — it delegates to `database.Connect()`

---

## Data Flow

### Connection Lifecycle

```
┌──────────────────────────────────────────────────────────┐
│  Application Boot                                        │
│                                                          │
│  1. ConfigProvider.Register()  → loads .env (DB_* vars)  │
│  2. DatabaseProvider.Register()→ registers Singleton     │
│     (factory stored, NOT executed)                       │
│                                                          │
│  ... other providers register ...                        │
│                                                          │
│  3. application.Boot()                                   │
│     - ConfigProvider.Boot()    → no-op                   │
│     - LoggerProvider.Boot()    → sets up slog            │
│     - DatabaseProvider.Boot()  → no-op                   │
│     - MiddlewareProvider.Boot()→ registers aliases       │
│     - RouterProvider.Boot()    → loads routes            │
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│  First Request Needing DB                                │
│                                                          │
│  handler calls: container.MustMake[*gorm.DB](c, "db")   │
│                                                          │
│  1. Singleton factory executes (first time only)         │
│  2. database.Connect() →                                 │
│     a. NewDBConfig()    → reads DB_* env vars            │
│     b. newDialector()   → creates postgres/mysql/sqlite  │
│     c. gorm.Open()      → establishes connection         │
│     d. Pool config      → MaxOpen, MaxIdle, Lifetime     │
│  3. *gorm.DB cached in container                         │
│  4. Subsequent Make("db") calls return cached instance   │
└──────────────────────────────────────────────────────────┘
```

### Provider Ordering (after this feature)

```go
application.Register(&providers.ConfigProvider{})      // 1. Config — loads .env
application.Register(&providers.LoggerProvider{})       // 2. Logger — uses config in Boot
application.Register(&providers.DatabaseProvider{})     // 3. Database — registers lazy singleton
application.Register(&providers.MiddlewareProvider{})   // 4. Middleware — registers aliases
application.Register(&providers.RouterProvider{})       // 5. Router — creates engine, defines routes
```

---

## Dependencies (New)

| Package | Version | Purpose |
|---|---|---|
| `gorm.io/gorm` | latest | Core ORM — `*gorm.DB`, `gorm.Config`, `gorm.Dialector` |
| `gorm.io/driver/postgres` | latest | PostgreSQL driver |
| `gorm.io/driver/mysql` | latest | MySQL driver |
| `github.com/glebarez/sqlite` | latest | SQLite driver (pure Go, no CGO) |

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_DRIVER` | (empty) | Database driver: `postgres`, `mysql`, or `sqlite` |
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_NAME` | `rgo_dev` | Database name (or SQLite file path) |
| `DB_USER` | (empty) | Database username |
| `DB_PASSWORD` | (empty) | Database password |
| `DB_SSL_MODE` | `disable` | PostgreSQL SSL mode |
| `DB_MAX_OPEN_CONNS` | `25` | Maximum open connections |
| `DB_MAX_IDLE_CONNS` | `10` | Maximum idle connections |
| `DB_CONN_MAX_LIFETIME` | `5` | Connection max lifetime (minutes) |
| `DB_CONN_MAX_IDLE_TIME` | `3` | Connection max idle time (minutes) |
