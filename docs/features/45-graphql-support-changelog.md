# 📝 Changelog: GraphQL Support

> **Feature**: `45` — GraphQL Support
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## Added

- `core/graphql/graphql.go` — GraphQL HTTP handler: `Handler()`, `Playground()`, `FromContext()`
- `core/graphql/graphql_test.go` — 12 unit tests for handler, playground, and context integration

## Dependencies

- `github.com/graphql-go/graphql` — GraphQL execution engine (MIT)

## Files

| File | Action |
|---|---|
| `core/graphql/graphql.go` | NEW |
| `core/graphql/graphql_test.go` | NEW |
| `go.mod` | MODIFIED (new dependency) |
| `go.sum` | MODIFIED (new checksums) |

## Migration Guide

- No migrations required
- No new environment variables
- No breaking changes — existing code is unaffected
- Import `core/graphql` and call `graphql.Handler(schema)` to serve a GraphQL endpoint
