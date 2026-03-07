package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	gql "github.com/graphql-go/graphql"
)

// Request represents a GraphQL query received from the client.
type Request struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                `json:"operationName"`
}

// contextKey is an unexported type for the gin.Context key to avoid collisions.
type contextKey struct{}

// FromContext extracts the *gin.Context from a resolver's context.Context.
// Returns the gin.Context and true if found, or nil and false otherwise.
func FromContext(ctx context.Context) (*gin.Context, bool) {
	c, ok := ctx.Value(contextKey{}).(*gin.Context)
	return c, ok
}

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

// Playground returns a gin.HandlerFunc that serves the GraphiQL IDE.
// title appears in the browser tab; endpoint is the URL of the GraphQL handler
// (e.g., "/graphql").
func Playground(title, endpoint string) gin.HandlerFunc {
	page := renderPlayground(title, endpoint)
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", page)
	}
}

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
