# Feature #53 — Database Read/Write Splitting: Architecture

> **Status**: ✅ SHIPPED
> **Depends On**: #09 (Database Connection)
> **Branch**: `docs/53-db-read-write-splitting`

---

## Package Location

`database/` — extends the existing database package. No new packages created.

## New File: `database/resolver.go`

### Structs

```go
// Resolver holds separate database connections for write and read operations.
type Resolver struct {
    writer *gorm.DB
    reader *gorm.DB
}
```

### Functions

```go
// NewResolver creates a Resolver with the given writer and reader connections.
// If reads should go to the same database, pass the same *gorm.DB for both.
func NewResolver(writer, reader *gorm.DB) *Resolver

// Writer returns the write (primary) database connection.
func (r *Resolver) Writer() *gorm.DB

// Reader returns the read (replica) database connection.
// Returns the writer if no separate reader was configured.
func (r *Resolver) Reader() *gorm.DB
```

## Modified File: `database/connection.go`

### New Function

```go
// NewReadDBConfig reads read-replica configuration from DB_READ_* environment
// variables. Each setting falls back to the corresponding DB_* value, then to
// the same defaults used by NewDBConfig.
func NewReadDBConfig() DBConfig
```

### Environment Variables

| Variable | Fallback | Default |
|---|---|---|
| `DB_READ_DRIVER` | `DB_DRIVER` | `""` |
| `DB_READ_HOST` | `DB_HOST` | `"localhost"` |
| `DB_READ_PORT` | `DB_PORT` | `"5432"` |
| `DB_READ_NAME` | `DB_NAME` | `"rapidgo_dev"` |
| `DB_READ_USER` | `DB_USER` | `""` |
| `DB_READ_PASSWORD` | `DB_PASSWORD` | `""` |
| `DB_READ_SSL_MODE` | `DB_SSL_MODE` | `"disable"` |
| `DB_READ_MAX_OPEN_CONNS` | `DB_MAX_OPEN_CONNS` | `25` |
| `DB_READ_MAX_IDLE_CONNS` | `DB_MAX_IDLE_CONNS` | `10` |
| `DB_READ_CONN_MAX_LIFETIME` | `DB_CONN_MAX_LIFETIME` | `5` (minutes) |
| `DB_READ_CONN_MAX_IDLE_TIME` | `DB_CONN_MAX_IDLE_TIME` | `3` (minutes) |

### Fallback Chain

Each `DB_READ_*` variable follows the same pattern:

```
DB_READ_HOST → DB_HOST → "localhost"
```

This means:
- If only `DB_HOST` is set, reads use the same host as writes (no splitting).
- If `DB_READ_HOST` is set, reads go to the replica.
- Pool settings can be tuned independently for the read connection.

## Modified File: `app/providers/database_provider.go`

### New Singleton

```go
// Inside Register(), after the existing "db" singleton:
c.Singleton("db.resolver", func(c *container.Container) interface{} {
    writer := c.Make("db").(*gorm.DB)
    if config.Env("DB_READ_HOST", "") == "" {
        return database.NewResolver(writer, writer)
    }
    reader, err := database.ConnectWithConfig(database.NewReadDBConfig())
    if err != nil {
        panic("read database connection failed: " + err.Error())
    }
    return database.NewResolver(writer, reader)
})
```

### Behavior

| `DB_READ_HOST` | Result |
|---|---|
| Not set / empty | `Reader()` returns the same `*gorm.DB` as `Writer()` |
| Set to a host | Separate `*gorm.DB` connection established for reads |

### Container Keys

| Key | Type | Description |
|---|---|---|
| `"db"` | `*gorm.DB` | Writer (primary) — unchanged |
| `"db.resolver"` | `*database.Resolver` | Writer + Reader pair |

## No Changes

- `database/transaction.go` — unchanged. `WithTransaction` accepts `*gorm.DB`; the caller passes `resolver.Writer()`.
- `core/health/health.go` — unchanged. Continues to ping the writer via `dbFn()`.
- `database/querybuilder/` — empty (`.gitkeep` only), not affected.
- `database/models/` — unchanged.
- No new dependencies.
- No migrations.

## Usage Pattern

```go
// Service constructor — accepts Resolver
type OrderService struct {
    db *database.Resolver
}

func NewOrderService(resolver *database.Resolver) *OrderService {
    return &OrderService{db: resolver}
}

// Read-heavy operation → use replica
func (s *OrderService) List() ([]Order, error) {
    var orders []Order
    return orders, s.db.Reader().Find(&orders).Error
}

// Write operation → use primary
func (s *OrderService) Create(o *Order) error {
    return s.db.Writer().Create(o).Error
}

// Transaction → always primary
func (s *OrderService) Transfer(fromID, toID uint) error {
    return database.WithTransaction(s.db.Writer(), func(tx *gorm.DB) error {
        // all ops use tx
        return nil
    })
}
```
