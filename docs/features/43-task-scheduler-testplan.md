# 🧪 Test Plan: Task Scheduler / Cron

> **Feature**: `43` — Task Scheduler / Cron
> **Tasks**: [`43-task-scheduler-tasks.md`](43-task-scheduler-tasks.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-07

---

## Test File

`core/scheduler/scheduler_test.go`

---

## Unit Tests

### 1. Scheduler Construction

| # | Test | Expectation |
|---|------|-------------|
| T01 | `TestNewScheduler` | `New(logger)` returns non-nil `*Scheduler`, `Tasks()` returns empty slice |
| T02 | `TestNewSchedulerNilLogger` | `New(nil)` uses `slog.Default()` — no panic |

### 2. Task Registration

| # | Test | Expectation |
|---|------|-------------|
| T03 | `TestAddTask` | `Add("*/5 * * * *", "cleanup", fn)` returns nil error, `Tasks()` has len 1 with correct Name and Schedule |
| T04 | `TestAddMultipleTasks` | Add 3 tasks, `Tasks()` returns 3 entries in order |
| T05 | `TestAddInvalidCronExpression` | `Add("not-a-cron", "bad", fn)` returns non-nil error |
| T06 | `TestAddEmptyName` | `Add("@every 1m", "", fn)` — accepted (name is optional for flexibility) |
| T07 | `TestAddDescriptors` | `Add("@every 5m", ...)` and `Add("@daily", ...)` both return nil error |

### 3. Task Execution

| # | Test | Expectation |
|---|------|-------------|
| T08 | `TestTaskExecutes` | Add `@every 1s` task that increments atomic counter, run for 2s, counter ≥ 1 |
| T09 | `TestMultipleTasksExecute` | Add 2 tasks with `@every 1s`, run for 2s, both counters ≥ 1 |
| T10 | `TestTaskReceivesContext` | Task receives non-nil `context.Context` |

### 4. Error Handling

| # | Test | Expectation |
|---|------|-------------|
| T11 | `TestTaskErrorDoesNotStopScheduler` | Task returns error, second task still executes on next tick |
| T12 | `TestPanicRecovery` | Task panics, scheduler continues running, other tasks still fire |

### 5. Graceful Shutdown

| # | Test | Expectation |
|---|------|-------------|
| T13 | `TestRunReturnsOnContextCancel` | Cancel ctx, `Run()` returns nil within 2s |
| T14 | `TestRunWaitsForRunningTask` | Task sleeps 500ms, cancel ctx during task execution, `Run()` waits for task completion then returns |

### 6. Logging

| # | Test | Expectation |
|---|------|-------------|
| T15 | `TestLogsTaskStarted` | Captured log output contains "task started" with task name |
| T16 | `TestLogsTaskCompleted` | Captured log output contains "task completed" with duration |
| T17 | `TestLogsTaskFailed` | Error-returning task's log contains "task failed" with error message |
| T18 | `TestLogsTaskPanicked` | Panicking task's log contains "task panicked" with panic value |

---

## Test Utilities

### Log Capture

Use `slog.New(slog.NewTextHandler(&buf, nil))` with a `bytes.Buffer` to capture structured log output for assertion.

### Timing

- Use short intervals (`@every 1s`) and small timeouts (2-3s)
- Use `sync/atomic` for counters to avoid data races
- Use `time.After` as a safety timeout to prevent hanging tests

---

## Coverage Matrix

| Component | Tests | Coverage |
|-----------|-------|----------|
| `New()` | T01–T02 | Construction, nil logger fallback |
| `Add()` | T03–T07 | Valid/invalid cron, multiple tasks, descriptors |
| `wrap()` | T11–T12, T15–T18 | Error handling, panic recovery, logging |
| `Run()` | T08–T10, T13–T14 | Execution, context, shutdown |
| `Tasks()` | T03–T04 | Introspection |

**Total: 18 tests**

---

## Next

Changelog → `43-task-scheduler-changelog.md`
