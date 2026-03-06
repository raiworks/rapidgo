# 🧪 Test Plan: Database Connection

> **Feature**: `09` — Database Connection
> **Tasks**: [`09-database-tasks.md`](09-database-tasks.md)
> **Date**: 2026-03-06

---

## Acceptance Criteria

- [ ] `NewDBConfig()` reads all `DB_*` environment variables correctly
- [ ] `NewDBConfig()` uses sensible fallback defaults when env vars are unset
- [ ] `DSN()` returns correct PostgreSQL connection string format
- [ ] `DSN()` returns correct MySQL connection string format
- [ ] `DSN()` returns the database name/path for SQLite
- [ ] `DSN()` returns empty string for unsupported drivers
- [ ] `ConnectWithConfig()` succeeds with SQLite `:memory:`
- [ ] `ConnectWithConfig()` returns error for unsupported driver
- [ ] Connection pool settings are applied correctly
- [ ] `DatabaseProvider` implements `Provider` interface
- [ ] `DatabaseProvider.Register()` registers `"db"` binding in container
- [ ] `DatabaseProvider.Boot()` is a no-op
- [ ] All tests pass with `go test ./database/...`
- [ ] All tests pass with `go test ./...` (full regression)
- [ ] `go vet ./...` reports no issues

---

## Test Cases

### TC-01: NewDBConfig reads all env vars

**File**: `database/database_test.go`
**Function**: `TestNewDBConfig_ReadsEnvVars`

| Step | Action | Expected |
|---|---|---|
| 1 | Set all `DB_*` env vars via `t.Setenv()` | Env vars set |
| 2 | Call `NewDBConfig()` | Returns `DBConfig` |
| 3 | Assert `cfg.Driver == "postgres"` | Match |
| 4 | Assert `cfg.Host == "testhost"` | Match |
| 5 | Assert `cfg.Port == "3306"` | Match |
| 6 | Assert `cfg.Name == "testdb"` | Match |
| 7 | Assert `cfg.User == "testuser"` | Match |
| 8 | Assert `cfg.Password == "testpass"` | Match |
| 9 | Assert `cfg.SSLMode == "require"` | Match |
| 10 | Assert `cfg.MaxOpenConns == 50` | Match |
| 11 | Assert `cfg.MaxIdleConns == 20` | Match |
| 12 | Assert `cfg.ConnMaxLifetime == 10 * time.Minute` | Match |
| 13 | Assert `cfg.ConnMaxIdleTime == 7 * time.Minute` | Match |

---

### TC-02: NewDBConfig uses defaults

**File**: `database/database_test.go`
**Function**: `TestNewDBConfig_Defaults`

| Step | Action | Expected |
|---|---|---|
| 1 | Do NOT set any `DB_*` env vars | Clean environment |
| 2 | Call `NewDBConfig()` | Returns `DBConfig` with defaults |
| 3 | Assert `cfg.Driver == ""` | Empty (forces explicit config) |
| 4 | Assert `cfg.Host == "localhost"` | Default |
| 5 | Assert `cfg.Port == "5432"` | Default |
| 6 | Assert `cfg.Name == "rgo_dev"` | Default |
| 7 | Assert `cfg.MaxOpenConns == 25` | Default |
| 8 | Assert `cfg.MaxIdleConns == 10` | Default |
| 9 | Assert `cfg.ConnMaxLifetime == 5 * time.Minute` | Default |
| 10 | Assert `cfg.ConnMaxIdleTime == 3 * time.Minute` | Default |

---

### TC-03: DSN postgres format

**File**: `database/database_test.go`
**Function**: `TestDSN_Postgres`

| Step | Action | Expected |
|---|---|---|
| 1 | Create `DBConfig{Driver: "postgres", Host: "db.example.com", User: "admin", Password: "s3cret", Name: "myapp", Port: "5432", SSLMode: "require"}` | Config created |
| 2 | Call `cfg.DSN()` | Returns `"host=db.example.com user=admin password=s3cret dbname=myapp port=5432 sslmode=require"` |

---

### TC-04: DSN mysql format

**File**: `database/database_test.go`
**Function**: `TestDSN_MySQL`

| Step | Action | Expected |
|---|---|---|
| 1 | Create `DBConfig{Driver: "mysql", Host: "db.example.com", User: "root", Password: "pass", Name: "myapp", Port: "3306"}` | Config created |
| 2 | Call `cfg.DSN()` | Returns `"root:pass@tcp(db.example.com:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local"` |

---

### TC-05: DSN sqlite format

**File**: `database/database_test.go`
**Function**: `TestDSN_SQLite`

| Step | Action | Expected |
|---|---|---|
| 1 | Create `DBConfig{Driver: "sqlite", Name: ":memory:"}` | Config created |
| 2 | Call `cfg.DSN()` | Returns `":memory:"` |

