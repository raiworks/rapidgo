# 📋 Tasks: Helpers

> **Feature**: `19` — Helpers
> **Architecture**: [`19-helpers-architecture.md`](19-helpers-architecture.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Phase A — Implementation

| # | Task | File | Est. |
|---|---|---|---|
| A1 | Create `password.go` — `HashPassword`, `CheckPassword` | `app/helpers/password.go` | 5m |
| A2 | Create `random.go` — `RandomString` | `app/helpers/random.go` | 3m |
| A3 | Create `string.go` — `Slugify`, `Truncate`, `Contains`, `Title`, `Excerpt`, `StripHTML`, `Mask` | `app/helpers/string.go` | 10m |
| A4 | Create `number.go` — `FormatBytes`, `Clamp` | `app/helpers/number.go` | 5m |
| A5 | Create `time.go` — `TimeAgo`, `FormatDate` | `app/helpers/time.go` | 5m |
| A6 | Create `data.go` — `StructToMap`, `MapKeys` | `app/helpers/data.go` | 5m |
| A7 | Create `env.go` — `Env` | `app/helpers/env.go` | 3m |
| A8 | Make `golang.org/x/crypto` a direct dependency | `go.mod` | 2m |

## Phase B — Tests

| # | Task | File | Est. |
|---|---|---|---|
| B1 | Create test file with tests for all 17 functions | `app/helpers/helpers_test.go` | 15m |

## Phase C — Quality Assurance

| # | Task | Est. |
|---|---|---|
| C1 | Run full test suite — 0 failures | 3m |
| C2 | Cross-check implementation vs architecture | 5m |
| C3 | Document deviations (if any) | 3m |
| C4 | Commit and push to `feature/19-helpers` | 2m |
| C5 | Merge to main and tag | 2m |
