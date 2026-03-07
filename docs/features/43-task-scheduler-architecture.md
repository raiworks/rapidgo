# 🏗️ Architecture: Task Scheduler / Cron

> **Feature**: `43` — Task Scheduler / Cron
> **Discussion**: [`43-task-scheduler-discussion.md`](43-task-scheduler-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-07

---

## Overview

A framework-native task scheduler built on `robfig/cron/v3`. Developers define named tasks with cron expressions in `app/schedule/schedule.go`. The scheduler runs as a separate long-running process via `rapidgo schedule:run`, with graceful shutdown, per-task panic recovery, and structured logging. Tasks can optionally dispatch work to the queue system (#42).

---

## File Structure

```
core/
├── scheduler/
│   ├── scheduler.go           # Scheduler struct, Task, Add(), Run(), wraps robfig/cron
│   └── scheduler_test.go      # Tests for scheduler
├── cli/
│   ├── schedule_run.go        # NEW — `rapidgo schedule:run` command
│   └── root.go                # MODIFY — register scheduleRunCmd

app/
├── schedule/
│   └── schedule.go            # NEW — RegisterSchedule() with example tasks

.env                           # MODIFY — no new vars needed (scheduled tasks use existing config)
```

**Total**: 3 new files, 1 modified file

---

## Component Design

### TaskFunc

**Package**: `core/scheduler`
**File**: `scheduler.go`

```go
// TaskFunc is a function executed on schedule.
type TaskFunc func(ctx context.Context) error
```

### Task

**Package**: `core/scheduler`
**File**: `scheduler.go`

```go
// Task is a named scheduled task with a cron expression.
type Task struct {
    Name     string
    Schedule string   // cron expression (e.g., "*/5 * * * *")
    Run      TaskFunc
}
```

### Scheduler

**Package**: `core/scheduler`
**File**: `scheduler.go`

```go
// Scheduler wraps robfig/cron and manages named tasks.
type Scheduler struct {
    cron  *cron.Cron
    tasks []Task
    log   *slog.Logger
}

// New creates a scheduler with the standard cron parser.
func New(log *slog.Logger) *Scheduler

// Add registers a named task with a cron expression.
func (s *Scheduler) Add(schedule, name string, fn TaskFunc) error

// Tasks returns all registered tasks (for display/debugging).
func (s *Scheduler) Tasks() []Task

// Run starts the cron engine and blocks until ctx is cancelled.
// On ctx.Done(), stops the cron engine and waits for running tasks to complete.
func (s *Scheduler) Run(ctx context.Context) error
```

**Implementation details**:

`New()`:
- Creates `cron.New(cron.WithParser(cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))` for standard 5-field cron + descriptors like `@every 5m`, `@daily`
- Stores `*slog.Logger` for structured logging

`Add()`:
- Validates cron expression via `cron.AddFunc`
- Wraps the user's `TaskFunc` with:
  1. Structured log entry (task start)
  2. Panic recovery
  3. Error logging
  4. Duration measurement
- Stores `Task` in internal slice for introspection

`Run()`:
- Calls `s.cron.Start()` to begin the cron engine
- Blocks on `<-ctx.Done()`
- Calls `s.cron.Stop()` which waits for running jobs to complete
- Returns nil (clean shutdown)

### Wrapped Task Execution

Each registered task is wrapped before passing to `cron.AddFunc`:

```go
func (s *Scheduler) wrap(name string, fn TaskFunc) func() {
    return func() {
        start := time.Now()
        s.log.Info("task started", "task", name)

        defer func() {
            if r := recover(); r != nil {
                s.log.Error("task panicked", "task", name, "panic", r,
                    "stack", string(debug.Stack()))
            }
        }()

        if err := fn(context.Background()); err != nil {
            s.log.Error("task failed", "task", name, "error", err,
                "duration", time.Since(start))
            return
        }

        s.log.Info("task completed", "task", name,
            "duration", time.Since(start))
    }
}
```

---

## CLI Command

```
rapidgo schedule:run

Start the task scheduler. Runs registered tasks on their cron schedules.
Blocks until SIGINT/SIGTERM.
```

No flags needed for MVP. Configuration is code-based (tasks registered in `app/schedule/schedule.go`).

### Bootstrap Flow

```go
config.Load()

// Minimal bootstrap — no HTTP providers.
application := app.New()
application.Register(&providers.ConfigProvider{})
application.Register(&providers.LoggerProvider{})
application.Register(&providers.DatabaseProvider{})
application.Register(&providers.RedisProvider{})
application.Register(&providers.QueueProvider{})
application.Boot()

// Create scheduler.
s := scheduler.New(slog.Default())

// Register application-defined tasks.
schedule.RegisterSchedule(s, application)

// Print banner.
fmt.Println("=================================")
fmt.Println("  RapidGo Task Scheduler")
fmt.Println("=================================")
for _, t := range s.Tasks() {
    fmt.Printf("  [%s] %s\n", t.Schedule, t.Name)
}
fmt.Println("=================================")

// Graceful shutdown.
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()

return s.Run(ctx)
```

### User-Land Task Registration

**File**: `app/schedule/schedule.go`

```go
package schedule

import (
    "context"
    "log/slog"

    "github.com/RAiWorks/RapidGo/core/app"
    "github.com/RAiWorks/RapidGo/core/scheduler"
)

// RegisterSchedule defines all scheduled tasks.
func RegisterSchedule(s *scheduler.Scheduler, application *app.App) {
    s.Add("@every 1m", "heartbeat", func(ctx context.Context) error {
        slog.Info("scheduler heartbeat")
        return nil
    })
}
```

The `application` parameter gives tasks access to the container if they need to resolve services (e.g., queue dispatcher, DB).

---

## Data Flow

```
`rapidgo schedule:run`
    → config.Load()
    → Minimal app bootstrap (Config, Logger, DB, Redis, Queue)
    → scheduler.New(slog.Default())
    → schedule.RegisterSchedule(s, app) — user registers tasks
    → s.Run(ctx):
        → cron.Start() — background goroutine per cron engine
        → On schedule:
            → wrap(name, fn)():
                → log "task started"
                → defer panic recovery
                → fn(ctx) — user task executes
                → log "task completed" or "task failed"
        → <-ctx.Done() — SIGINT/SIGTERM
        → cron.Stop() — wait for running tasks
        → return nil
```

---

## Configuration

No new `.env` variables needed. The scheduler uses existing config infrastructure:
- Scheduled tasks use existing services (DB, Redis, Queue) which are already configured
- Task definitions are in code, not config files

---

## Security Considerations

- **Panic recovery**: each task is wrapped with `recover()` — one crashing task doesn't kill the scheduler
- **No external scheduling interface**: tasks are defined in compiled Go code, no user-input parsing, no injection vectors
- **Queue integration**: scheduled tasks dispatching to queue inherit the queue's security model

---

## Trade-offs & Alternatives

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| `robfig/cron/v3` wrapper | Battle-tested, cron expressions, timezone support, graceful stop | External dependency | ✅ Selected — recommended by roadmap |
| Custom `time.Ticker` | No dependency | Can't do "every Monday at 3am", needs DST handling, reinvents the wheel | ❌ Too limited |
| System crontab | OS-native, reliable | Not portable, not in version control, no structured logging | ❌ Against framework philosophy |
| Embedded in serve | Single process | Violates SRP, serve is already complex, crash coupling | ❌ Wrong separation |

---

## Future Iterations (NOT in this feature)

- Overlap prevention (`WithoutOverlapping()`)
- Distributed locking (only one instance runs a task)
- Task output capture
- Fluent API (`s.Command("cleanup").EveryFiveMinutes()`)
- Event hooks (before/after task)
- `schedule:list` command to show registered tasks

---

## Next

Tasks doc → `43-task-scheduler-tasks.md`
