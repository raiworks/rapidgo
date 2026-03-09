---
title: "Seeders"
version: "1.0.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Seeders

## Abstract

This document covers database seeding — populating the database with
initial or test data using the interface-based seeder registry and CLI
commands.

RapidGo uses a **registry-based seeder system** where each seeder
implements the `Seeder` interface. Seeders are registered globally and
executed via CLI. You can run all seeders at once or target a specific
one by name.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Architecture](#2-architecture)
3. [Seeder Interface](#3-seeder-interface)
4. [Writing Seeders](#4-writing-seeders)
5. [Registering Seeders](#5-registering-seeders)
6. [CLI Commands](#6-cli-commands)
7. [Wiring in main.go](#7-wiring-in-maingo)
8. [Best Practices](#8-best-practices)
9. [Security Considerations](#9-security-considerations)
10. [References](#10-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Seeder** — A struct implementing the `Seeder` interface that
  inserts predefined records into the database.
- **Registry** — A global list of all registered seeders; seeders
  execute in registration order.
- **`FirstOrCreate`** — GORM method that inserts only if a matching
  record doesn't exist, making seeders idempotent.

## 2. Architecture

RapidGo's seeder system consists of three layers:

```
┌─────────────────────────────────────────────────┐
│           CLI: rapidgo db:seed                  │
│           (--seeder flag for single run)         │
├─────────────────────────────────────────────────┤
│         Seeder Registry (database/seeders)       │
│  Register() → RunAll() / RunByName() → Names()  │
├─────────────────────────────────────────────────┤
│         Seeder Implementations (app layer)       │
│  UserSeeder, PostSeeder, etc.                    │
└─────────────────────────────────────────────────┘
```

- **Framework provides**: `Seeder` interface, global registry,
  `RunAll()`, `RunByName()`, `Names()`, and the `db:seed` CLI command
- **Application provides**: Concrete seeder structs that implement
  the `Seeder` interface

## 3. Seeder Interface

The seeder engine lives in `database/seeders/seeder.go`:

```go
// Seeder defines the interface for database seeders.
type Seeder interface {
    // Name returns the seeder's unique name (used with --seeder flag).
    Name() string
    // Seed populates the database with data.
    Seed(db *gorm.DB) error
}
```

### Registry Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `Register` | `Register(s Seeder)` | Adds a seeder to the global registry |
| `RunAll` | `RunAll(db *gorm.DB) error` | Runs all registered seeders in order |
| `RunByName` | `RunByName(db *gorm.DB, name string) error` | Runs a single seeder by name |
| `Names` | `Names() []string` | Returns names of all registered seeders |
| `ResetRegistry` | `ResetRegistry()` | Clears all seeders (test-only) |

### Execution Behavior

- `RunAll` executes seeders **in registration order** and **stops on
  first error**, returning a wrapped error with the seeder name.
- `RunByName` returns `"seeder %q not found"` if no matching seeder
  is registered.

## 4. Writing Seeders

Each seeder is a struct that implements the `Seeder` interface. Place
seeder files in your application's `database/seeders/` directory:

```go
package seeders

import (
    "github.com/RAiWorks/RapidGo/v2/core/crypto"
    "gorm.io/gorm"
)

// User is a local model reference for seeding.
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string `gorm:"size:255"`
    Email    string `gorm:"size:255;uniqueIndex"`
    Password string `gorm:"size:255"`
    Role     string `gorm:"size:50"`
}

// UserSeeder seeds initial user accounts.
type UserSeeder struct{}

func (s *UserSeeder) Name() string { return "users" }

func (s *UserSeeder) Seed(db *gorm.DB) error {
    users := []struct {
        Name     string
        Email    string
        Password string
        Role     string
    }{
        {"Admin", "admin@example.com", "Admin@Str0ng!Pass", "admin"},
        {"User", "user@example.com", "User@Str0ng!Pass", "user"},
    }

    for _, u := range users {
        hashed, err := crypto.HashPassword(u.Password)
        if err != nil {
            return err
        }
        record := User{
            Name:     u.Name,
            Email:    u.Email,
            Password: hashed,
            Role:     u.Role,
        }
        if err := db.Where(User{Email: u.Email}).FirstOrCreate(&record).Error; err != nil {
            return err
        }
    }
    return nil
}
```

Key practices:
- **Hash passwords** before insertion — never store plain text.
- **Use `FirstOrCreate`** to make seeders idempotent (safe to run
  multiple times).
- **Return errors** instead of logging and continuing — the registry
  handles error reporting.
- **Use realistic data** for development; avoid production secrets
  in seed files.

## 5. Registering Seeders

Register all seeders in an `init()` function or a dedicated
registration function:

```go
package seeders

import (
    dbseeders "github.com/RAiWorks/RapidGo/v2/database/seeders"
)

func init() {
    dbseeders.Register(&UserSeeder{})
    dbseeders.Register(&PostSeeder{})
    // Add more seeders here — execution order matches registration order
}
```

## 6. CLI Commands

### Run All Seeders

```bash
rapidgo db:seed
```

Executes all registered seeders in registration order.

### Run a Specific Seeder

```bash
rapidgo db:seed --seeder users
```

Runs only the seeder whose `Name()` returns `"users"`.

### How It Works

The `db:seed` command:
1. Bootstraps the application (loads config, connects to database)
2. Calls the seeder function set via `cli.SetSeeder()`
3. If `--seeder` flag is provided, calls `RunByName(db, name)`
4. Otherwise, calls `RunAll(db)`
5. Prints a success message on completion

## 7. Wiring in main.go

Connect the seeder system to the CLI using the `SetSeeder` hook:

```go
cli.SetSeeder(func(db *gorm.DB, name string) error {
    if name != "" {
        return seeders.RunByName(db, name)
    }
    return seeders.RunAll(db)
})
```

This gives the CLI access to your application's registered seeders
without the framework depending on your application code directly.

## 8. Best Practices

| Practice | Why |
|----------|-----|
| One seeder per entity | Keeps seeders focused and independently runnable |
| Use `FirstOrCreate` | Makes seeders idempotent — safe to run repeatedly |
| Hash credentials | Never store plain-text passwords in the database |
| Use strong dev passwords | Even development data should follow security practices |
| Registration order matters | Dependencies should be registered first (e.g., users before posts) |
| Return errors, don't swallow | Let the registry report failures with the seeder name |

## 9. Security Considerations

- Seed data for admin accounts **MUST** use strong passwords, even
  in development.
- Seeder files **MUST NOT** be deployed to production unless
  specifically intended for initial data setup.
- Passwords in seed data **MUST** be hashed — never stored as
  plain text.
- Seed files **MUST NOT** contain real credentials, API keys, or
  production secrets.

## 10. References

- [Database](database.md)
- [Models](models.md)
- [Migrations](migrations.md)
- [CLI Overview](../cli/cli-overview.md)
- [Crypto](../security/crypto.md) — for `HashPassword()`

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
| 1.0.0 | 2026-03-10 | RAiWorks | Rewritten to match interface-based registry implementation. Added architecture diagram, registration, CLI flags, wiring hook, best practices. |
