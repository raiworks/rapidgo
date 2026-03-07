# 📋 Tasks: GraphQL Support

> **Feature**: `45` — GraphQL Support
> **Architecture**: [`45-graphql-support-architecture.md`](45-graphql-support-architecture.md)
> **Status**: ✅ SHIPPED
> **Date**: 2026-03-07

---

## Phase A — Core Package

| # | Task | Detail |
|---|---|---|
| A1 | Add `graphql-go/graphql` dependency | `go get github.com/graphql-go/graphql` |
| A2 | Create `core/graphql/graphql.go` | `Request` struct, `contextKey` type, `FromContext()`, `Handler()`, `Playground()`, `renderPlayground()`, `playgroundTemplate` |

**Exit**: `go build` succeeds; package compiles with no errors

---

## Phase B — Tests & Verification

| # | Task | Detail |
|---|---|---|
| B1 | Create `core/graphql/graphql_test.go` | All tests from testplan (T01–T12) using `httptest` and a test schema |
| B2 | Run `go test ./core/graphql/... -v` | All 12 tests pass |
| B3 | Run `go test ./... -count=1` | All packages pass (no regressions) |
| B4 | Run `go build -o bin/rapidgo.exe ./cmd` | Binary builds successfully |

**Exit**: All tests pass, binary builds, no regressions

---

## Next

Test plan → `45-graphql-support-testplan.md`
