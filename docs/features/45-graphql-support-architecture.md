# 📐 Architecture: GraphQL Support

> **Feature**: `45` — GraphQL Support
> **Discussion**: [`45-graphql-support-discussion.md`](45-graphql-support-discussion.md)
> **Status**: ✅ SHIPPED
> **Date**: 2026-03-07

---

## Overview

Feature #45 adds a `core/graphql` package that provides an HTTP handler for executing GraphQL queries against a `graphql-go/graphql` schema, a Playground handler serving the GraphiQL IDE, and a `FromContext` helper for accessing the Gin context from within resolvers. The package handles request parsing (POST JSON body and GET query parameters), executes queries via the `graphql-go/graphql` engine, and returns JSON responses per the GraphQL over HTTP specification.

---

## File Structure

```
core/
  graphql/
    graphql.go       ← NEW — Handler, Playground, FromContext, Request
    graphql_test.go  ← NEW — unit tests (12 tests)
```

No existing files are modified. No migrations. No model changes.

---

## Component Design

### 1. Request Type

**Package**: `core/graphql`
**File**: `graphql.go`

```go
// Request represents a GraphQL query received from the client.
type Request struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                `json:"operationName"`
}
```

| Field | Type | Purpose |
|---|---|---|
| `Query` | `string` | The GraphQL query string |
| `Variables` | `map[string]interface{}` | Optional query variables |
| `OperationName` | `string` | Optional operation name for multi-operation documents |

---

### 2. Context Key & FromContext

```go
// contextKey is an unexported type for the gin.Context key to avoid collisions.
type contextKey struct{}

// FromContext extracts the *gin.Context from a resolver's context.Context.
// Returns the gin.Context and true if found, or nil and false otherwise.
func FromContext(ctx context.Context) (*gin.Context, bool) {
	c, ok := ctx.Value(contextKey{}).(*gin.Context)
	return c, ok
}
```

**Purpose**: Resolvers need access to request-level context (authenticated user, headers, request ID). The Handler injects the `gin.Context` into the standard `context.Context` using an unexported key. Resolvers call `FromContext(p.Context)` to retrieve it.

**Note on imports**: Since the package is named `graphql`, the external `graphql-go/graphql` library is imported with an alias throughout this file: `gql "github.com/graphql-go/graphql"`.

---

### 3. Handler

```go
// Handler returns a gin.HandlerFunc that executes GraphQL queries against
// the given schema. Supports POST (JSON body) and GET (query parameters).
//
// POST body format:
//
//	{"query": "...", "variables": {...}, "operationName": "..."}
//
// GET query parameters: ?query=...&variables={...}&operationName=...
//
// Responses follow the GraphQL over HTTP specification:
//
//	{"data": ..., "errors": [...]}
func Handler(schema gql.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Request

		if c.Request.Method == http.MethodGet {
			req.Query = c.Query("query")
			req.OperationName = c.Query("operationName")
			if v := c.Query("variables"); v != "" {
				json.Unmarshal([]byte(v), &req.Variables)
			}
		} else {
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"errors": []gin.H{{"message": "invalid request body"}},
				})
				return
			}
		}

		// Inject gin.Context into context.Context for resolvers
		ctx := context.WithValue(c.Request.Context(), contextKey{}, c)

		result := gql.Do(gql.Params{
			Schema:         schema,
			RequestString:  req.Query,
			VariableValues: req.Variables,
			OperationName:  req.OperationName,
			Context:        ctx,
		})

		c.JSON(http.StatusOK, result)
	}
}
```

**Behavior**:
- **POST**: Parses JSON body into `Request`; returns 400 if body is not valid JSON
- **GET**: Reads `query`, `variables`, `operationName` from URL query parameters; malformed `variables` JSON is silently ignored (treated as no variables), consistent with standard GraphQL HTTP implementations
- **Response**: Always HTTP 200 with `{"data": ..., "errors": [...]}` per GraphQL spec (except for transport-level errors like invalid JSON)
- **Context**: Injects `gin.Context` into `context.Context` so resolvers can access it via `FromContext`

