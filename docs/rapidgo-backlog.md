# 🔧 RapidGo Framework Backlog

> **Framework**: RapidGo v2
> **Repo**: https://github.com/RAiWorks/RapidGo
> **Last Updated**: 2026-03-11
> **Purpose**: Track features, fixes, and improvements needed in the RapidGo framework — discovered while building SpecialMedia. Work on these in the framework repo separately.

---

## How This Doc Works

As we build SpecialMedia, we'll discover things the framework doesn't support yet, bugs, or improvements that belong in the framework — not in the app. Log them here. Each item gets a priority and status.

### Priority Levels

| Priority | Meaning |
|---|---|
| 🔴 Blocker | SpecialMedia cannot proceed without this. Fix immediately in the framework. |
| 🟠 High | Needed for current or next feature. Should be done soon. |
| 🟡 Medium | Would improve DX or performance. Can wait until current feature ships. |
| 🟢 Low | Nice to have. Do when there's time. |

### Status

| Status | Meaning |
|---|---|
| ⬜ Open | Identified, not started |
| 🟡 In Progress | Being worked on in the framework repo |
| ✅ Done | Fixed/added in framework, available in latest version |
| ❌ Won't Do | Decided against — documented why |

---

## Missing Features (Framework — v2.1.0)

Features that belong in the framework. Targeted for v2.1.0 release.

| # | Feature | Priority | Status | Needed By (Feature #) | Notes |
|---|---|---|---|---|---|
| F-002 | Cursor-based pagination helper | 🟠 High | ✅ Done | #08 Home Feed | Added `database/pagination.go` with `Paginate()` (offset-based) and `CursorPaginate()` (cursor-based, base64 cursors, next/prev direction). 10 tests. |
| F-004 | Notification system abstraction | 🟠 High | ✅ Done | #11 Notifications | Added `core/notification/` package: Notification, Channel, Notifiable interfaces. Built-in DatabaseChannel (GORM) and MailChannel. Notifier dispatches to multiple channels. 6 tests. |

## Project-Level Features (SpecialMedia)

Features that are app-specific, not framework concerns. Handle in the SpecialMedia repo.

| # | Feature | Priority | Reason Not Framework |
|---|---|---|---|
| F-001 | Image processing (resize, crop, format conversion) | 🟠 High | Too app-specific. Use `disintegration/imaging` directly in SpecialMedia. Framework stores files, doesn't process them. |
| F-003 | Polymorphic relationship helpers for GORM | 🟡 Medium | "Likeable", "Commentable" are business domain patterns. GORM already supports polymorphic associations natively. |
| F-005 | Soft delete scopes on relationships | 🟢 Low | GORM API usage (`Unscoped()`, preload callbacks). Document the pattern, don't wrap it. |
| F-006 | Admin panel scaffolding improvements | 🟢 Low | Filters, search, bulk actions are app-specific UI. `make:admin` generates the starting point; customization is project-level. |

---

## Bugs / Fixes

Issues discovered while building SpecialMedia.

| # | Bug | Priority | Status | Found In (Feature #) | Description |
|---|---|---|---|---|---|
| B-001 | — | — | — | — | No bugs found yet. Will be logged as discovered. |

---

## Improvements / Refactors (Framework — v2.1.0)

Things that work but could be better. Targeted for v2.1.0 release.

| # | Improvement | Priority | Status | Context | Notes |
|---|---|---|---|---|---|
| I-001 | WebSocket room cleanup on disconnect | 🟡 Medium | ✅ Done | #10 Messaging | Added `HubConfig` (PingInterval/PongTimeout), `NewHubWithConfig()`, heartbeat goroutine with Ping/Pong, `OnJoin`/`OnLeave` callbacks. Dead connections auto-removed. 4 new tests. |
| I-002 | Queue job retry with backoff | 🟡 Medium | ✅ Done | #04 Media Upload | Added `BackoffSeconds []uint` to `Job`, `DispatchWithBackoff()` on `Dispatcher`, per-job `retryDelay()` in worker with fallback to config. 4 new tests. |
| I-003 | Rate limiter per-user (not just per-IP) | 🟡 Medium | ✅ Done | #02 Auth | Added `RateLimitConfig` struct with `Rate` and `KeyFunc`, plus `RateLimitWithConfig()` middleware. Supports per-user or any custom key. |
| I-005 | GORM query logging in development | 🟡 Medium | ✅ Done | General | Added `newGormLogger()` — auto-enables in dev mode, configurable via `DB_LOG` and `DB_SLOW_THRESHOLD_MS` env vars. Silent in production. |

## Project-Level Improvements (SpecialMedia)

| # | Improvement | Priority | Reason Not Framework |
|---|---|---|---|
| I-004 | Config hot-reload for branding | 🟢 Low | `branding.json` is a SpecialMedia concept. Implement a file watcher in the app if needed. |

---

## Documentation Gaps

Framework docs that are missing or incomplete.

| # | Doc Needed | Priority | Status | Notes |
|---|---|---|---|---|
| D-001 | WebSocket rooms/channels usage guide | 🟠 High | ✅ Done | Created `docs/framework/guides/websocket-rooms.md` — Hub, rooms, broadcast, heartbeat, callbacks, full chat example. |
| D-002 | Queue job creation and registration guide | 🟡 Medium | ✅ Done | Created `docs/framework/guides/queue-jobs.md` — handlers, dispatch, backoff, drivers, workers, full example. |
| D-003 | OAuth2 provider setup guide (Google, GitHub, Facebook) | 🟡 Medium | ✅ Done | Created `docs/framework/guides/oauth-setup.md` — credentials, redirect/callback routes, custom providers, security. |
| D-004 | Multi-port serving (service mode) guide | 🟢 Low | ✅ Done | Created `docs/framework/guides/service-mode.md` — modes, multi-port, Docker Compose, Nginx example. |

---

## How to Work on These

1. Pick an item from this doc
2. Open the RapidGo framework repo (`github.com/RAiWorks/RapidGo`)
3. Create a branch: `feature/F-001-image-processing` or `fix/B-001-description`
4. Implement, test, merge to framework main
5. Update this doc: mark status as ✅ Done, note the version/commit
6. In SpecialMedia: `go get github.com/RAiWorks/RapidGo/v2@latest` to pull the update
7. Continue building the feature that needed it

---

## Change Log

| Date | Change |
|---|---|
| 2026-03-11 | Initial backlog created with 6 features, 5 improvements, 4 doc gaps identified from architecture analysis |
| 2026-03-11 | Classified items: 7 framework (v2.1.0), 5 project-level (SpecialMedia). Created Feature #59 for v2.1.0 release. |
| 2026-03-11 | Built I-005 (GORM logging), I-001 (WebSocket heartbeat + callbacks), I-002 (Queue backoff), I-003 (Rate limiter per-user). All compiled, vetted, tested. |
| 2026-03-11 | Built F-002 (Cursor pagination, 10 tests), F-004 (Notification system, 6 tests). Created D-001–D-004 guides (websocket-rooms, queue-jobs, oauth-setup, service-mode). All v2.1.0 items complete. |

---

*This document is updated continuously as SpecialMedia development reveals framework needs.*
