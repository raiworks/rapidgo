# RapidGo v2 — Roadmap

> **Project**: RapidGo Framework  
> **Target**: v2.0.0  
> **Base**: v1.0.0 (tagged 2026-03-07, commit `9c0a22a`)  
> **Author**: RAi Works  
> **Date**: 2026-03-08  

---

## Goal

Transform RapidGo from a monolithic starter (clone-and-modify) into an importable Go library (`go get github.com/RAiWorks/RapidGo`) plus a companion starter template (`RapidGo-starter`).

---

## Phase Overview

| Phase | Name | Steps | Key Milestone | Breaking? |
|-------|------|:-----:|---------------|:---------:|
| **A** | Foundation | 2 | `hooks.go` exists, `AuditLog` moved | No |
| **B** | Decouple | 4 | Zero `app/` imports in `core/` | No* |
| **C** | Split | 2 | Library + Starter repos separate | Yes |
| **D** | Polish | 2 | `rapidgo new myapp` works | No |

*\*Phase B is non-breaking because `cmd/main.go` wires the hooks — the monolith keeps working.*

---

## Milestone Timeline

```
v1.0.0 (frozen) ─── main branch ──────────────────────────────────────
                     │
                     └── v2 branch created ────────────────────────────
                          │
                          ├─ Phase A: Foundation ──────────────────────
                          │   ├─ A1: hooks.go (6 types + Set*() funcs)
                          │   └─ A2: AuditLog → core/audit/model.go
                          │   └── ✓ Checkpoint: go build + go test pass
                          │
                          ├─ Phase B: Decouple ────────────────────────
                          │   ├─ B1: root.go → SetBootstrap
                          │   ├─ B2: serve.go → SetRoutes
                          │   ├─ B3: work.go + schedule_run.go → hooks
                          │   └─ B4: migrate.go + seed.go → hooks
                          │         + test file refactoring
                          │   └── ✓ Checkpoint: zero app/ imports in core/
                          │
                          ├─ Phase C: Split ───────────────────────────
                          │   ├─ C1: Delete app code from library
                          │   └─ C2: Create RapidGo-starter repo
                          │   └── ✓ Checkpoint: both repos build + test
                          │
                          ├─ Phase D: Polish ──────────────────────────
                          │   ├─ D1: rapidgo new myapp command
                          │   └─ D2: READMEs for both repos
                          │   └── ✓ Checkpoint: scaffolded app runs
                          │
                          └── v2.0.0 tag ──────────────────────────────
```

---

## Phase Dependencies

```
Phase A ──→ Phase B ──→ Phase C ──→ Phase D
  │            │            │
  │            │            └─ C2 depends on C1 (can't create starter until library is clean)
  │            │
  │            ├─ B1 must be first (root.go/NewApp is shared by B2-B4)
  │            ├─ B2, B3, B4 are independent of each other
  │            └─ B4 includes test refactoring (blocks Phase C)
  │
  ├─ A1 must be first (hooks.go used by everything after)
  └─ A2 is independent of A1 (but done on same branch lineage)
```

### Parallelization Opportunities

- **B2, B3, B4** can be developed in parallel after B1 merges (but must merge sequentially for clean git)
- **D1 and D2** are independent — can be done in parallel

---

## Branch Strategy

| Branch | Created From | Merges Into | Deleted After Merge |
|--------|-------------|-------------|:-------------------:|
| `v2` | `main` | Becomes default | No |
| `feature/v2-01-hooks-foundation` | `v2` | `v2` | Yes |
| `feature/v2-02-audit-decouple` | `v2` | `v2` | Yes |
| `feature/v2-03-root-decouple` | `v2` | `v2` | Yes |
| `feature/v2-04-serve-decouple` | `v2` | `v2` | Yes |
| `feature/v2-05-worker-decouple` | `v2` | `v2` | Yes |
| `feature/v2-06-migrate-decouple` | `v2` | `v2` | Yes |
| `feature/v2-07-remove-app-code` | `v2` | `v2` | Yes |
| `feature/v2-08-starter-repo` | `v2` | `v2` | Yes |
| `feature/v2-09-rapidgo-new-cmd` | `v2` | `v2` | Yes |
| `feature/v2-10-library-readme` | `v2` | `v2` | Yes |

### Branch Rules

- `main` stays at v1.0.0 — **frozen, no changes**
- Each feature branch is deleted after merge (clean branch strategy)
- `v2` becomes the default GitHub branch when v2.0.0 is tagged

---

## Verification Gates

Each phase must pass ALL gates before proceeding to the next phase.

### Gate A (after Phase A)

- [ ] `core/cli/hooks.go` exists with 6 type definitions and `Set*()` functions
- [ ] `hooks_test.go` passes — all setters store callbacks, defaults are nil
- [ ] `AuditLog` struct lives in `core/audit/model.go`
- [ ] `core/audit/audit.go` uses local `AuditLog` (no `database/models` import)
- [ ] `database/models/audit_log.go` re-exports `core/audit.AuditLog` for backward compat
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes (all 30+ test packages)

### Gate B (after Phase B)

- [ ] `core/cli/root.go` — no `app/providers` import
- [ ] `core/cli/serve.go` — no `routes` import
- [ ] `core/cli/work.go` — no `app/jobs`, `app/providers` imports
- [ ] `core/cli/schedule_run.go` — no `app/providers`, `app/schedule` imports
- [ ] `core/cli/migrate.go` — no `database/models` import
- [ ] `core/cli/seed.go` — no `database/seeders` import
- [ ] `cmd/main.go` — wires all 6 `cli.Set*()` hooks
- [ ] `database/models/models_test.go` — uses test-only model struct (no `User`/`Post`)
- [ ] `database/models/scopes_test.go` — uses test-only model struct
- [ ] `database/migrations/migrations_test.go` — uses test-only model struct
- [ ] `database/seeders/seeders_test.go` — uses test-only model struct
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
- [ ] `rapidgo serve` starts and serves routes correctly

### Gate C (after Phase C)

- [ ] Library repo has NO `app/`, `routes/`, `http/`, `plugins/` directories
- [ ] Library `go build ./...` passes standalone
- [ ] Library `go test ./...` passes standalone
- [ ] Library `go vet ./...` passes
- [ ] `go.mod` has no unused dependencies
- [ ] Starter repo `go build ./...` passes
- [ ] Starter `go run cmd/main.go serve` starts correctly
- [ ] Starter `go run cmd/main.go migrate` runs migrations
- [ ] Starter `go run cmd/main.go db:seed` runs seeders

### Gate D (after Phase D — release)

- [ ] `rapidgo new myapp` creates a working project
- [ ] Scaffolded project builds and runs
- [ ] Library README: `go get` instructions, package index, API overview
- [ ] Starter README: clone instructions, getting-started guide, hook wiring explained
- [ ] v2.0.0 tag pushed to GitHub
- [ ] `v2` set as default GitHub branch

---

## Document Index

| # | Document | Purpose |
|---|----------|---------|
| 1 | `v2-importable-library-master.md` | Master reference — coupling analysis, hooks, file disposition, risks |
| 2 | `v2-roadmap.md` | This document — timeline, milestones, gates |
| 3 | `v2-architecture.md` | Two-repo topology, dependency graphs, public API |
| 4 | `v2-phase-a-foundation.md` | Phase A tasks, exact code, test plan |
| 5 | `v2-phase-b-decouple.md` | Phase B tasks, exact code, test plan |
| 6 | `v2-phase-c-split.md` | Phase C tasks, file lists, test plan |
| 7 | `v2-phase-d-polish.md` | Phase D tasks, test plan |
