# 🏗️ Architecture: Queue Workers / Background Jobs

> **Feature**: `42` — Queue Workers / Background Jobs
> **Discussion**: [`42-queue-workers-discussion.md`](42-queue-workers-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-07

---

## Overview

A framework-native queue system with pluggable backends (database, Redis, memory, sync), a CLI worker command, automatic retries, failed job tracking, and graceful shutdown. Jobs are dispatched from application code via a `Dispatcher`, persisted by a `Driver`, and processed by workers started via `rapidgo work`. The system follows existing container/provider/config patterns.

---

## File Structure

```
core/
├── queue/
│   ├── queue.go                # Job, Driver interface, Dispatcher, handler registry
│   ├── database.go             # Database driver (GORM-backed)
│   ├── redis.go                # Redis driver (go-redis/v9)
│   ├── memory.go               # In-memory driver (channel-based)
│   ├── sync.go                 # Synchronous driver (immediate execution)
│   ├── worker.go               # Worker pool with graceful shutdown
│   └── queue_test.go           # Tests for all components
├── cli/
│   ├── work.go                 # NEW — `rapidgo work` command
│   └── root.go                 # MODIFY — register QueueProvider conditionally

app/
├── providers/
│   ├── queue_provider.go       # NEW — QueueProvider (registers queue in container)
│   └── redis_provider.go       # NEW — RedisProvider (registers *redis.Client)
├── jobs/
│   └── example_job.go          # NEW — Example job handler

database/
├── migrations/
│   └── XXXXXX_create_jobs_table.go  # NEW — jobs + failed_jobs migration

.env                            # MODIFY — add queue env vars
```

**Total**: 10 new files, 2 modified files

---

## Data Model

### `jobs` Table

| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uint64 | PK, auto | Unique identifier |
| queue | varchar(255) | NOT NULL, INDEX | Queue name (e.g., "default", "emails") |
| type | varchar(255) | NOT NULL | Job handler type name |
| payload | text | NOT NULL | JSON-encoded job data |
| attempts | uint | NOT NULL, DEFAULT 0 | Number of times this job has been attempted |
| max_attempts | uint | NOT NULL, DEFAULT 3 | Maximum retry attempts |
| available_at | datetime | NOT NULL, INDEX | When this job becomes available for processing |
| reserved_at | datetime | NULL | When a worker reserved this job (NULL = available) |
| created_at | datetime | NOT NULL | When the job was dispatched |

### `failed_jobs` Table

| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uint64 | PK, auto | Unique identifier |
| queue | varchar(255) | NOT NULL | Queue name |
| type | varchar(255) | NOT NULL | Job handler type name |
| payload | text | NOT NULL | JSON-encoded job data |
| error | text | NOT NULL | Error message + stack trace |
| failed_at | datetime | NOT NULL | When the job failed |

### Relationships

None — jobs are standalone records. No foreign keys to other application tables.

---

## Component Design

### Job

**Package**: `core/queue`
**File**: `queue.go`

```go
// Job represents a unit of work to be processed by a worker.
type Job struct {
    ID          uint64
    Queue       string
    Type        string
    Payload     json.RawMessage
    Attempts    uint
    MaxAttempts uint
    AvailableAt time.Time
    ReservedAt  *time.Time
    CreatedAt   time.Time
}
```

### HandlerFunc & Registry

**Package**: `core/queue`
**File**: `queue.go`

```go
// HandlerFunc processes a job. Receives the raw JSON payload.
type HandlerFunc func(ctx context.Context, payload json.RawMessage) error

// RegisterHandler maps a type name to a handler function.
func RegisterHandler(typeName string, handler HandlerFunc)

// ResolveHandler returns the handler for a type name, or nil.
func ResolveHandler(typeName string) HandlerFunc

// ResetHandlers clears the registry. For testing only.
func ResetHandlers()
```

Pattern: Identical to `middleware.RegisterAlias` / `middleware.Resolve`.

### Driver Interface

**Package**: `core/queue`
**File**: `queue.go`

```go
// Driver is the storage backend for the queue system.
type Driver interface {
    // Push adds a job to the queue.
    Push(ctx context.Context, job *Job) error

    // Pop retrieves and reserves the next available job from the given queue.
    // Returns nil, nil if no job is available.
    Pop(ctx context.Context, queue string) (*Job, error)

    // Delete removes a completed job.
    Delete(ctx context.Context, job *Job) error

    // Release puts a reserved job back into the queue for retry.
    // The job's AvailableAt is set to the given delay from now.
    Release(ctx context.Context, job *Job, delay time.Duration) error

    // Fail moves a job to the failed jobs storage.
    Fail(ctx context.Context, job *Job, jobErr error) error

    // Size returns the number of pending jobs in a queue.
    Size(ctx context.Context, queue string) (int64, error)
}
```

### Dispatcher

**Package**: `core/queue`
**File**: `queue.go`

```go
// Dispatcher is the public API for dispatching jobs.
type Dispatcher struct {
    driver Driver
}

// NewDispatcher creates a dispatcher with the given driver.
func NewDispatcher(driver Driver) *Dispatcher

// Dispatch pushes a job onto the named queue.
func (d *Dispatcher) Dispatch(ctx context.Context, queue, typeName string, payload interface{}) error

// DispatchDelayed pushes a job onto the queue with a delay.
func (d *Dispatcher) DispatchDelayed(ctx context.Context, queue, typeName string, payload interface{}, delay time.Duration) error

// Driver returns the underlying driver (for Size checks, etc.).
func (d *Dispatcher) Driver() Driver
```

`Dispatch` serializes the payload to JSON, constructs a `Job`, and calls `driver.Push()`.

### Database Driver

**Package**: `core/queue`
**File**: `database.go`

```go
// DatabaseDriver implements Driver using GORM.
type DatabaseDriver struct {
    db          *gorm.DB
    table       string // from QUEUE_TABLE, default "jobs"
    failedTable string // from QUEUE_FAILED_TABLE, default "failed_jobs"
}

func NewDatabaseDriver(db *gorm.DB, table, failedTable string) *DatabaseDriver
```

**GORM model structs** (internal to database.go, not exported):

```go
// jobModel is the GORM model for the jobs table.
type jobModel struct {
    ID          uint64          `gorm:"primaryKey;autoIncrement"`
    Queue       string          `gorm:"size:255;not null;index"`
    Type        string          `gorm:"size:255;not null"`
    Payload     string          `gorm:"type:text;not null"`
    Attempts    uint            `gorm:"not null;default:0"`
    MaxAttempts uint            `gorm:"not null;default:3"`
    AvailableAt time.Time       `gorm:"not null;index"`
    ReservedAt  *time.Time      `gorm:"index"`
    CreatedAt   time.Time       `gorm:"not null"`
}

func (j jobModel) TableName() string { return d.table }

// failedJobModel is the GORM model for the failed_jobs table.
type failedJobModel struct {
    ID       uint64    `gorm:"primaryKey;autoIncrement"`
    Queue    string    `gorm:"size:255;not null"`
    Type     string    `gorm:"size:255;not null"`
    Payload  string    `gorm:"type:text;not null"`
    Error    string    `gorm:"type:text;not null"`
    FailedAt time.Time `gorm:"not null"`
}

func (f failedJobModel) TableName() string { return d.failedTable }
```

**Note**: `uint64` for IDs is deliberate — jobs accumulate faster than domain models. The GORM models convert to/from the `Job` struct internally.

**Pop implementation**: Uses `SELECT ... WHERE queue = ? AND available_at <= NOW() AND reserved_at IS NULL ORDER BY id ASC LIMIT 1 FOR UPDATE SKIP LOCKED` to atomically reserve a job. Falls back to non-SKIP LOCKED for SQLite (single-worker only).

**Fail implementation**: Inserts into `failed_jobs` table, deletes from `jobs` table.

### Redis Driver

**Package**: `core/queue`
**File**: `redis.go`

```go
// RedisDriver implements Driver using Redis lists.
type RedisDriver struct {
    client *redis.Client
    prefix string // key prefix, default "rapidgo:queue:"
}

func NewRedisDriver(client *redis.Client) *RedisDriver
```

**Storage**:
- Pending jobs: `rapidgo:queue:{queueName}` — Redis List (LPUSH/BRPOP)
- Job data: `rapidgo:queue:job:{id}` — Redis Hash with TTL
- Failed jobs: `rapidgo:queue:failed:{queueName}` — Redis List

**Pop implementation**: `BRPOP` with 2-second timeout (non-blocking poll loop). Job data stored as JSON in the list entry itself (no separate hash needed — simpler).

### Memory Driver

**Package**: `core/queue`
**File**: `memory.go`

```go
// MemoryDriver implements Driver using in-process channels.
type MemoryDriver struct {
    mu     sync.Mutex
    queues map[string][]*Job
    failed []*Job
    nextID uint64
}

func NewMemoryDriver() *MemoryDriver
```

In-process only. No persistence. Useful for testing.

### Sync Driver

**Package**: `core/queue`
**File**: `sync.go`

```go
// SyncDriver executes jobs immediately when dispatched.
// No worker needed. Useful for local development and testing.
type SyncDriver struct{}

func NewSyncDriver() *SyncDriver
```

`Push` immediately resolves the handler and executes it. Returns the handler's error directly. No retries.

### Worker Pool

**Package**: `core/queue`
**File**: `worker.go`

```go
// WorkerConfig configures the worker pool.
type WorkerConfig struct {
    Queues      []string      // Queue names to consume (default: ["default"])
    Concurrency int           // Number of worker goroutines (default: 1)
    PollInterval time.Duration // How often to check for jobs (default: 3s)
    MaxAttempts uint          // Default max attempts for retried jobs (default: 3)
    RetryDelay  time.Duration // Delay before retry (default: 30s)
    Timeout     time.Duration // Max job processing time (default: 60s)
}

// Worker manages a pool of goroutines that process jobs.
type Worker struct {
    driver Driver
    config WorkerConfig
}

// NewWorker creates a worker with the given driver and config.
func NewWorker(driver Driver, config WorkerConfig) *Worker

// Run starts the worker pool and blocks until ctx is cancelled.
// Returns nil on clean shutdown.
func (w *Worker) Run(ctx context.Context) error
```

**Run flow**:
1. Start `config.Concurrency` goroutines
2. Each goroutine loops: Pop → process → Delete/Release/Fail
3. If Pop returns nil (no job), sleep `PollInterval`
4. On `ctx.Done()`, finish current job, then exit
5. Wait for all goroutines to finish (sync.WaitGroup)

**Job processing per goroutine**:
1. `Pop(ctx, queue)` — get next job
2. Resolve handler by `job.Type`
3. Create job context with `Timeout` deadline
4. Call handler with panic recovery
5. On success: `Delete(job)`
6. On failure: if `attempts < maxAttempts` → `Release(job, retryDelay)`, else → `Fail(job, err)`

---

## Data Flow

### Dispatch (application code → queue)

```
Controller/Service
    → queue.Dispatch(ctx, "emails", "send_welcome", payload)
        → JSON marshal payload
        → driver.Push(ctx, &Job{Queue: "emails", Type: "send_welcome", ...})
            → INSERT INTO jobs ... (database)
            → LPUSH rapidgo:queue:emails ... (redis)
```

### Processing (worker → job handler)

```
`rapidgo work --queues=emails,default --workers=4`
    → Bootstrap app (NewApp with QueueProvider)
    → Resolve Dispatcher from container
    → Create Worker(driver, config)
    → Worker.Run(ctx):
        → goroutine per worker:
            → driver.Pop(ctx, "emails") → *Job
            → ResolveHandler("send_welcome") → HandlerFunc
            → handler(ctx, job.Payload)
            → success: driver.Delete(job)
            → failure: driver.Release(job, delay) or driver.Fail(job, err)
        → ctx.Done() → finish current job → exit goroutine
    → all goroutines done → clean exit
```

---

## Configuration

| Key | Type | Default | Description |
|---|---|---|---|
| `QUEUE_DRIVER` | string | `"database"` | Backend: `database`, `redis`, `memory`, `sync` |
| `QUEUE_TABLE` | string | `"jobs"` | Database table name for jobs |
| `QUEUE_FAILED_TABLE` | string | `"failed_jobs"` | Database table name for failed jobs |
| `QUEUE_DEFAULT` | string | `"default"` | Default queue name when none specified |
| `QUEUE_MAX_ATTEMPTS` | int | `3` | Default max retry attempts |
| `QUEUE_RETRY_DELAY` | int | `30` | Seconds to wait before retrying a failed job |
| `QUEUE_TIMEOUT` | int | `60` | Max seconds a job can run before timeout |
| `QUEUE_POLL_INTERVAL` | int | `3` | Seconds between polling for new jobs (database driver) |
| `REDIS_HOST` | string | `"localhost"` | Redis host (shared with cache) |
| `REDIS_PORT` | string | `"6379"` | Redis port (shared with cache) |
| `REDIS_PASSWORD` | string | `""` | Redis password (shared with cache) |

**Note**: Redis env vars are shared with the existing cache system — no new Redis-specific vars needed.

---

## CLI Command

```
rapidgo work [flags]

Flags:
  --queues, -q    string   Comma-separated queue names (default: "default")
  --workers, -w   int      Number of concurrent workers (default: 1)
  --timeout       int      Max job processing time in seconds (default: from QUEUE_TIMEOUT)
```

**Bootstrap flow**:
1. `config.Load()` — load .env
2. Create `app.New()` and register only the providers needed for queue work:
   - `ConfigProvider`, `LoggerProvider`, `DatabaseProvider`, `RedisProvider`, `QueueProvider`
   - No `SessionProvider`, `MiddlewareProvider`, or `RouterProvider` (unnecessary for workers)
3. `Boot()` the app
4. Resolve `*Dispatcher` from container as `"queue"`
5. Register job handlers (application-defined)
6. Create `Worker` with config from flags + env
7. `worker.Run(ctx)` — blocks until signal

**Note**: The `work` command does NOT use `NewApp(mode)` because workers don't need HTTP providers. It builds a minimal app with only config, logging, database, Redis, and queue.

---

## Provider Design

### RedisProvider

**File**: `app/providers/redis_provider.go`

```go
type RedisProvider struct{}

func (p *RedisProvider) Register(c *container.Container) {
    c.Singleton("redis", func(c *container.Container) interface{} {
        return redis.NewClient(&redis.Options{
            Addr:     config.Env("REDIS_HOST", "localhost") + ":" + config.Env("REDIS_PORT", "6379"),
            Password: config.Env("REDIS_PASSWORD", ""),
        })
    })
}

func (p *RedisProvider) Boot(c *container.Container) {}
```

### QueueProvider

**File**: `app/providers/queue_provider.go`

```go
type QueueProvider struct{}

func (p *QueueProvider) Register(c *container.Container) {
    c.Singleton("queue", func(c *container.Container) interface{} {
        driver := config.Env("QUEUE_DRIVER", "database")
        switch driver {
        case "database":
            db := container.MustMake[*gorm.DB](c, "db")
            table := config.Env("QUEUE_TABLE", "jobs")
            failedTable := config.Env("QUEUE_FAILED_TABLE", "failed_jobs")
            return queue.NewDispatcher(queue.NewDatabaseDriver(db, table, failedTable))
        case "redis":
            client := container.MustMake[*redis.Client](c, "redis")
            return queue.NewDispatcher(queue.NewRedisDriver(client))
        case "memory":
            return queue.NewDispatcher(queue.NewMemoryDriver())
        case "sync":
            return queue.NewDispatcher(queue.NewSyncDriver())
        default:
            panic("queue: unsupported driver: " + driver)
        }
    })
}

func (p *QueueProvider) Boot(c *container.Container) {}
```

### Registration in `root.go`

`NewApp(mode)` is NOT modified. The `work` command builds its own minimal app to avoid registering unnecessary HTTP providers. However, `RedisProvider` and `QueueProvider` are also registered in `NewApp` for HTTP modes that need to dispatch jobs:

```go
// NewApp(mode) — add RedisProvider + QueueProvider for dispatch from HTTP handlers
a.Register(&ConfigProvider{})
a.Register(&LoggerProvider{})
if mode.Has(service.ModeWeb) || mode.Has(service.ModeAPI) || mode.Has(service.ModeWS) {
    a.Register(&DatabaseProvider{})
}
a.Register(&RedisProvider{})     // NEW — shared Redis client
a.Register(&QueueProvider{})     // NEW — queue dispatcher (dispatch-only from HTTP)
if mode.Has(service.ModeWeb) {
    a.Register(&SessionProvider{})
}
a.Register(&MiddlewareProvider{Mode: mode})
a.Register(&RouterProvider{Mode: mode})
```

**Worker-specific bootstrap** (in `core/cli/work.go`):

```go
// Minimal app for worker — no HTTP providers
app := app.New()
app.Register(&providers.ConfigProvider{})
app.Register(&providers.LoggerProvider{})
app.Register(&providers.DatabaseProvider{})
app.Register(&providers.RedisProvider{})
app.Register(&providers.QueueProvider{})
app.Boot()
```

---

## Security Considerations

- **Job payloads are plain JSON** — do not store secrets (passwords, API keys) in job payloads. Pass references (user IDs, resource IDs) and look up sensitive data in the handler.
- **SQL injection** — GORM parameterized queries prevent injection in the database driver.
- **Panic recovery** — worker goroutines recover from panics to prevent one bad job from crashing the pool.
- **Timeout enforcement** — per-job context deadline prevents hung jobs from blocking workers.

---

## Trade-offs & Alternatives

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| Custom queue implementation | Framework-native, consistent patterns, no external deps, full control | More code to maintain, not battle-tested at scale | ✅ Selected |
| External library (e.g., riverqueue) | Battle-tested, feature-rich, active maintenance | External dependency, different API patterns, heavier | ❌ Adds dependency against framework philosophy |
| AMQP / RabbitMQ driver | Industry standard, message routing, pub/sub | Requires RabbitMQ infrastructure, complex setup | ❌ Too heavy for default, could be a plugin later |
| goroutine-only (no persistence) | Simple, no backend needed | Jobs lost on crash, no retries, no visibility | ❌ Not suitable for production work |

---

## Next

Create tasks doc → `42-queue-workers-tasks.md`
