# ✅ Tasks: Queue Workers / Background Jobs

> **Feature**: `42` — Queue Workers / Background Jobs
> **Architecture**: [`42-queue-workers-architecture.md`](42-queue-workers-architecture.md)
> **Branch**: `feature/42-queue-workers`
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/20 tasks complete

---

## Pre-Flight Checklist

- [ ] Discussion doc is marked COMPLETE
- [ ] Architecture doc is FINALIZED
- [ ] Feature branch created from latest `main`
- [ ] Dependent features are merged to `main`

---

## Phase A — Queue Core & Job Types

> Foundation: Job struct, Driver interface, handler registry, Dispatcher

- [ ] **A.1** — Create `core/queue/queue.go`
  - [ ] `Job` struct (ID, Queue, Type, Payload, Attempts, MaxAttempts, AvailableAt, ReservedAt, CreatedAt)
  - [ ] `HandlerFunc` type: `func(ctx context.Context, payload json.RawMessage) error`
  - [ ] Handler registry: `RegisterHandler(typeName, handler)`, `ResolveHandler(typeName)`, `ResetHandlers()`
  - [ ] `Driver` interface (Push, Pop, Delete, Release, Fail, Size)
  - [ ] `Dispatcher` struct with `NewDispatcher(driver)`, `Dispatch(ctx, queue, typeName, payload)`, `DispatchDelayed(ctx, queue, typeName, payload, delay)`, `Driver()`
  - [ ] `Dispatch` validates handler exists, marshals payload to JSON, constructs Job, calls `driver.Push()`
- [ ] 📍 **Checkpoint A** — `go build ./core/queue/...` compiles

---

## Phase B — Drivers

> Implement all four backends

- [ ] **B.1** — Create `core/queue/memory.go`
  - [ ] `MemoryDriver` struct with mutex-protected job slices per queue
  - [ ] `Push` — append to queue slice, assign incremented ID
  - [ ] `Pop` — find first job where `AvailableAt <= now` and `ReservedAt == nil`, set ReservedAt
  - [ ] `Delete` — remove from slice
  - [ ] `Release` — clear ReservedAt, set new AvailableAt, increment Attempts
  - [ ] `Fail` — move to failed slice, remove from queue slice
  - [ ] `Size` — count pending jobs in queue
- [ ] **B.2** — Create `core/queue/sync.go`
  - [ ] `SyncDriver` struct (stateless)
  - [ ] `Push` — immediately resolve handler, execute with context, return handler error
  - [ ] `Pop` — return nil, nil (no-op, jobs never queued)
  - [ ] `Delete`, `Release`, `Fail` — no-op
  - [ ] `Size` — always 0
- [ ] **B.3** — Create `core/queue/database.go`
  - [ ] `DatabaseDriver` struct with `*gorm.DB`
  - [ ] GORM models: `jobModel`, `failedJobModel` (mapped to table names from config)
  - [ ] `Push` — INSERT into jobs table
  - [ ] `Pop` — transaction with `FOR UPDATE SKIP LOCKED` (fallback for SQLite)
  - [ ] `Delete` — DELETE from jobs table
  - [ ] `Release` — UPDATE reserved_at=NULL, available_at=now+delay, attempts++
  - [ ] `Fail` — INSERT into failed_jobs, DELETE from jobs
  - [ ] `Size` — COUNT WHERE queue=? AND reserved_at IS NULL
- [ ] **B.4** — Create `core/queue/redis.go`
  - [ ] `RedisDriver` struct with `*redis.Client` and prefix
  - [ ] `Push` — JSON marshal job, LPUSH to `{prefix}{queue}`
  - [ ] `Pop` — RPOP from `{prefix}{queue}`, JSON unmarshal
  - [ ] `Delete` — no-op (already removed by RPOP)
  - [ ] `Release` — LPUSH back to queue with updated AvailableAt and Attempts
  - [ ] `Fail` — LPUSH to `{prefix}failed:{queue}`
  - [ ] `Size` — LLEN `{prefix}{queue}`
- [ ] 📍 **Checkpoint B** — `go build ./core/queue/...` compiles, all driver types implement `Driver` interface

---

## Phase C — Worker Pool

> Worker goroutines that process jobs with retry logic and graceful shutdown

- [ ] **C.1** — Create `core/queue/worker.go`
  - [ ] `WorkerConfig` struct (Queues, Concurrency, PollInterval, MaxAttempts, RetryDelay, Timeout)
  - [ ] `Worker` struct with `NewWorker(driver, config)`
  - [ ] `Run(ctx context.Context) error` — start goroutine pool, block until ctx done
  - [ ] Per-goroutine loop: Pop → resolve handler → execute with timeout → Delete/Release/Fail
  - [ ] Panic recovery in job execution
  - [ ] Round-robin across configured queues
  - [ ] Sleep `PollInterval` when all queues empty
  - [ ] `sync.WaitGroup` for clean goroutine shutdown
- [ ] 📍 **Checkpoint C** — Worker can process in-memory jobs in a test

---

## Phase D — Providers & CLI

> Integration with framework: providers, CLI command, migration

