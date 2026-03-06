# 📝 Changelog: Controllers

> **Feature**: `15` — Controllers
> **Branch**: `feature/15-controllers`
> **Started**: —
> **Completed**: —

---

## Log

- **Phase A** — Created `http/controllers/home_controller.go` (`Home`) and `http/controllers/post_controller.go` (`PostController` with 7 methods). Updated `routes/web.go` and `routes/api.go`. Build + vet clean.
- **Phase B** — Created `http/controllers/controllers_test.go` with 9 tests (TC-01 through TC-09). All pass. Full regression: 164 tests, 0 failures.
- **Phase C** — Cross-check passed. One deviation noted.

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| External test package | `package controllers` | `package controllers_test` | `routes` imports `controllers` — using internal test package with `routes` import creates a cycle; external test package breaks it |

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
| External test package | Import cycle: `controllers_test` → `routes` → `controllers`. Using `package controllers_test` avoids the cycle while preserving all test coverage. | 2026-03-06 |
