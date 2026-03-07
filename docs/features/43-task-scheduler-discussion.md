# 💬 Discussion: Task Scheduler / Cron

> **Feature**: `43` — Task Scheduler / Cron
> **Status**: 🟢 COMPLETE
> **Date**: 2026-03-07

---

## What Are We Building?

A framework-native task scheduler that allows developers to define recurring tasks using cron expressions. Tasks are registered in application code and executed automatically on schedule via the `rapidgo schedule:run` CLI command. The scheduler runs as a long-running process (separate from the HTTP server), supports graceful shutdown, and integrates with the existing queue system (#42) for heavy work delegation.

---

## Why?

- **Background maintenance**: cache cleanup, report generation, data aggregation — every web application needs recurring tasks
- **Laravel parity**: `php artisan schedule:run` / `schedule:work` is one of the most-used features
- **Framework-native**: defining scheduled tasks in Go code (not crontab) keeps configuration in the codebase, version-controlled and type-safe
- **Queue integration**: scheduled tasks can dispatch to the queue for heavy work, keeping the scheduler lightweight

---

## Prior Art

| System | Approach | Notes |
|---|---|---|
| Laravel | `schedule:run` + `schedule:work` | Kernel-based task registration, cron expression DSL, overlap prevention |
| Spring | `@Scheduled` annotation | Annotation-based, embedded in app process |
| Node (node-cron) | Library-based | Simple cron wrapper |
| robfig/cron v3 | Go library | 14k+ stars, battle-tested, cron expression parser, timezone support |

---

## Constraints

1. **Separate process** — scheduler runs as `rapidgo schedule:run`, not embedded in the HTTP server (same pattern as `rapidgo work`)
2. **robfig/cron/v3** — use the industry-standard Go cron library, as recommended by the framework roadmap
3. **No new provider** — scheduler is a CLI concern, created directly in the command (like `Worker` in `work.go`)
4. **Queue-optional** — tasks CAN dispatch to queue but are NOT forced through it; simple tasks run inline
5. **Panic recovery** — one crashing task must not kill the scheduler
6. **MVP scope** — cron expressions, named tasks, logging, panic recovery. NO overlap prevention, distributed locking, or output capture in this iteration

---

## Decision Log

| # | Decision | Rationale |
|---|---|---|
| 1 | Separate CLI command, not embedded in serve | Same pattern as `work`, separation of concerns, independent process lifecycle |
| 2 | Use `robfig/cron/v3` | Roadmap recommends it, industry standard, handles timezone/DST, minimal deps |
| 3 | Named tasks with `func(ctx context.Context) error` | Consistent with queue `HandlerFunc` pattern, supports timeout/cancellation |
| 4 | User registers tasks in `app/schedule/schedule.go` | Mirrors `app/jobs/example_job.go` pattern |
| 5 | No container binding for scheduler | Short-lived CLI concern, not a shared service |
| 6 | Bootstrap includes QueueProvider | Allows scheduled tasks to dispatch queue jobs |

---

## Open Questions

_None — all resolved._

---

## Next

Architecture doc → `43-task-scheduler-architecture.md`
