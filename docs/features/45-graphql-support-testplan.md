# 🧪 Test Plan: GraphQL Support

> **Feature**: `45` — GraphQL Support
> **Tasks**: [`45-graphql-support-tasks.md`](45-graphql-support-tasks.md)
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## Test File

- `core/graphql/graphql_test.go`

---

## Unit Tests

### 1. Handler Construction

| # | Test | Expectation |
|---|---|---|
| T01 | `TestHandler_ReturnsHandlerFunc` | `Handler(schema)` returns a non-nil `gin.HandlerFunc` |

### 2. POST Requests

| # | Test | Expectation |
|---|---|---|
| T02 | `TestHandler_PostValidQuery` | POST `{"query": "{ hello }"}` → 200 with `{"data": {"hello": "world"}}` |
| T03 | `TestHandler_PostWithVariables` | POST with variables → resolver receives variable values correctly |
| T04 | `TestHandler_PostInvalidJSON` | POST with malformed JSON body → 400 with `{"errors": [{"message": "invalid request body"}]}` |
| T05 | `TestHandler_PostWithOperationName` | POST with `operationName` → correct operation selected from multi-operation document |

### 3. GET Requests

| # | Test | Expectation |
|---|---|---|
| T06 | `TestHandler_GetWithQueryParam` | GET `?query={hello}` → 200 with `{"data": {"hello": "world"}}` |

### 4. Error Handling

| # | Test | Expectation |
|---|---|---|
| T07 | `TestHandler_ResolverError` | Resolver returns error → 200 with `"errors"` array containing error message, `"data"` is null |

### 5. Context Integration

| # | Test | Expectation |
|---|---|---|
| T08 | `TestHandler_FromContext` | Resolver calls `FromContext(p.Context)` → returns the `*gin.Context` with `ok == true` |
| T09 | `TestFromContext_NoGinContext` | `FromContext(context.Background())` → returns `nil, false` |

### 6. Playground

| # | Test | Expectation |
|---|---|---|
| T10 | `TestPlayground_ReturnsHTML` | GET playground endpoint → 200, Content-Type `text/html; charset=utf-8` |
| T11 | `TestPlayground_ContainsEndpoint` | Response body contains the endpoint URL passed to `Playground()` |
| T12 | `TestPlayground_ContainsTitle` | Response body contains the title passed to `Playground()` |

---

## Acceptance Criteria

1. All 12 tests pass
2. All existing tests across all packages still pass (`go test ./...`)
3. `core/graphql` package provides `Handler`, `Playground`, `FromContext`
4. `graphql-go/graphql` added to `go.mod` as a direct dependency
5. `go build` succeeds with no errors

---

## Next

Changelog → `45-graphql-support-changelog.md`
