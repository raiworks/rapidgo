# 📋 Tasks: API Versioning

> **Feature**: `47` — API Versioning
> **Architecture**: [`47-api-versioning-architecture.md`](47-api-versioning-architecture.md)
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## Phase A — Core Package

| # | Task | Detail |
|---|---|---|
| A1 | Create `core/router/version.go` | `Version()`, `DeprecatedVersion()`, `deprecationHeaders()` |

**Exit**: `go build` succeeds; package compiles with no errors

---

## Phase B — Tests & Verification

| # | Task | Detail |
|---|---|---|
| B1 | Create `core/router/version_test.go` | All tests from testplan (T01–T10) using `httptest` |
| B2 | Run `go test ./core/router/... -v` | All 10 tests pass |
| B3 | Run `go test ./... -count=1` | All packages pass (no regressions) |
| B4 | Run `go build -o bin/rapidgo.exe ./cmd` | Binary builds successfully |

**Exit**: All tests pass, binary builds, no regressions

---

## Next

Test plan → `47-api-versioning-testplan.md`
