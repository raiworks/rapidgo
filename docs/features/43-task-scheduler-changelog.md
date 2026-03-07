# 📝 Changelog: Task Scheduler / Cron

> **Feature**: `43` — Task Scheduler / Cron
> **Status**: � SHIPPED
> **Date**: 2026-03-07
> **Commit**: `85c1438` (merge to main)

---

## Added

- `core/scheduler/scheduler.go` — `Scheduler` struct wrapping `robfig/cron/v3` with `TaskFunc`, `Task`, `New()`, `Add()`, `Tasks()`, `Run(ctx)`, per-task panic recovery, structured logging with duration
- `core/scheduler/scheduler_test.go` — 18 tests covering construction, registration, execution, error handling, graceful shutdown, and logging
- `core/cli/schedule_run.go` — `rapidgo schedule:run` Cobra command with minimal bootstrap (Config, Logger, DB, Redis, Queue), banner, graceful shutdown via SIGINT/SIGTERM
- `app/schedule/schedule.go` — `RegisterSchedule()` with example heartbeat task
- `github.com/robfig/cron/v3` v3.0.1 dependency

## Changed

- `core/cli/root.go` — added `scheduleRunCmd` to `init()`
- `go.mod` / `go.sum` — added `robfig/cron/v3`

## Files

| File | Action |
|------|--------|
| `core/scheduler/scheduler.go` | NEW |
| `core/scheduler/scheduler_test.go` | NEW |
| `core/cli/schedule_run.go` | NEW |
| `app/schedule/schedule.go` | NEW |
| `core/cli/root.go` | MODIFIED |
| `go.mod` | MODIFIED |
| `go.sum` | MODIFIED |

## Migration Guide

No migration needed. This is a new feature with no breaking changes.
