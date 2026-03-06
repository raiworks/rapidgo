# 🏗️ Architecture: Database Migrations

> **Feature**: `12` — Database Migrations
> **Discussion**: [`12-database-migrations-discussion.md`](12-database-migrations-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

The Database Migrations feature provides a two-tier migration system: (1) GORM `AutoMigrate` applies model structs to the database, and (2) file-based migrations handle changes that AutoMigrate cannot (column renames, drops, data migrations). A central model registry keeps model-to-table mapping in one place. Four CLI commands drive the system: `migrate`, `make:migration`, `migrate:rollback`, and `migrate:status`.

## File Structure

```
database/models/
└── registry.go          # All() function — returns models for AutoMigrate

database/migrations/
├── migrator.go          # Migration engine: Run, Rollback, Status, Register
└── migrations_test.go   # Tests for the migration engine

core/cli/
├── migrate.go           # migrate command — AutoMigrate + pending migrations
├── migrate_rollback.go  # migrate:rollback command
├── migrate_status.go    # migrate:status command
└── make_migration.go    # make:migration command — generates migration file
```

### Files Created (6)
| File | Package | Lines (est.) |
|---|---|---|
| `database/models/registry.go` | `models` | ~10 |
| `database/migrations/migrator.go` | `migrations` | ~120 |
| `core/cli/migrate.go` | `cli` | ~40 |
| `core/cli/migrate_rollback.go` | `cli` | ~30 |
| `core/cli/migrate_status.go` | `cli` | ~40 |
| `core/cli/make_migration.go` | `cli` | ~60 |

### Files Modified (1)
| File | Change |
|---|---|
| `core/cli/root.go` | Register 4 new subcommands in `init()` |

---

## Component Design

### Model Registry (`database/models/registry.go`)

**Responsibility**: Provide a single list of all GORM models for AutoMigrate.
**Package**: `models`

```go
package models

// All returns all model structs for GORM AutoMigrate.
// Add new models here as they are created.
func All() []interface{} {
	return []interface{}{
		&User{},
		&Post{},
	}
}
```

**Design notes**:
- Returns `[]interface{}` — GORM's `AutoMigrate` accepts `...interface{}`
- Single source of truth — no scattered AutoMigrate calls

### Migration Engine (`database/migrations/migrator.go`)

**Responsibility**: Track, run, and rollback file-based migrations.
**Package**: `migrations`

```go
package migrations

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

// SchemaMigration tracks applied migrations in the database.
type SchemaMigration struct {
	ID        uint   `gorm:"primaryKey"`
	Version   string `gorm:"size:255;uniqueIndex;not null"`
	Batch     int    `gorm:"not null"`
	AppliedAt time.Time
}

// MigrationFunc is a function that performs a migration step.
type MigrationFunc func(db *gorm.DB) error

// Migration represents a single migration with up and down functions.
type Migration struct {
	Version string
	Up      MigrationFunc
	Down    MigrationFunc
}

// registry holds all registered migrations (populated via Register).
var registry []Migration

// Register adds a migration to the global registry.
// Called from migration files (typically in init()).
func Register(m Migration) {
	registry = append(registry, m)
}

// Migrator manages database migrations.
type Migrator struct {
	DB *gorm.DB
}

// NewMigrator creates a Migrator and ensures the schema_migrations table exists.
func NewMigrator(db *gorm.DB) (*Migrator, error) {
	if err := db.AutoMigrate(&SchemaMigration{}); err != nil {
		return nil, fmt.Errorf("failed to create schema_migrations table: %w", err)
	}
	return &Migrator{DB: db}, nil
}

// Run applies all pending migrations in version order.
// Returns the number of migrations applied.
func (m *Migrator) Run() (int, error)

// Rollback undoes the last batch of applied migrations.
// Returns the number of migrations rolled back.
func (m *Migrator) Rollback() (int, error)

// Status returns the status of all registered migrations.
func (m *Migrator) Status() ([]MigrationStatus, error)
```

**MigrationStatus** for the status command:
```go
// MigrationStatus represents the status of a single migration.
type MigrationStatus struct {
	Version string
	Applied bool
	Batch   int
}
```

**Design notes**:
- `SchemaMigration` table — auto-created, tracks version string + batch number
- `registry` is package-level — migrations self-register via `Register()` in `init()`
- `Run()` — finds pending (registered but not in schema_migrations), applies in sorted order, same batch number
- `Rollback()` — finds highest batch, runs Down() in reverse version order, deletes from schema_migrations
- `Status()` — merges registry with schema_migrations, returns sorted list
- `Migrator` takes `*gorm.DB` — testable with SQLite in-memory

### Migrate Command (`core/cli/migrate.go`)

**Responsibility**: Run AutoMigrate for all models, then apply pending file-based migrations.
**Package**: `cli`

```go
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := NewApp()
		db := container.MustMake[*gorm.DB](app.Container, "db")

		// Step 1: AutoMigrate all models
		if err := db.AutoMigrate(models.All()...); err != nil {
			return fmt.Errorf("auto-migrate failed: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "AutoMigrate complete.")

		// Step 2: Run pending file-based migrations
		migrator, err := migrations.NewMigrator(db)
		if err != nil {
			return err
		}
		n, err := migrator.Run()
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
		if n == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Nothing to migrate.")
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Applied %d migration(s).\n", n)
		}
		return nil
	},
}
```

### Rollback Command (`core/cli/migrate_rollback.go`)

```go
var migrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last batch of migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := NewApp()
		db := container.MustMake[*gorm.DB](app.Container, "db")

		migrator, err := migrations.NewMigrator(db)
		if err != nil {
			return err
		}
		n, err := migrator.Rollback()
		if err != nil {
			return fmt.Errorf("rollback failed: %w", err)
		}
		if n == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Nothing to rollback.")
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Rolled back %d migration(s).\n", n)
		}
		return nil
	},
}
```

### Status Command (`core/cli/migrate_status.go`)

```go
var migrateStatusCmd = &cobra.Command{
	Use:   "migrate:status",
	Short: "Show the status of all migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := NewApp()
		db := container.MustMake[*gorm.DB](app.Container, "db")

		migrator, err := migrations.NewMigrator(db)
		if err != nil {
			return err
		}
		statuses, err := migrator.Status()
		if err != nil {
			return err
		}
		if len(statuses) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No migrations registered.")
			return nil
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Migration                | Status  | Batch")
		fmt.Fprintln(cmd.OutOrStdout(), "-------------------------+---------+------")
		for _, s := range statuses {
			status := "Pending"
			batch := ""
			if s.Applied {
				status = "Applied"
				batch = fmt.Sprintf("%d", s.Batch)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%-25s| %-8s| %s\n", s.Version, status, batch)
		}
		return nil
	},
}
```

### Make Migration Command (`core/cli/make_migration.go`)

**Responsibility**: Generate a timestamped migration file with Up/Down stubs.
**Package**: `cli`

```go
var makeMigrationCmd = &cobra.Command{
	Use:   "make:migration [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		timestamp := time.Now().Format("20060102150405")
		version := timestamp + "_" + toSnakeCase(name)
		filename := version + ".go"
		path := filepath.Join("database", "migrations", filename)

		if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()

		t := template.Must(template.New("migration").Parse(migrationTpl))
		if err := t.Execute(f, map[string]string{"Version": version}); err != nil {
			return fmt.Errorf("failed to write template: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Migration created: %s\n", path)
		return nil
	},
}
```

Migration template:
```go
var migrationTpl = `package migrations

import "gorm.io/gorm"

func init() {
	Register(Migration{
		Version: "{{.Version}}",
		Up: func(db *gorm.DB) error {
			// TODO: implement migration
			return nil
		},
		Down: func(db *gorm.DB) error {
			// TODO: implement rollback
			return nil
		},
	})
}
`
```

**`toSnakeCase` helper** (internal to the file):
```go
func toSnakeCase(s string) string {
	// Convert PascalCase/camelCase to snake_case
	// "CreateUsersTable" → "create_users_table"
}
```

### Root Command Registration (`core/cli/root.go`)

Add to `init()`:
```go
rootCmd.AddCommand(migrateCmd)
rootCmd.AddCommand(migrateRollbackCmd)
rootCmd.AddCommand(migrateStatusCmd)
rootCmd.AddCommand(makeMigrationCmd)
```

---

## Data Flow

### `rgo migrate`
```
CLI → NewApp() → get *gorm.DB from container
    → db.AutoMigrate(models.All()...)       ← step 1: sync model structs
    → NewMigrator(db)                       ← ensures schema_migrations table
    → migrator.Run()                        ← step 2: apply pending migrations
        → query schema_migrations for applied versions
        → filter registry for pending
        → sort by version
        → for each: call Up(db), insert SchemaMigration row
    → print results
```

### `rgo migrate:rollback`
```
CLI → NewApp() → get *gorm.DB
    → NewMigrator(db) → migrator.Rollback()
        → find max batch in schema_migrations
        → get all migrations in that batch
        → sort by version descending
        → for each: call Down(db), delete SchemaMigration row
    → print results
```

### `rgo make:migration CreateSessionsTable`
```
CLI → format timestamp + snake_case name
    → create database/migrations/20260306120000_create_sessions_table.go
    → write template with init() → Register(Migration{...})
```

---

## Constraints & Invariants

1. **schema_migrations** table is auto-created by `NewMigrator` — no manual setup
2. Migrations run in **sorted version order** (lexicographic = chronological)
3. Rollback only undoes the **last batch** — not all migrations
4. Each `Run()` call assigns a new batch number (max(batch) + 1)
5. `AutoMigrate` runs before file-based migrations — always
6. Generated migration files must compile — they import the `migrations` package
7. `make:migration` uses `os.OpenFile` with `O_CREATE|O_EXCL` — fails if file already exists (safety net, though timestamp guarantees uniqueness)
