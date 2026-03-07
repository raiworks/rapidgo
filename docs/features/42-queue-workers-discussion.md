# 💬 Discussion: Queue Workers / Background Jobs

> **Feature**: `42` — Queue Workers / Background Jobs
> **Status**: 🟢 COMPLETE
> **Branch**: `feature/42-queue-workers`
> **Depends On**: #05 (Service Container), #09 (Database Connection)
> **Date Started**: 2026-03-07
> **Date Completed**: 2026-03-07

---

## Summary

Add a background job queue system to the RapidGo framework, enabling developers to dispatch work (emails, notifications, image processing, API calls) to be processed asynchronously by worker processes. The system supports multiple backends (database, Redis, memory, sync) with automatic retries, failure tracking, and graceful shutdown — following the same provider/container conventions as the rest of the framework.

---

## Functional Requirements

- As a developer, I want to **dispatch jobs** to a queue so that heavy work runs outside the HTTP request cycle
- As a developer, I want to **define job handlers** as simple Go structs so that I can encapsulate job logic
- As a developer, I want to **run workers** via a CLI command (`rapidgo work`) so that I can start processing jobs
- As a developer, I want **multiple queue names** (e.g., "default", "emails", "notifications") so that I can prioritize different types of work
- As a developer, I want **automatic retries** with configurable max attempts so that transient failures don't lose work
- As a developer, I want **failed job tracking** so that I can inspect and retry failed jobs later
- As a developer, I want **graceful shutdown** so that in-progress jobs finish before the worker exits
- As a developer, I want to choose a **queue backend** (database, Redis, memory, sync) via config so that I can match my infrastructure
- As a developer, I want a **sync driver** for local development so that jobs execute immediately during testing without needing external services

---

## Current State / Reference

### What Exists

- **Events system** (`core/events/events.go`) — in-process async dispatch via goroutines. Fire-and-forget, no persistence, no retries. Not suitable for queue work but demonstrates the async dispatch pattern.
- **Redis client** (`redis/go-redis/v9`) — already a dependency. Used by cache system. No shared Redis provider.
- **Database connection** (`database/connection.go`) — GORM-based, supports postgres/mysql/sqlite. Fully reusable for DB-backed queue.
- **Graceful shutdown** (`core/server/server.go`) — `signal.NotifyContext` + per-server shutdown with timeout. Exact pattern needed for workers.
- **CLI commands** (`core/cli/`) — Cobra-based. Workers will follow the same `rootCmd.AddCommand` pattern.
- **Service container** — Singleton/factory pattern with Provider interface (Register + Boot).
- **miniredis** — Already available for Redis testing without external services.

### What Works Well

- Provider/container pattern — clean separation of registration and initialization
- Config via `config.Env()` — consistent env var access
- Graceful shutdown via `signal.NotifyContext` — proven, reusable
- Cobra CLI pattern — consistent command registration
- Service mode system — establishes the pattern for running different process types

### What Needs Improvement

- No shared Redis provider — cache creates its own client. Queue would also need one. Consider a shared Redis provider in this feature or as prerequisite.

---

## Proposed Approach

### Build a framework-native queue system with pluggable backends

**Core concepts**:
1. **Job** — a serializable unit of work with a type name and JSON payload
2. **Queue** — a named FIFO list that stores pending jobs (backed by DB, Redis, or memory)
3. **Worker** — a loop that dequeues and processes jobs, with retry logic and failure tracking
4. **Dispatcher** — the public API used by application code to push jobs onto a queue
5. **Driver** — the backend implementation (database, Redis, memory, sync)

**Architecture layers**:
```
Application Code → Dispatcher → Driver (DB/Redis/Memory/Sync) → Storage
CLI `work` command → Worker Pool → Driver → Job Handler → Application Logic
```

**Key decisions**:
- Custom implementation (no external queue library) — keeps framework self-contained and follows the same pattern as session, cache, and mail
- JSON serialization for job payloads — simple, debuggable, cross-language compatible
- Database as the default driver — works everywhere, no extra infrastructure needed
- Worker runs as a separate CLI command (`rapidgo work`), not embedded in the HTTP server
- No `ModeWorker` in service mode — workers are a separate process, not an HTTP service

### Drivers

| Driver | Use Case | Persistence | Requires |
|---|---|---|---|
| `database` | Production default | ✅ Persistent | GORM + migrations |
| `redis` | High-throughput | ✅ Persistent (Redis persistence) | Redis server |
| `memory` | Testing, development | ❌ In-process only | Nothing |
| `sync` | Local dev, testing | ❌ Executes immediately | Nothing |

