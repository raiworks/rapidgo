# 📋 Tasks: Task Scheduler / Cron

> **Feature**: `43` — Task Scheduler / Cron
> **Architecture**: [`43-task-scheduler-architecture.md`](43-task-scheduler-architecture.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-07

---

## Phase A — Core Scheduler (`core/scheduler/scheduler.go`)

| # | Task | Detail |
|---|------|--------|
| A1 | Add `robfig/cron/v3` dependency | `go get github.com/robfig/cron/v3` |
| A2 | Create `core/scheduler/scheduler.go` | `TaskFunc`, `Task`, `Scheduler` struct |
| A3 | Implement `New()` | Create cron engine with 5-field + descriptor parser, store logger |
| A4 | Implement `Add()` | Validate cron expression, wrap handler, store in tasks slice |
| A5 | Implement `wrap()` | Panic recovery, structured logging (start/end/fail/panic), duration |
| A6 | Implement `Tasks()` | Return slice of registered tasks for introspection |
| A7 | Implement `Run(ctx)` | Start cron, block on ctx.Done(), stop cron (wait), return nil |

**Exit**: `Scheduler` compiles, ready for tests.

---

## Phase B — Tests (`core/scheduler/scheduler_test.go`)

| # | Task | Detail |
|---|------|--------|
| B1 | `TestNewScheduler` | New returns non-nil Scheduler with no tasks |
| B2 | `TestAddTask` | Add a task, verify Tasks() returns it with correct name and schedule |
| B3 | `TestAddMultipleTasks` | Add 3 tasks, verify Tasks() returns all 3 |
| B4 | `TestAddInvalidCron` | Add with invalid cron expression returns error |
| B5 | `TestTaskExecution` | Add `@every 1s` task, Run for 2s, verify task ran at least once |
| B6 | `TestGracefulShutdown` | Start Run, cancel ctx, verify Run returns nil without hanging |
| B7 | `TestPanicRecovery` | Task panics, scheduler continues running (next task still fires) |
| B8 | `TestTaskError` | Task returns error, scheduler continues (error logged, not fatal) |
| B9 | `TestTaskDuration` | Verify structured log includes duration field |
| B10 | `TestDescriptorParser` | Add with `@daily`, `@every 5m` — no error returned |

**Exit**: All tests pass.

---

## Phase C — CLI Command (`core/cli/schedule_run.go`)

| # | Task | Detail |
|---|------|--------|
| C1 | Create `core/cli/schedule_run.go` | `scheduleRunCmd` Cobra command |
| C2 | Minimal bootstrap | Config, Logger, DB, Redis, Queue — no HTTP providers |
| C3 | Create scheduler, call `RegisterSchedule()` | Hook into user-land |
| C4 | Print banner with task table | Show registered tasks on startup |
| C5 | Signal handling + graceful shutdown | SIGINT/SIGTERM via `signal.NotifyContext` |

**Exit**: `rapidgo schedule:run` starts and stops cleanly.

---

## Phase D — Integration & Wiring

| # | Task | Detail |
|---|------|--------|
| D1 | Create `app/schedule/schedule.go` | `RegisterSchedule()` with example heartbeat task |
| D2 | Update `core/cli/root.go` | Add `scheduleRunCmd` to init() |

**Exit**: Full pipeline works end-to-end.

---

## Phase E — Verification

| # | Task | Detail |
|---|------|--------|
| E1 | Run scheduler tests | `go test ./core/scheduler/...` — all pass |
| E2 | Run full test suite | `go test ./...` — all 32+ packages pass |
| E3 | Manual smoke test (optional) | `rapidgo schedule:run` starts, shows banner, ctrl-C stops cleanly |

**Exit**: Feature is complete and tested.

---

## Summary

| Phase | Files | Tasks |
|-------|-------|-------|
| A — Core | 1 new | 7 |
| B — Tests | 1 new | 10 |
| C — CLI | 1 new | 5 |
| D — Wiring | 1 new, 1 mod | 2 |
| E — Verify | — | 3 |
| **Total** | **4 new, 1 mod** | **27** |

---

## Next

Test plan → `43-task-scheduler-testplan.md`
