package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	gql "github.com/graphql-go/graphql"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// testSchema returns a simple schema with a "hello" query that returns "world".
func testSchema() gql.Schema {
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
	return schema
}

// T01
func TestHandler_ReturnsHandlerFunc(t *testing.T) {
	handler := Handler(testSchema())
	if handler == nil {
		t.Fatal("expected non-nil handler")
	}
}

// T02
func TestHandler_PostValidQuery(t *testing.T) {
	r := gin.New()
	r.POST("/graphql", Handler(testSchema()))

	body := `{"query": "{ hello }"}`
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	data := result["data"].(map[string]interface{})
	if data["hello"] != "world" {
		t.Fatalf("expected hello=world, got %v", data["hello"])
	}
}

// T03
func TestHandler_PostWithVariables(t *testing.T) {
	fields := gql.Fields{
		"greet": &gql.Field{
			Type: gql.String,
			Args: gql.FieldConfigArgument{
				"name": &gql.ArgumentConfig{Type: gql.String},
			},
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				return "hello " + name, nil
			},
		},
	}
	rootQuery := gql.NewObject(gql.ObjectConfig{Name: "Query", Fields: fields})
	schema, _ := gql.NewSchema(gql.SchemaConfig{Query: rootQuery})

	r := gin.New()
	r.POST("/graphql", Handler(schema))

	body := `{"query": "query ($n: String) { greet(name: $n) }", "variables": {"n": "Go"}}`
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	data := result["data"].(map[string]interface{})
	if data["greet"] != "hello Go" {
		t.Fatalf("expected 'hello Go', got %v", data["greet"])
	}
}

// T04
func TestHandler_PostInvalidJSON(t *testing.T) {
	r := gin.New()
	r.POST("/graphql", Handler(testSchema()))

	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	errs, ok := result["errors"].([]interface{})
	if !ok || len(errs) == 0 {
		t.Fatal("expected errors array")
	}
	msg := errs[0].(map[string]interface{})["message"]
	if msg != "invalid request body" {
		t.Fatalf("expected 'invalid request body', got %v", msg)
	}
}

// T05
func TestHandler_PostWithOperationName(t *testing.T) {
	r := gin.New()
	r.POST("/graphql", Handler(testSchema()))

	body := `{"query": "query Hello { hello }", "operationName": "Hello"}`
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	data := result["data"].(map[string]interface{})
	if data["hello"] != "world" {
		t.Fatalf("expected hello=world, got %v", data["hello"])
	}
}

// T06
func TestHandler_GetWithQueryParam(t *testing.T) {
	r := gin.New()
	r.GET("/graphql", Handler(testSchema()))

	req := httptest.NewRequest(http.MethodGet, "/graphql?query={hello}", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	data := result["data"].(map[string]interface{})
	if data["hello"] != "world" {
		t.Fatalf("expected hello=world, got %v", data["hello"])
	}
}

// T07
func TestHandler_ResolverError(t *testing.T) {
	fields := gql.Fields{
		"fail": &gql.Field{
			Type: gql.String,
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				return nil, errors.New("something went wrong")
			},
		},
	}
	rootQuery := gql.NewObject(gql.ObjectConfig{Name: "Query", Fields: fields})
	schema, _ := gql.NewSchema(gql.SchemaConfig{Query: rootQuery})

	r := gin.New()
	r.POST("/graphql", Handler(schema))

	body := `{"query": "{ fail }"}`
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	errs, ok := result["errors"].([]interface{})
	if !ok || len(errs) == 0 {
		t.Fatal("expected errors array with resolver error")
	}
	msg := errs[0].(map[string]interface{})["message"].(string)
	if !strings.Contains(msg, "something went wrong") {
		t.Fatalf("expected error message, got %v", msg)
	}
}

// T08
func TestHandler_FromContext(t *testing.T) {
	var captured bool
	fields := gql.Fields{
		"ctx": &gql.Field{
			Type: gql.String,
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				gc, ok := FromContext(p.Context)
				captured = ok && gc != nil
				return "ok", nil
			},
		},
	}
	rootQuery := gql.NewObject(gql.ObjectConfig{Name: "Query", Fields: fields})
	schema, _ := gql.NewSchema(gql.SchemaConfig{Query: rootQuery})

	r := gin.New()
	r.POST("/graphql", Handler(schema))

	body := `{"query": "{ ctx }"}`
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !captured {
		t.Fatal("expected FromContext to return gin.Context, but got nil/false")
	}
}

// T09
func TestFromContext_NoGinContext(t *testing.T) {
	gc, ok := FromContext(context.Background())
	if ok {
		t.Fatal("expected ok=false for background context")
	}
	if gc != nil {
		t.Fatal("expected nil gin.Context")
	}
}

// T10
func TestPlayground_ReturnsHTML(t *testing.T) {
	r := gin.New()
	r.GET("/playground", Playground("Test", "/graphql"))

	req := httptest.NewRequest(http.MethodGet, "/playground", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if ct != "text/html; charset=utf-8" {
		t.Fatalf("expected text/html; charset=utf-8, got %s", ct)
	}
}

// T11
func TestPlayground_ContainsEndpoint(t *testing.T) {
	r := gin.New()
	r.GET("/playground", Playground("Test", "/my-graphql"))

	req := httptest.NewRequest(http.MethodGet, "/playground", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !bytes.Contains(w.Body.Bytes(), []byte("/my-graphql")) {
		t.Fatal("expected response to contain endpoint URL")
	}
}

// T12
func TestPlayground_ContainsTitle(t *testing.T) {
	r := gin.New()
	r.GET("/playground", Playground("My App GraphQL", "/graphql"))

	req := httptest.NewRequest(http.MethodGet, "/playground", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !bytes.Contains(w.Body.Bytes(), []byte("My App GraphQL")) {
		t.Fatal("expected response to contain title")
	}
}
