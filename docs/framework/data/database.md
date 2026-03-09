---
title: "Database"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Database

## Abstract

This document covers database connectivity using GORM — supported
drivers, connection factory, connection pool tuning, and multi-driver
configuration via `.env`.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Supported Drivers](#2-supported-drivers)
3. [Configuration](#3-configuration)
4. [Connection Factory](#4-connection-factory)
5. [Connection Pool](#5-connection-pool)
6. [Provider Registration](#6-provider-registration)
7. [Security Considerations](#7-security-considerations)
8. [References](#8-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **ORM** — Object-Relational Mapping; maps Go structs to database
  tables.
- **DSN** — Data Source Name; the connection string for a database.

## 2. Supported Drivers

| Driver | Library | Config Value |
|--------|---------|-------------|
| PostgreSQL | `gorm.io/driver/postgres` | `postgres` |
| MySQL | `gorm.io/driver/mysql` | `mysql` |
| SQLite | `gorm.io/driver/sqlite` | `sqlite` |

The driver is selected via the `DB_DRIVER` environment variable.

## 3. Configuration

`.env` database variables:

```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=secret
DB_NAME=myapp
```

## 4. Connection Factory

The `Connect()` function resolves the correct driver from `DB_DRIVER`
and returns a configured `*gorm.DB`:

```go
package database

import (
    "fmt"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
    driver := os.Getenv("DB_DRIVER")

    var dialector gorm.Dialector

    switch driver {
    case "postgres":
        dsn := fmt.Sprintf(
            "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
            os.Getenv("DB_HOST"), os.Getenv("DB_USER"),
            os.Getenv("DB_PASS"), os.Getenv("DB_NAME"),
            os.Getenv("DB_PORT"),
        )
        dialector = postgres.Open(dsn)
    case "mysql":
        dsn := fmt.Sprintf(
            "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
            os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
            os.Getenv("DB_NAME"),
        )
        dialector = mysql.Open(dsn)
    case "sqlite":
        dialector = sqlite.Open(os.Getenv("DB_NAME"))
    default:
        return nil, fmt.Errorf("unsupported DB_DRIVER: %s", driver)
    }

    db, err := gorm.Open(dialector, &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("database connection failed: %w", err)
    }

    // Connection pool settings
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)
    sqlDB.SetConnMaxIdleTime(3 * time.Minute)

    return db, nil
}
```

## 5. Connection Pool

Tune pool settings via `.env` for production:

| Setting | Default | Description |
|---------|---------|-------------|
| `MaxOpenConns` | 25 | Maximum simultaneous connections |
| `MaxIdleConns` | 10 | Idle connections kept alive |
| `ConnMaxLifetime` | 5 min | Recycle connections after this duration |
| `ConnMaxIdleTime` | 3 min | Close idle connections after this duration |

### Recommendations

- **Development:** Defaults are sufficient.
- **Production:** Tune based on your database server capacity and
  expected concurrency. Monitor connection usage under load.

## 6. Provider Registration

The `DatabaseProvider` registers the connection as a singleton in the
service container:

```go
type DatabaseProvider struct{}

func (p *DatabaseProvider) Register(c *container.Container) {
    c.Singleton("db", func(c *container.Container) interface{} {
        db, err := database.Connect()
        if err != nil {
            panic("database connection failed: " + err.Error())
        }
        return db
    })
}

func (p *DatabaseProvider) Boot(c *container.Container) {
    // Run auto-migrations if enabled
}
```

Resolve in controllers or services:

```go
db := container.MustMake[*gorm.DB](app.Container, "db")
```

## 7. Security Considerations

- Database credentials **MUST** be stored in `.env` and never
  committed to version control.
- PostgreSQL connections **SHOULD** use `sslmode=require` in
  production.
- Connection pool limits prevent exhaustion attacks from overloading
  the database.
- Use parameterized queries (GORM's default) to prevent SQL
  injection.

## 8. References

- [Models](models.md)
- [Migrations](migrations.md)
- [Transactions](transactions.md)
- [Configuration](../core/configuration.md)
- [Service Providers](../core/service-providers.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
