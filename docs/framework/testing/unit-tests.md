---
title: "Unit Tests"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Unit Tests

## Abstract

This document covers unit testing patterns for services and helpers —
basic assertions, table-driven tests, and naming conventions.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Helper Testing](#2-helper-testing)
3. [Table-Driven Tests](#3-table-driven-tests)
4. [Naming Conventions](#4-naming-conventions)
5. [Assertions and Error Checking](#5-assertions-and-error-checking)
6. [References](#6-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. Helper Testing

Test pure helper functions directly:

```go
package unit_test

import (
    "testing"

    "yourframework/app/helpers"
)

func TestHashPassword(t *testing.T) {
    hash, err := helpers.HashPassword("secret123")
    if err != nil {
        t.Fatalf("HashPassword failed: %v", err)
    }
    if !helpers.CheckPassword(hash, "secret123") {
        t.Error("CheckPassword should return true for correct password")
    }
    if helpers.CheckPassword(hash, "wrong") {
        t.Error("CheckPassword should return false for wrong password")
    }
}
```

## 3. Table-Driven Tests

Use table-driven tests for functions with multiple input/output
scenarios:

```go
func TestSlugify(t *testing.T) {
    tests := []struct {
        input string
        want  string
    }{
        {"Hello World", "hello-world"},
        {"Go Web Framework!", "go-web-framework"},
        {"  spaces  ", "spaces"},
    }
    for _, tt := range tests {
        got := helpers.Slugify(tt.input)
        if got != tt.want {
            t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
        }
    }
}
```

## 4. Naming Conventions

| Convention | Example |
|------------|---------|
| Test function prefix | `Test` + function name: `TestHashPassword` |
| Subtest naming | `t.Run("empty input", ...)` |
| File naming | `<package>_test.go`: `helpers_test.go` |
| Package naming | `<package>_test`: `unit_test` |

## 5. Assertions and Error Checking

Go's `testing` package uses manual assertions:

```go
// Fatal — stops test immediately
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}

// Error — records failure, continues
if got != want {
    t.Errorf("got %v, want %v", got, want)
}

// Skip — skip test conditionally
if os.Getenv("CI") == "" {
    t.Skip("skipping integration test in local env")
}
```

## 6. References

- [Testing Overview](testing-overview.md)
- [Integration Tests](integration-tests.md)
- [Helpers Reference](../reference/helpers-reference.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