- [ ] **D.1** — Create `app/providers/redis_provider.go`
  - [ ] `RedisProvider` struct implementing `container.Provider`
  - [ ] `Register` — Singleton `"redis"` with `*redis.Client` from REDIS_HOST/PORT/PASSWORD
  - [ ] `Boot` — no-op
- [ ] **D.2** — Create `app/providers/queue_provider.go`
  - [ ] `QueueProvider` struct implementing `container.Provider`
  - [ ] `Register` — Singleton `"queue"` that reads QUEUE_DRIVER and returns `*Dispatcher`
  - [ ] Switch on driver: database→DatabaseDriver, redis→RedisDriver, memory→MemoryDriver, sync→SyncDriver
  - [ ] `Boot` — no-op
- [ ] **D.3** — Update `core/cli/root.go`
  - [ ] Add `RedisProvider{}` and `QueueProvider{}` registration in `NewApp()` (after DatabaseProvider, unconditional — lazy singletons don't connect until used)
  - [ ] Add `rootCmd.AddCommand(workCmd)` in `init()`
- [ ] **D.4** — Create `core/cli/work.go`
  - [ ] `workCmd` Cobra command: `rapidgo work`
  - [ ] Flags: `--queues/-q` (default "default"), `--workers/-w` (default 1), `--timeout` (default from env)
  - [ ] Minimal bootstrap: `config.Load()`, build app with only ConfigProvider, LoggerProvider, DatabaseProvider, RedisProvider, QueueProvider (no HTTP providers)
  - [ ] Register application job handlers
  - [ ] Create `Worker` with config from flags + env
  - [ ] `worker.Run(ctx)` with signal handling (SIGINT/SIGTERM)
  - [ ] Banner: show queue names, worker count, driver
  - [ ] Register command in `init()`
- [ ] **D.5** — Create jobs migration
  - [ ] Create migration file with `jobs` and `failed_jobs` table schemas
  - [ ] Up: CREATE TABLE jobs, CREATE TABLE failed_jobs
  - [ ] Down: DROP TABLE failed_jobs, DROP TABLE jobs
- [ ] **D.6** — Update `.env`
  - [ ] Add QUEUE_DRIVER, QUEUE_DEFAULT, QUEUE_MAX_ATTEMPTS, QUEUE_RETRY_DELAY, QUEUE_TIMEOUT, QUEUE_POLL_INTERVAL
- [ ] **D.7** — Create `app/jobs/example_job.go`
  - [ ] Example job handler showing the pattern
  - [ ] `RegisterJobs()` function that registers all app job handlers
- [ ] 📍 **Checkpoint D** — `go build ./...` compiles, `rapidgo work --help` shows flags, migration runs

---

## Phase E — Tests

> Comprehensive test coverage for all components

- [ ] **E.1** — Create `core/queue/queue_test.go`
  - [ ] Test handler registry: register, resolve, resolve unknown, reset
  - [ ] Test Dispatcher.Dispatch: marshals payload, calls driver.Push
  - [ ] Test Dispatcher.DispatchDelayed: sets future AvailableAt
  - [ ] Test MemoryDriver: full Push→Pop→Delete cycle
  - [ ] Test MemoryDriver: Release increments attempts, sets future AvailableAt
  - [ ] Test MemoryDriver: Fail moves to failed storage
  - [ ] Test MemoryDriver: Pop respects AvailableAt (delayed jobs)
  - [ ] Test MemoryDriver: Size counts correctly
  - [ ] Test SyncDriver: Push executes handler immediately
  - [ ] Test SyncDriver: Push returns handler error
  - [ ] Test SyncDriver: Push with unknown handler returns error
  - [ ] Test DatabaseDriver: Push/Pop/Delete cycle (SQLite)
  - [ ] Test DatabaseDriver: Release and retry
  - [ ] Test DatabaseDriver: Fail moves to failed_jobs
  - [ ] Test RedisDriver: Push/Pop/Delete cycle (miniredis)
  - [ ] Test RedisDriver: Release pushes back
  - [ ] Test RedisDriver: Fail stores in failed list
  - [ ] Test RedisDriver: Size returns list length
  - [ ] Test Worker: processes jobs from memory driver
  - [ ] Test Worker: retries on failure up to max attempts
  - [ ] Test Worker: moves to failed after max attempts
  - [ ] Test Worker: recovers from handler panic
  - [ ] Test Worker: stops on context cancellation (graceful shutdown)
  - [ ] Test Worker: respects job timeout
- [ ] 📍 **Checkpoint E** — `go test ./core/queue/... -count=1` all pass, `go test ./... -count=1` all pass

---

## Ship 🚀

- [ ] All phases complete
- [ ] All checkpoints verified
- [ ] `go test ./... -count=1` — all packages pass
- [ ] `rapidgo work --help` shows correct flags and usage
- [ ] Final commit with descriptive message
- [ ] Push to feature branch
- [ ] Merge to `main`
- [ ] Push `main`
- [ ] **Keep the feature branch**
- [ ] Create review doc → `42-queue-workers-review.md`
