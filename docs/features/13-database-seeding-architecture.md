# 🏗️ Architecture: Database Seeding

> **Feature**: `13` — Database Seeding
> **Discussion**: [`13-database-seeding-discussion.md`](13-database-seeding-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

The Database Seeding feature provides a registry-based seeder system. Each seeder implements the `Seeder` interface, self-registers, and is invoked via the `rgo db:seed` CLI command. A sample `UserSeeder` demonstrates the pattern with idempotent user creation.

## File Structure

```
database/seeders/
├── seeder.go          # Seeder interface, registry, RunAll/RunByName
└── user_seeder.go     # UserSeeder — creates admin + regular user

core/cli/
└── seed.go            # db:seed command with --seeder flag

database/seeders/
└── seeders_test.go    # Tests for seeder system
```

### Files Created (3)
| File | Package | Lines (est.) |
|---|---|---|
| `database/seeders/seeder.go` | `seeders` | ~60 |
| `database/seeders/user_seeder.go` | `seeders` | ~35 |
| `core/cli/seed.go` | `cli` | ~35 |

### Files Modified (1)
| File | Change |
|---|---|
| `core/cli/root.go` | Register `dbSeedCmd` in `init()` |

---

## Component Design

### Seeder Interface & Registry (`database/seeders/seeder.go`)

**Responsibility**: Define the seeder contract and provide registration/execution.
**Package**: `seeders`

```go
package seeders

import (
	"fmt"

	"gorm.io/gorm"
)

// Seeder defines the interface for database seeders.
type Seeder interface {
	// Name returns the seeder's unique name (used with --seeder flag).
	Name() string
	// Seed populates the database with data.
	Seed(db *gorm.DB) error
}

// registry holds all registered seeders.
var registry []Seeder

// Register adds a seeder to the global registry.
func Register(s Seeder) {
	registry = append(registry, s)
}

// ResetRegistry clears all registered seeders. Used in tests only.
func ResetRegistry() {
	registry = nil
}

// RunAll executes all registered seeders in registration order.
func RunAll(db *gorm.DB) error {
	for _, s := range registry {
		if err := s.Seed(db); err != nil {
			return fmt.Errorf("seeder %s failed: %w", s.Name(), err)
		}
	}
	return nil
}

// RunByName executes a single seeder by name.
func RunByName(db *gorm.DB, name string) error {
	for _, s := range registry {
		if s.Name() == name {
			return s.Seed(db)
		}
	}
	return fmt.Errorf("seeder %q not found", name)
}

// Names returns the names of all registered seeders.
func Names() []string {
	names := make([]string, len(registry))
	for i, s := range registry {
		names[i] = s.Name()
	}
	return names
}
```

**Design notes**:
- `Seeder` interface — `Name()` + `Seed(db)`, minimal contract
- `registry` is package-level — consistent with migrations pattern
- `RunAll` stops on first error — fail-fast behavior
- `RunByName` enables `--seeder` flag on the CLI command
- `Names()` useful for status/listing
- `ResetRegistry()` for test isolation (same pattern as migrations)

### UserSeeder (`database/seeders/user_seeder.go`)

**Responsibility**: Create default admin and regular user accounts.
**Package**: `seeders`

```go
package seeders

import (
	"github.com/RAiWorks/RGo/database/models"
	"gorm.io/gorm"
)

func init() {
	Register(&UserSeeder{})
}

// UserSeeder creates default user accounts.
type UserSeeder struct{}

func (s *UserSeeder) Name() string { return "UserSeeder" }

func (s *UserSeeder) Seed(db *gorm.DB) error {
	users := []models.User{
		{Name: "Admin", Email: "admin@example.com", Password: "password123", Role: "admin"},
		{Name: "User", Email: "user@example.com", Password: "password123", Role: "user"},
	}
	for _, u := range users {
		// TODO: Hash password when Feature #19/#22 ships
		if err := db.FirstOrCreate(&u, models.User{Email: u.Email}).Error; err != nil {
			return err
		}
	}
	return nil
}
```

**Design notes**:
- `init()` registers the seeder automatically when the package is imported
- `FirstOrCreate` with Email match — idempotent, per blueprint
- Plaintext passwords — documented TODO, will be updated with crypto feature
- Follows blueprint example closely (admin + regular user)

### Seed Command (`core/cli/seed.go`)

**Responsibility**: CLI command to run database seeders.
**Package**: `cli`

```go
package cli

import (
	"fmt"

	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/database/seeders"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var dbSeedCmd = &cobra.Command{
	Use:   "db:seed",
	Short: "Seed the database with records",
	RunE: func(cmd *cobra.Command, args []string) error {
		application := NewApp()
		db := container.MustMake[*gorm.DB](application.Container, "db")

		name, _ := cmd.Flags().GetString("seeder")
		if name != "" {
			if err := seeders.RunByName(db, name); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Seeder %s complete.\n", name)
			return nil
		}

		if err := seeders.RunAll(db); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Database seeding complete.")
		return nil
	},
}

func init() {
	dbSeedCmd.Flags().String("seeder", "", "Run a specific seeder by name")
}
```

**Design notes**:
- `--seeder` flag — optional, runs single seeder when provided
- Uses `NewApp()` to bootstrap (same pattern as migrate commands)
- `application.Container` — field access, not method call
- `init()` registers the `--seeder` flag on the command

### Root Command Registration (`core/cli/root.go`)

Add to `init()`:
```go
rootCmd.AddCommand(dbSeedCmd)
```

---

## Data Flow

### `rgo db:seed`
```
CLI → NewApp() → get *gorm.DB from container
    → seeders.RunAll(db)
        → for each registered seeder: call Seed(db)
        → UserSeeder.Seed(db): FirstOrCreate admin, FirstOrCreate user
    → print "Database seeding complete."
```

### `rgo db:seed --seeder UserSeeder`
```
CLI → NewApp() → get *gorm.DB
    → seeders.RunByName(db, "UserSeeder")
        → find seeder with Name() == "UserSeeder"
        → call Seed(db)
    → print "Seeder UserSeeder complete."
```

---

## Constraints & Invariants

1. Seeders are **idempotent** — `FirstOrCreate` prevents duplicates
2. `RunAll` executes in **registration order** (deterministic via `init()`)
3. `RunAll` **stops on first error** — fail-fast
4. `--seeder` flag must match a registered name exactly (case-sensitive)
5. Passwords are **plaintext** until Feature #19/#22 ships
6. `UserSeeder` self-registers via `init()` — imported by the `seed.go` CLI command
