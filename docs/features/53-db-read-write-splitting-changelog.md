# Feature #53 — Database Read/Write Splitting: Changelog

> **Status**: ✅ SHIPPED
> **Depends On**: #09 (Database Connection)
> **Branch**: `docs/53-db-read-write-splitting`

---

## Files Changed

| File | Action | Description |
|---|---|---|
| `database/resolver.go` | **NEW** | `Resolver` struct, `NewResolver`, `Writer`, `Reader` |
| `database/resolver_test.go` | **NEW** | 10 test cases (T01–T10) |
| `database/connection.go` | MODIFIED | Add `NewReadDBConfig()` function |
| `app/providers/database_provider.go` | MODIFIED | Register `"db.resolver"` singleton |

## New Environment Variables

| Variable | Fallback | Default | Purpose |
|---|---|---|---|
| `DB_READ_DRIVER` | `DB_DRIVER` | `""` | Read replica database driver |
| `DB_READ_HOST` | `DB_HOST` | `"localhost"` | Read replica host (opt-in trigger) |
| `DB_READ_PORT` | `DB_PORT` | `"5432"` | Read replica port |
| `DB_READ_NAME` | `DB_NAME` | `"rapidgo_dev"` | Read replica database name |
| `DB_READ_USER` | `DB_USER` | `""` | Read replica user |
| `DB_READ_PASSWORD` | `DB_PASSWORD` | `""` | Read replica password |
| `DB_READ_SSL_MODE` | `DB_SSL_MODE` | `"disable"` | Read replica SSL mode |
| `DB_READ_MAX_OPEN_CONNS` | `DB_MAX_OPEN_CONNS` | `25` | Read replica max open connections |
| `DB_READ_MAX_IDLE_CONNS` | `DB_MAX_IDLE_CONNS` | `10` | Read replica max idle connections |
| `DB_READ_CONN_MAX_LIFETIME` | `DB_CONN_MAX_LIFETIME` | `5` (min) | Read replica connection max lifetime |
| `DB_READ_CONN_MAX_IDLE_TIME` | `DB_CONN_MAX_IDLE_TIME` | `3` (min) | Read replica connection max idle time |

## New Container Keys

| Key | Type | Description |
|---|---|---|
| `"db.resolver"` | `*database.Resolver` | Writer + Reader connection pair |

## Dependencies

No new dependencies. Uses existing `gorm.io/gorm` and driver packages.

## Migrations

None.

## Breaking Changes

None. The existing `"db"` singleton is unchanged. All existing code continues to work.