---

### TC-06: DSN unsupported driver

**File**: `database/database_test.go`
**Function**: `TestDSN_UnsupportedDriver`

| Step | Action | Expected |
|---|---|---|
| 1 | Create `DBConfig{Driver: "oracle"}` | Config created |
| 2 | Call `cfg.DSN()` | Returns `""` (empty string) |

---

### TC-07: ConnectWithConfig succeeds with SQLite in-memory

**File**: `database/database_test.go`
**Function**: `TestConnectWithConfig_SQLiteMemory`

| Step | Action | Expected |
|---|---|---|
| 1 | Create `DBConfig{Driver: "sqlite", Name: ":memory:", MaxOpenConns: 1, MaxIdleConns: 1, ConnMaxLifetime: time.Minute, ConnMaxIdleTime: time.Minute}` | Config created |
| 2 | Call `ConnectWithConfig(cfg)` | Returns `*gorm.DB`, no error |
| 3 | Assert `db != nil` | Not nil |
| 4 | Get underlying `sql.DB` via `db.DB()` | No error |
| 5 | Call `sqlDB.Ping()` | No error — connection is alive |

---

### TC-08: ConnectWithConfig fails with unsupported driver

**File**: `database/database_test.go`
**Function**: `TestConnectWithConfig_UnsupportedDriver`

| Step | Action | Expected |
|---|---|---|
| 1 | Create `DBConfig{Driver: "oracle"}` | Config created |
| 2 | Call `ConnectWithConfig(cfg)` | Returns nil `*gorm.DB`, non-nil error |
| 3 | Assert error message contains `"unsupported DB_DRIVER: oracle"` | Match |

---

### TC-09: Pool settings are applied

**File**: `database/database_test.go`
**Function**: `TestConnectWithConfig_PoolSettings`

| Step | Action | Expected |
|---|---|---|
| 1 | Create `DBConfig{Driver: "sqlite", Name: ":memory:", MaxOpenConns: 42, MaxIdleConns: 7, ConnMaxLifetime: 10 * time.Minute, ConnMaxIdleTime: 5 * time.Minute}` | Config created |
| 2 | Call `ConnectWithConfig(cfg)` | Returns `*gorm.DB`, no error |
| 3 | Get underlying `sql.DB` via `db.DB()` | No error |
| 4 | Call `sqlDB.Stats()` | Returns `sql.DBStats` |
| 5 | Assert `stats.MaxOpenConnections == 42` | Pool setting applied |

---

### TC-10: DatabaseProvider implements Provider interface (compile-time)

**File**: `app/providers/providers_test.go`
**Declaration**: `var _ container.Provider = (*DatabaseProvider)(nil)`

| Step | Action | Expected |
|---|---|---|
| 1 | Compile-time interface assertion | Compiles without error |

---

### TC-11: DatabaseProvider.Register registers "db" binding

**File**: `app/providers/providers_test.go`
**Function**: `TestDatabaseProvider_RegistersBinding`

| Step | Action | Expected |
|---|---|---|
| 1 | Create new `container.New()` | Empty container |
| 2 | Call `DatabaseProvider{}.Register(c)` | Registers singleton factory |
| 3 | Assert `c.Has("db") == true` | Binding registered |

---

### TC-12: Full bootstrap with DatabaseProvider (SQLite integration)

**File**: `app/providers/providers_test.go`
**Function**: `TestDatabaseProvider_FullBootstrap`

| Step | Action | Expected |
|---|---|---|
| 1 | Set `DB_DRIVER=sqlite`, `DB_NAME=:memory:` via `t.Setenv()` | Env configured for SQLite |
| 2 | Create `app.New()`, register ConfigProvider + DatabaseProvider | App created |
| 3 | Call `app.Boot()` | No panic |
| 4 | Call `container.Make("db")` | Returns `*gorm.DB` |
| 5 | Assert result is not nil | Connection established |

---

## Test Matrix Summary

| Test Case | File | Type | Needs DB |
|---|---|---|---|
| TC-01 | `database/database_test.go` | Unit | No |
| TC-02 | `database/database_test.go` | Unit | No |
| TC-03 | `database/database_test.go` | Unit | No |
| TC-04 | `database/database_test.go` | Unit | No |
| TC-05 | `database/database_test.go` | Unit | No |
| TC-06 | `database/database_test.go` | Unit | No |
| TC-07 | `database/database_test.go` | Integration | SQLite `:memory:` |
| TC-08 | `database/database_test.go` | Unit | No |
| TC-09 | `database/database_test.go` | Integration | SQLite `:memory:` |
| TC-10 | `app/providers/providers_test.go` | Compile-time | No |
| TC-11 | `app/providers/providers_test.go` | Unit | No |
| TC-12 | `app/providers/providers_test.go` | Integration | SQLite `:memory:` |
