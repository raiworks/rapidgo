# 📋 Changelog: Queue Workers / Background Jobs

> **Feature**: `42` — Queue Workers / Background Jobs
> **Status**: � Complete

---

## [Unreleased] — Feature #42

### Added

- `core/queue/queue.go` — Job struct, Driver interface, HandlerFunc, handler registry, Dispatcher
- `core/queue/memory.go` — In-memory queue driver (mutex-protected slices)
- `core/queue/sync.go` — Synchronous driver (immediate execution, no worker needed)
- `core/queue/database.go` — Database driver (GORM, FOR UPDATE SKIP LOCKED, SQLite fallback)
- `core/queue/redis.go` — Redis driver (go-redis lists, JSON serialization)
- `core/queue/worker.go` — Worker pool (goroutine-per-worker, retry, panic recovery, timeout, graceful shutdown)
- `core/queue/queue_test.go` — 24 test cases covering all components
- `app/providers/redis_provider.go` — RedisProvider (shared `*redis.Client` singleton)
- `app/providers/queue_provider.go` — QueueProvider (Dispatcher singleton, driver switch)
- `core/cli/work.go` — `rapidgo work` CLI command with --queues, --workers, --timeout flags
- `app/jobs/example_job.go` — Example job handler with RegisterJobs()
- `database/migrations/20260307000001_create_jobs_tables.go` — jobs + failed_jobs migration

### Changed

- `core/cli/root.go` — Added `workCmd`, `RedisProvider`, `QueueProvider` registration
- `.env` — Added QUEUE_DRIVER, QUEUE_DEFAULT, QUEUE_TABLE, QUEUE_FAILED_TABLE, QUEUE_MAX_ATTEMPTS, QUEUE_RETRY_DELAY, QUEUE_TIMEOUT, QUEUE_POLL_INTERVAL

### Configuration

| Key | Default | Description |
|---|---|---|
| QUEUE_DRIVER | database | Backend: database, redis, memory, sync |
| QUEUE_DEFAULT | default | Default queue name |
| QUEUE_TABLE | jobs | Database table for pending jobs |
| QUEUE_FAILED_TABLE | failed_jobs | Database table for failed jobs |
| QUEUE_MAX_ATTEMPTS | 3 | Max retry attempts |
| QUEUE_RETRY_DELAY | 30 | Seconds before retry |
| QUEUE_TIMEOUT | 60 | Max job processing seconds |
| QUEUE_POLL_INTERVAL | 3 | Seconds between polls |
