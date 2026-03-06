# 💬 Discussion: Database Seeding

> **Feature**: `13` — Database Seeding
> **Status**: 🟢 COMPLETE
> **Date**: 2026-03-06

---

## What Are We Building?

A seeder system that populates the database with initial/test data. Seeders live in `database/seeders/`, are registered via a central registry, and are invoked via the `rgo db:seed` CLI command. The system supports running all seeders or a specific one by name.

## Blueprint References

The blueprint specifies (lines 3032–3065):

1. **Directory**: `database/seeders/` (line 72)
2. **CLI command**: `framework db:seed` (line 496)
3. **Pattern**: Individual seed functions (`SeedUsers(db)`) + `RunAll(db)` orchestrator
4. **Technique**: `db.FirstOrCreate()` for idempotent seeding
5. **Example**: Seed admin and regular users with hashed passwords

## Scope for Feature #13

### In Scope
- Seeder interface — `Seeder` with `Name()` and `Seed(db)` methods
- Seeder registry — `Register()` and `RunAll()` in `database/seeders/`
- `UserSeeder` — sample seeder creating admin + regular user
- `rgo db:seed` CLI command — runs all registered seeders
- `--seeder` flag — run a specific seeder by name

### Out of Scope (deferred)
- Password hashing — `helpers.HashPassword()` is Feature #19/#22. Seeders store plaintext for now; will be updated when crypto ships.
- `make:seeder` command — scaffolding is a future feature
- Post seeder — can be added later as an example

## Key Design Decisions

### 1. Seeder Interface vs Bare Functions
The blueprint uses bare functions (`SeedUsers(db)`). We use a `Seeder` interface instead:
```go
type Seeder interface {
    Name() string
    Seed(db *gorm.DB) error
}
```
This enables the `--seeder` flag (run by name), consistent registration, and better testability. Minor adaptation from blueprint.

### 2. Idempotent Seeding with FirstOrCreate
Following the blueprint, seeders use `db.FirstOrCreate()` so running `db:seed` multiple times doesn't create duplicates.

### 3. Plaintext Passwords (Temporary)
The blueprint calls `helpers.HashPassword()`, but Helpers (#19) and Crypto (#22) aren't built yet. Passwords stored as plaintext in seeders for now — clearly documented. Will be updated when those features ship.

### 4. Registry Pattern (Consistent with Migrations)
Same pattern as Feature #12: seeders self-register via `Register()`, central `RunAll()` orchestrates. Consistent API across the data layer.

## Dependencies

| Dependency | Status | Notes |
|---|---|---|
| Feature #09 — Database Connection | ✅ Done | Provides `*gorm.DB` |
| Feature #10 — CLI Foundation | ✅ Done | Provides Cobra command registration |
| Feature #11 — Models (GORM) | ✅ Done | Provides User, Post models |

## Discussion Complete ✅
