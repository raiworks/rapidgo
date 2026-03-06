# 💬 Discussion: Database Migrations

> **Feature**: `12` — Database Migrations
> **Status**: 🟢 COMPLETE
> **Date**: 2026-03-06

---

## What Are We Building?

A CLI-driven database migration system that uses GORM's `AutoMigrate` to apply model schemas to the database. The `migrate` command connects to the configured database and runs `AutoMigrate` against all registered models.

## Blueprint References

The blueprint specifies:

1. **CLI command**: `framework migrate` (line 495)
2. **Directory**: `database/migrations/` in the project structure (line 71)
3. **DatabaseProvider.Boot()**: "Run auto-migrations if enabled" (line 907)
4. **Auto-migration approach**: The blueprint uses `db.AutoMigrate()` throughout (lines 907, 1753)

## Scope for Feature #12

### In Scope
- `migrate` CLI subcommand — registered under root, runs GORM AutoMigrate
- Centralized model registry — single function that returns all models to migrate
- `make:migration` CLI subcommand — generates a timestamped migration file with Up/Down stubs
- Migration runner — scans `database/migrations/`, runs pending Up() functions in order
- Migration tracking table — `schema_migrations` table records which migrations have been applied
- `migrate:rollback` CLI subcommand — rolls back the last batch of migrations
- `migrate:status` CLI subcommand — shows which migrations have been applied

### Out of Scope (deferred)
- `make:model` (Feature #15+ scaffolding commands)
- `db:seed` (Feature #13)
- Fresh/reset/wipe commands (future feature)
- Down migrations for AutoMigrate (GORM AutoMigrate only adds, never drops)

## Key Design Decisions

### 1. GORM AutoMigrate as Foundation
The blueprint consistently uses `AutoMigrate`. This is idiomatic for Go/GORM — no separate SQL migration files for schema creation. AutoMigrate creates tables, adds missing columns, and creates indexes. It never drops columns or tables (safe by design).

**Blueprint deviation**: The blueprint places AutoMigrate in `DatabaseProvider.Boot()`. We instead run it via the `rgo migrate` CLI command. This is deliberate — auto-migrating on every boot is risky in production. Explicit CLI invocation gives the developer control over when schema changes apply.

### 2. Custom Migrations for Changes AutoMigrate Can't Handle
AutoMigrate can't rename columns, drop columns, change types, or add data. For these operations, we need file-based migrations with Up/Down functions. This is standard in Laravel (`php artisan migrate`) and Rails-style frameworks.

### 3. Migration File Pattern
Timestamped Go files in `database/migrations/`:
```
database/migrations/
├── 20260306120000_create_sessions_table.go
├── 20260306120100_add_avatar_to_users.go
└── registry.go  # Registers all migration functions
```

Each migration file exports an `init()` that registers Up/Down functions with the migration registry.

### 4. CLI Commands
- `rgo migrate` — runs AutoMigrate for all registered models, then runs pending file-based migrations
- `rgo make:migration <name>` — generates a timestamped migration file with Up/Down stubs
- `rgo migrate:rollback` — undoes the last batch of applied migrations
- `rgo migrate:status` — lists all migrations and their applied/pending status

### 5. Model Registry
A central `database/models/registry.go` file exports `All()` returning `[]interface{}`. The migrate command calls `db.AutoMigrate(models.All()...)`. This keeps model registration in one place.

## Dependencies

| Dependency | Status | Notes |
|---|---|---|
| Feature #09 — Database Connection | ✅ Done | Provides `*gorm.DB` |
| Feature #10 — CLI Foundation | ✅ Done | Provides Cobra command registration |
| Feature #11 — Models (GORM) | ✅ Done | Provides User, Post models to migrate |

## Questions Resolved

| Question | Answer |
|---|---|
| AutoMigrate or SQL files? | Both: AutoMigrate for model schemas + file-based for custom changes |
| Where do migration files live? | `database/migrations/` per blueprint |
| How to track applied migrations? | `schema_migrations` table with timestamp and batch number |
| How to register models? | `models.All()` function in `database/models/registry.go` |

## Discussion Complete ✅
