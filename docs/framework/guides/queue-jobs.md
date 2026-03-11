---
title: "Queue Jobs Guide"
version: "0.1.0"
status: "Final"
date: "2026-03-11"
last_updated: "2026-03-11"
authors:
  - "RAiWorks"
supersedes: ""
---

# Queue Jobs Guide

## Abstract

How to create, dispatch, and process background jobs using the RapidGo
queue system — including drivers, retry strategies, backoff, and workers.

## Table of Contents

1. [Overview](#1-overview)
2. [Registering Job Handlers](#2-registering-job-handlers)
3. [Dispatching Jobs](#3-dispatching-jobs)
4. [Retry and Backoff](#4-retry-and-backoff)
5. [Delayed Jobs](#5-delayed-jobs)
6. [Queue Drivers](#6-queue-drivers)
7. [Running Workers](#7-running-workers)
8. [Worker Configuration](#8-worker-configuration)
9. [Full Example: Email Queue](#9-full-example-email-queue)
10. [References](#10-references)

---

## 1. Overview

The `core/queue` package provides:

- **Dispatcher** — pushes jobs onto queues
- **Worker** — polls queues and processes jobs with retry logic
- **Drivers** — Memory, Database (GORM), Redis, Sync (immediate)
- **Handler registry** — maps job type names to handler functions

## 2. Registering Job Handlers

Register a handler for each job type at application startup:

```go
import "github.com/RAiWorks/RapidGo/v2/core/queue"

queue.RegisterHandler("send_email", func(ctx context.Context, payload json.RawMessage) error {
    var data struct {
        To      string `json:"to"`
        Subject string `json:"subject"`
        Body    string `json:"body"`
    }
    if err := json.Unmarshal(payload, &data); err != nil {
        return err
    }
    return mailer.Send(data.To, data.Subject, data.Body)
})
```

Handlers receive a `context.Context` (with timeout) and the raw JSON
payload. Return `nil` on success or an `error` to trigger retry.

## 3. Dispatching Jobs

```go
dispatcher := queue.NewDispatcher(driver)

// Basic dispatch (MaxAttempts=3, available immediately)
err := dispatcher.Dispatch(ctx, "default", "send_email", map[string]string{
    "to":      "user@example.com",
    "subject": "Welcome!",
    "body":    "<p>Hello</p>",
})
```

Parameters: `(ctx, queueName, typeName, payload)`

The payload can be any value that marshals to JSON.

## 4. Retry and Backoff

### Default Retry

By default, failed jobs are retried up to 3 times with the worker's
flat `RetryDelay` (default 30s).

### Per-Job Backoff

Use `DispatchWithBackoff` for exponential or custom backoff schedules:

```go
// 5s → 30s → 120s between retries. MaxAttempts = len(backoff) + 1 = 4.
err := dispatcher.DispatchWithBackoff(ctx, "default", "process_image",
    payload, []uint{5, 30, 120})
```

The `BackoffSeconds` slice maps to retry attempts:
- Attempt 1 fails → wait `backoff[0]` seconds (5s)
- Attempt 2 fails → wait `backoff[1]` seconds (30s)
- Attempt 3 fails → wait `backoff[2]` seconds (120s)
- Attempt 4 fails → job moves to failed storage

If a job has more retries than backoff entries, the last value is reused.

## 5. Delayed Jobs

Schedule a job to become available after a delay:

```go
err := dispatcher.DispatchDelayed(ctx, "default", "send_reminder",
    payload, 1*time.Hour)
```

## 6. Queue Drivers

### Memory (testing/development)

```go
driver := queue.NewMemoryDriver()
```

### Database (GORM)

```go
driver := queue.NewDatabaseDriver(db) // *gorm.DB
```

Requires the `jobs` and `failed_jobs` tables (auto-migrated).

### Redis

```go
driver := queue.NewRedisDriver(redisClient) // *redis.Client
```

### Sync (immediate execution, no queue)

```go
driver := queue.NewSyncDriver()
```

Executes the handler immediately in the current goroutine. Useful for
testing or simple setups where background processing isn't needed.

## 7. Running Workers

Workers are started via the CLI:

```bash
go run cmd/main.go work
```

Or programmatically:

```go
worker := queue.NewWorker(driver, queue.WorkerConfig{
    Queues:       []string{"default", "emails"},
    Concurrency:  4,
    PollInterval: 3 * time.Second,
}, logger)

ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
defer cancel()
worker.Run(ctx) // blocks until ctx is cancelled
```

## 8. Worker Configuration

| Field | Default | Description |
|---|---|---|
| `Queues` | `["default"]` | Queue names to poll, in priority order |
| `Concurrency` | `1` | Number of goroutines processing jobs |
| `PollInterval` | `3s` | How often to check for new jobs |
| `MaxAttempts` | `3` | Default max attempts (overridden by per-job setting) |
| `RetryDelay` | `30s` | Default delay between retries (overridden by per-job backoff) |
| `Timeout` | `60s` | Max execution time per job |

## 9. Full Example: Email Queue

```go
// Register handler
queue.RegisterHandler("welcome_email", func(ctx context.Context, payload json.RawMessage) error {
    var p struct{ Email string `json:"email"` }
    json.Unmarshal(payload, &p)
    return mailer.Send(p.Email, "Welcome!", "<h1>Welcome to our app!</h1>")
})

// Dispatch with retry backoff
driver := queue.NewRedisDriver(redisClient)
dispatcher := queue.NewDispatcher(driver)
dispatcher.DispatchWithBackoff(ctx, "emails", "welcome_email",
    map[string]string{"email": "new@user.com"},
    []uint{10, 60, 300}, // 10s, 1m, 5m
)

// Run worker
worker := queue.NewWorker(driver, queue.WorkerConfig{
    Queues:      []string{"emails"},
    Concurrency: 2,
    Timeout:     30 * time.Second,
}, slog.Default())
worker.Run(ctx)
```

## 10. References

- [Queue source](../../core/queue/) — implementation
- [Worker CLI](../cli/) — `work` command documentation
