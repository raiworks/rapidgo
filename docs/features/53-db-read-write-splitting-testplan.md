# Feature #53 — Database Read/Write Splitting: Test Plan

> **Status**: 🟡 IN PROGRESS
> **Depends On**: #09 (Database Connection)
> **Branch**: `docs/53-db-read-write-splitting`

---

## Test File: `database/resolver_test.go`

### Test Cases (10 total)

| ID | Test Name | What It Verifies |
|---|---|---|
| T01 | `TestNewResolver` | `NewResolver(w, r)` stores writer and reader correctly |
| T02 | `TestWriter` | `Writer()` returns the writer `*gorm.DB` |
| T03 | `TestReader` | `Reader()` returns the reader `*gorm.DB` |
| T04 | `TestReaderFallback` | When `NewResolver(w, w)`, `Reader()` returns same instance as `Writer()` |
| T05 | `TestNewReadDBConfig_Explicit` | `NewReadDBConfig()` reads `DB_READ_*` env vars |
| T06 | `TestNewReadDBConfig_Fallback` | Without `DB_READ_*`, falls back to `DB_*` values |
| T07 | `TestNewReadDBConfig_Defaults` | Without any env vars, returns same defaults as `NewDBConfig()` |
| T08 | `TestResolverWithSQLite` | Create two SQLite in-memory connections, verify `Writer()` and `Reader()` are distinct and functional |
| T09 | `TestResolverSameConnection` | Create one SQLite connection, pass as both writer and reader, verify queries work |
| T10 | `TestNewReadDBConfig_PoolSettings` | `DB_READ_MAX_OPEN_CONNS` etc. are read independently from `DB_MAX_OPEN_CONNS` |

### Test Approach

- **T01–T04**: Pure unit tests — create mock `*gorm.DB` instances (SQLite in-memory), verify pointer identity.
- **T05–T07, T10**: Env var tests — use `t.Setenv()` to set/unset `DB_READ_*` and `DB_*` values, call `NewReadDBConfig()`, assert each field.
- **T08–T09**: Integration-style — use `ConnectWithConfig` with SQLite driver to create real GORM connections, execute queries through `Writer()` and `Reader()`.

### Acceptance Criteria

1. All 10 tests pass with `go test ./database/...`
2. No changes to existing `database_test.go` tests
3. Full regression passes: `go test ./...`
4. Binary builds clean: `go build -o bin/rapidgo.exe ./cmd`
