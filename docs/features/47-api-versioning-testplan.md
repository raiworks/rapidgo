# 🧪 Test Plan: API Versioning

> **Feature**: `47` — API Versioning
> **Tasks**: [`47-api-versioning-tasks.md`](47-api-versioning-tasks.md)
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## Test File

- `core/router/version_test.go`

---

## Unit Tests

### 1. Version Method

| # | Test | Expectation |
|---|---|---|
| T01 | `TestVersion_ReturnsRouteGroup` | `Version("v1")` returns a non-nil `*RouteGroup` |
| T02 | `TestVersion_PrefixesRoutes` | Route registered on `Version("v1")` is reachable at `/api/v1/{path}` |
| T03 | `TestVersion_MultipleVersions` | `Version("v1")` and `Version("v2")` create independent groups; routes don't overlap |
| T04 | `TestVersion_SupportsAllMethods` | GET, POST, PUT, DELETE all work on a versioned group |
| T05 | `TestVersion_NoDeprecationHeaders` | Requests to a non-deprecated version have no `Sunset` or `X-API-Deprecated` headers |

### 2. DeprecatedVersion Method

| # | Test | Expectation |
|---|---|---|
| T06 | `TestDeprecatedVersion_SetsSunsetHeader` | Response includes `Sunset` header with the provided date string |
| T07 | `TestDeprecatedVersion_SetsDeprecatedHeader` | Response includes `X-API-Deprecated: true` header |
| T08 | `TestDeprecatedVersion_PrefixesRoutes` | Route registered on `DeprecatedVersion("v1", ...)` is reachable at `/api/v1/{path}` |
| T09 | `TestDeprecatedVersion_HeadersOnAllRoutes` | Multiple routes in a deprecated version all receive deprecation headers |

### 3. Nested Groups

| # | Test | Expectation |
|---|---|---|
| T10 | `TestVersion_NestedGroup` | `v1.Group("/admin")` creates routes at `/api/v1/admin/{path}` |

---

## Acceptance Criteria

1. All 10 tests pass
2. All existing tests across all packages still pass (`go test ./...`)
3. `core/router` package exports `Version` and `DeprecatedVersion` methods on `*Router`
4. No new dependencies added
5. `go build` succeeds with no errors

---

## Next

Changelog → `47-api-versioning-changelog.md`
