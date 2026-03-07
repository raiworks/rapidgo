# Feature #53 — Database Read/Write Splitting: Tasks

> **Status**: 🟡 IN PROGRESS
> **Depends On**: #09 (Database Connection)
> **Branch**: `docs/53-db-read-write-splitting`

---

## Task List

### T1 — Add `NewReadDBConfig` to `database/connection.go`

- [ ] Add `NewReadDBConfig() DBConfig` function
- [ ] Each field reads `DB_READ_*` env var, falling back to `config.Env("DB_*", default)`
- [ ] Same field set as `NewDBConfig`: Driver, Host, Port, Name, User, Password, SSLMode, MaxOpenConns, MaxIdleConns, ConnMaxLifetime, ConnMaxIdleTime
- [ ] No changes to existing functions

### T2 — Create `database/resolver.go`

- [ ] Create new file `database/resolver.go`
- [ ] Define `Resolver` struct with unexported `writer` and `reader` fields (both `*gorm.DB`)
- [ ] Implement `NewResolver(writer, reader *gorm.DB) *Resolver`
- [ ] Implement `Writer() *gorm.DB` — returns the writer connection
- [ ] Implement `Reader() *gorm.DB` — returns the reader connection

### T3 — Register `"db.resolver"` in `DatabaseProvider`

- [ ] Modify `app/providers/database_provider.go`
- [ ] Add `"db.resolver"` singleton registration inside `Register()`
- [ ] When `DB_READ_HOST` is empty: `NewResolver(writer, writer)` — reader mirrors writer
- [ ] When `DB_READ_HOST` is set: connect via `ConnectWithConfig(NewReadDBConfig())`, pass as reader
- [ ] Panic on read connection failure (consistent with writer behavior)
- [ ] Add `config` and `database` imports as needed

### T4 — Write tests in `database/resolver_test.go`

- [ ] Create `database/resolver_test.go`
- [ ] 10 test cases (T01–T10) per test plan
- [ ] Verify NewResolver, Writer, Reader, NewReadDBConfig, fallback behavior
- [ ] Use SQLite in-memory for integration-style tests (no external DB required)
- [ ] Use `t.Setenv` for env var manipulation

### T5 — Verify full regression

- [ ] Run `go test ./...` — all packages pass
- [ ] Run `go build -o bin/rapidgo.exe ./cmd` — binary builds clean
