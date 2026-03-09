---
title: "Migrations"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Migrations

## Abstract

This document covers database schema management using GORM's
auto-migration and the CLI `migrate` command.

## Table of Contents

1. [Terminology](#1-terminology)
2. [GORM AutoMigrate](#2-gorm-automigrate)
3. [CLI Command](#3-cli-command)
4. [Migration Strategy](#4-migration-strategy)
5. [Security Considerations](#5-security-considerations)
6. [References](#6-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Migration** — A versioned change to the database schema.
- **AutoMigrate** — GORM's built-in method that creates or updates
  tables to match model struct definitions.

## 2. GORM AutoMigrate

GORM can automatically create and alter tables based on model
definitions:

```go
func (p *DatabaseProvider) Boot(c *container.Container) {
    db := container.MustMake[*gorm.DB](c, "db")
    db.AutoMigrate(
        &models.User{},
        &models.Post{},
        &session.SessionRecord{},
    )
}
```

`AutoMigrate`:
- Creates tables if they don't exist
- Adds missing columns
- Adds missing indexes
- Does **NOT** delete columns or change column types

### Limitations

AutoMigrate is suitable for development but has limitations:

| Capability | Supported |
|-----------|-----------|
| Create tables | Yes |
| Add columns | Yes |
| Add indexes | Yes |
| Drop columns | No |
| Change column types | No |
| Rename columns | No |
| Complex alterations | No |

For production schema changes that require column drops or type
changes, use manual SQL migrations.

## 3. CLI Command

Run migrations via the CLI:

```text
framework migrate
```

This calls the boot phase of the `DatabaseProvider`, which triggers
`AutoMigrate` for all registered models.

## 4. Migration Strategy

### Development

Use GORM AutoMigrate freely — it safely handles additive changes.

### Production

- **Additive changes** (new tables, columns, indexes): AutoMigrate
  is safe.
- **Destructive changes** (drop columns, rename, type changes):
  Write manual SQL migration scripts in `database/migrations/` and
  apply them before deploying.
- **Always back up** the database before running migrations in
  production.

### Migration Files

For manual migrations, place SQL files in `database/migrations/`:

```text
database/
└── migrations/
    ├── 001_create_users.sql
    ├── 002_create_posts.sql
    └── 003_add_role_to_users.sql
```

## 5. Security Considerations

- Migrations **MUST** not run with superuser/admin database
  credentials in production. Use a dedicated migration user with
  schema-alter permissions only.
- Always review generated SQL before applying to production.
- Migration scripts **MUST NOT** contain sensitive data (passwords,
  keys).

## 6. References

- [Database](database.md)
- [Models](models.md)
- [Seeders](seeders.md)
- [CLI Overview](../cli/cli-overview.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