---

## Edge Cases & Risks

- [x] **Job serialization** — payloads must be JSON-serializable. Non-serializable types (channels, funcs) will fail at dispatch time, not at processing time. Validate early.
- [x] **Worker crash mid-job** — database driver can use a "reserved" status with timeout to re-queue jobs if the worker dies. Redis uses visibility timeout or BRPOPLPUSH.
- [x] **Database table locking** — `SELECT ... FOR UPDATE SKIP LOCKED` prevents multiple workers from grabbing the same job. PostgreSQL and MySQL both support this. SQLite doesn't support `SKIP LOCKED` — use single-worker mode for SQLite.
- [x] **Redis connection loss** — worker should reconnect with backoff, not crash.
- [x] **Graceful shutdown** — worker must finish current job before exiting, with a hard timeout if the job takes too long.
- [x] **Job handler panics** — recover from panics in job handlers, mark job as failed, continue processing other jobs.
- [x] **Empty queue polling** — don't busy-loop. Use `BRPOP` (Redis) or polling interval with backoff (database).
- [x] **Concurrent workers** — support configurable pool size. Each worker goroutine processes one job at a time.
- [x] **Failed job storage** — record failed jobs with error message and stack trace so developers can debug and retry.
- [x] **Sync driver in tests** — must execute the job immediately and return any error, so tests can assert on outcomes.

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Feature #05 — Service Container | Feature | ✅ Done |
| Feature #09 — Database Connection | Feature | ✅ Done |
| Feature #10 — CLI Foundation | Feature | ✅ Done |
| Feature #12 — Database Migrations | Feature | ✅ Done |
| `redis/go-redis/v9` | External | ✅ Already in go.mod |
| `alicebob/miniredis/v2` | External (test) | ✅ Already in go.mod |
| `gorm.io/gorm` | External | ✅ Already in go.mod |

---

## Open Questions

- [x] **Should we build custom or use a library?** → Custom. Keeps the framework self-contained, follows existing patterns (session, cache, mail all custom). No external queue library needed.
- [x] **Should workers be part of service mode?** → No. Workers don't listen on HTTP ports. Separate `rapidgo work` CLI command.
- [x] **What's the default driver?** → Database. Works everywhere, no extra infrastructure. Redis is the performance upgrade path.
- [x] **Should we create a shared Redis provider?** → Yes, in this feature. Register `*redis.Client` as `"redis"` singleton in container. Cache can migrate to use it later (separate cleanup, not in scope).
- [x] **How do jobs reference their handler?** → By type name string. A global registry maps type names to handler functions, similar to middleware registry pattern.
- [x] **Should we support delayed/scheduled jobs?** → Yes, basic delay support (dispatch with "run after" timestamp). Full cron scheduling is Feature #43.
- [x] **Should we support job priorities?** → No. Use separate named queues instead. Simpler, more predictable.
- [x] **What serialization format?** → JSON. Simple, debuggable, standard library support.

---

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-07 | Custom queue implementation, no external library | Consistent with framework philosophy — cache, session, mail all custom. Keeps dependency tree lean. |
| 2026-03-07 | Database as default driver | Zero extra infrastructure. Works with existing DB. Redis as opt-in upgrade. |
| 2026-03-07 | Separate `work` CLI command, not service mode | Workers don't serve HTTP. Different lifecycle. Same pattern as `migrate`, `db:seed`. |
| 2026-03-07 | JSON payloads | Simple, debuggable, standard library. Binary formats (protobuf, msgpack) overkill for this use case. |
| 2026-03-07 | Job handler registry by type name | Same pattern as middleware registry. String-based lookup. Register at boot time. |
| 2026-03-07 | Create Redis provider as part of this feature | Queue needs Redis client. Registering it in container enables future sharing with cache. |
| 2026-03-07 | Support delayed jobs, not priorities | Delay is common (send email in 5 minutes). Priorities add complexity — use named queues instead. |
| 2026-03-07 | `SELECT FOR UPDATE SKIP LOCKED` for DB driver | Standard pattern for concurrent job processing. SQLite falls back to single-worker. |

---

## Discussion Complete ✅

**Summary**: Build a framework-native queue system with 4 drivers (database, Redis, memory, sync), a CLI worker command, automatic retries, failed job tracking, graceful shutdown, and a shared Redis provider — following the existing provider/container/config patterns.
**Completed**: 2026-03-07
**Next**: Create architecture doc → `42-queue-workers-architecture.md`
