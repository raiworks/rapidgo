---
title: "Integration Tests"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Integration Tests

## Abstract

This document covers HTTP handler integration testing with
`net/http/httptest`, including testing public and protected routes.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Basic HTTP Test](#2-basic-http-test)
3. [Request and Recorder Pattern](#3-request-and-recorder-pattern)
4. [Testing Protected Routes](#4-testing-protected-routes)
5. [References](#5-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. Basic HTTP Test

```go
package integration_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "yourframework/core/router"
)

func TestHealthEndpoint(t *testing.T) {
    r := router.SetupRouter()

    req := httptest.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("GET /health returned %d, want %d", w.Code, http.StatusOK)
    }
}
```

## 3. Request and Recorder Pattern

Every integration test follows the same pattern:

```go
// 1. Set up the router
r := router.SetupRouter()

// 2. Create a request
req := httptest.NewRequest(method, path, body)
req.Header.Set("Content-Type", "application/json")

// 3. Create a response recorder
w := httptest.NewRecorder()

// 4. Serve the request
r.ServeHTTP(w, req)

// 5. Assert the response
if w.Code != expectedStatus {
    t.Errorf("got %d, want %d", w.Code, expectedStatus)
}
```

### Checking Response Body

```go
var response map[string]interface{}
json.NewDecoder(w.Body).Decode(&response)

if response["status"] != "ok" {
    t.Errorf("expected status ok, got %v", response["status"])
}
```

## 4. Testing Protected Routes

For routes behind JWT auth, include a valid token in the request:

```go
func TestProtectedEndpoint(t *testing.T) {
    r := router.SetupRouter()

    // Generate a test token
    token, _ := auth.GenerateToken(1, "test@example.com")

    req := httptest.NewRequest("GET", "/api/v1/profile", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("GET /api/v1/profile returned %d, want %d",
            w.Code, http.StatusOK)
    }
}
```

### Testing POST with JSON Body

```go
func TestCreateArticle(t *testing.T) {
    r := router.SetupRouter()
    token, _ := auth.GenerateToken(1, "test@example.com")

    body := strings.NewReader(`{"title":"Test","body":"Content"}`)
    req := httptest.NewRequest("POST", "/api/v1/articles", body)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusCreated {
        t.Errorf("POST /api/v1/articles returned %d, want %d",
            w.Code, http.StatusCreated)
    }
}
```

## 5. References

- [Testing Overview](testing-overview.md)
- [Unit Tests](unit-tests.md)
- [Authentication](../security/authentication.md)
- [Responses](../http/responses.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