---

### 4. Playground

```go
// Playground returns a gin.HandlerFunc that serves the GraphiQL IDE.
// title appears in the browser tab; endpoint is the URL of the GraphQL handler
// (e.g., "/graphql").
func Playground(title, endpoint string) gin.HandlerFunc {
	// Pre-render the HTML at setup time (not per-request)
	page := renderPlayground(title, endpoint)
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", page)
	}
}
```

The `renderPlayground` function uses `html/template` to safely render the title and endpoint into an HTML page that loads GraphiQL from the unpkg.com CDN:

```go
var playgroundTemplate = template.Must(template.New("playground").Parse(`<!DOCTYPE html>
<html>
<head>
  <title>{{.Title}}</title>
  <link rel="stylesheet" href="https://unpkg.com/graphiql/graphiql.min.css" />
</head>
<body style="margin:0;">
  <div id="graphiql" style="height:100vh;"></div>
  <script crossorigin src="https://unpkg.com/react/umd/react.production.min.js"></script>
  <script crossorigin src="https://unpkg.com/react-dom/umd/react-dom.production.min.js"></script>
  <script crossorigin src="https://unpkg.com/graphiql/graphiql.min.js"></script>
  <script>
    const fetcher = GraphiQL.createFetcher({url: '{{.Endpoint}}'});
    ReactDOM.render(
      React.createElement(GraphiQL, {fetcher: fetcher}),
      document.getElementById('graphiql'),
    );
  </script>
</body>
</html>`))

func renderPlayground(title, endpoint string) []byte {
	var buf bytes.Buffer
	playgroundTemplate.Execute(&buf, map[string]string{
		"Title":    title,
		"Endpoint": endpoint,
	})
	return buf.Bytes()
}
```

**Design notes**:
- HTML is rendered once at setup time (in `Playground()`), not per-request — zero allocation on each hit
- Uses `html/template` for safe escaping of title and endpoint values
- GraphiQL loaded from CDN — no embedded JavaScript in the binary

---

## Public API Summary

| Function | Signature | Purpose |
|---|---|---|
| `Handler` | `Handler(schema gql.Schema) gin.HandlerFunc` | Mount a GraphQL HTTP endpoint |
| `Playground` | `Playground(title, endpoint string) gin.HandlerFunc` | Mount the GraphiQL IDE |
| `FromContext` | `FromContext(ctx context.Context) (*gin.Context, bool)` | Extract Gin context in a resolver |

---

## Usage Example

```go
// routes/graphql.go
package routes

import (
	"github.com/RAiWorks/RapidGo/core/graphql"
	"github.com/RAiWorks/RapidGo/core/router"
	gql "github.com/graphql-go/graphql"
)

func RegisterGraphQL(r *router.Router) {
	// Define schema (application-level concern)
	fields := gql.Fields{
		"hello": &gql.Field{
			Type: gql.String,
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}
	rootQuery := gql.NewObject(gql.ObjectConfig{Name: "Query", Fields: fields})
	schema, _ := gql.NewSchema(gql.SchemaConfig{Query: rootQuery})

	// Mount handler and playground
	r.Post("/graphql", graphql.Handler(schema))
	r.Get("/graphql", graphql.Handler(schema))
	r.Get("/graphql/playground", graphql.Playground("RapidGo GraphQL", "/graphql"))
}
```

---

## Dependencies

### New External Dependency

| Package | Version | License | Purpose |
|---|---|---|---|
| `github.com/graphql-go/graphql` | latest | MIT | GraphQL schema definition and query execution engine |

### Framework Dependencies (existing)

- `github.com/gin-gonic/gin` — HTTP router (already in go.mod)

---

## Environment Variables

None. The GraphQL endpoint path and Playground availability are configured in route definitions, not environment variables.

---

## Next

Tasks → `45-graphql-support-tasks.md`
