# Feature #53 — Database Read/Write Splitting: Discussion

> **Status**: 🟡 IN PROGRESS
> **Depends On**: #09 (Database Connection)
> **Branch**: `docs/53-db-read-write-splitting`

---

## What Problem Does This Solve?

Production applications with high read traffic can overwhelm a single database server. Most web workloads are read-heavy (80-90% reads), so distributing reads to one or more replicas while directing writes to the primary server significantly improves throughput and reduces latency.

Currently, RapidGo registers a single `*gorm.DB` instance as the `"db"` singleton. All queries — reads and writes — flow through the same connection. There is no mechanism to route read queries to a separate replica.

## What Does This Feature Add?

A `Resolver` struct in the `database` package that holds two `*gorm.DB` connections — one for writes (primary) and one for reads (replica). Developers explicitly choose which connection to use via `Writer()` and `Reader()` methods.

### Key Design Decisions

1. **Explicit routing, not implicit** — The developer calls `resolver.Writer()` or `resolver.Reader()` rather than the framework guessing based on query type. This avoids subtle bugs (e.g., reading stale data from a replica immediately after a write).

2. **Opt-in via env var** — When `DB_READ_HOST` is set, a separate read connection is established. When unset, `Reader()` returns the same connection as `Writer()` — zero overhead, no code changes required.

3. **Backward compatible** — The existing `"db"` singleton remains unchanged and always returns the writer. Existing code continues to work unmodified.

4. **No new dependencies** — Uses existing GORM and driver packages. No external plugins.

## What's Out of Scope?

- **Automatic query routing** — No parsing SQL to detect SELECT vs INSERT. Developers choose explicitly.
- **Multiple read replicas with load balancing** — Single reader connection. Multiple replicas can be fronted by a load balancer (e.g., PgBouncer, ProxySQL) at the infrastructure level.
- **Replication lag detection** — No staleness checks. The developer decides when to read from primary vs replica.
- **Query builder changes** — No modifications to `database/querybuilder/`.
- **Health check changes** — The existing `/health/ready` endpoint continues to ping the writer. Developers can add custom health endpoints for replicas if needed.

## How Will Developers Use It?

```go
// In a service — resolve from container
resolver := container.MustMake[*database.Resolver](c, "db.resolver")

// Reads go to the replica
var users []models.User
resolver.Reader().Find(&users)

// Writes go to the primary
resolver.Writer().Create(&newUser)

// Transactions always use the writer
database.WithTransaction(resolver.Writer(), func(tx *gorm.DB) error {
    // ...
    return nil
})
```

When `DB_READ_HOST` is not set, both `Reader()` and `Writer()` return the same `*gorm.DB`, so the code works identically in development (single DB) and production (split DB).
