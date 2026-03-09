---
title: "Testing Overview"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Testing Overview

## Abstract

This document describes the testing strategy — Go's built-in
`testing` package, `httptest` for HTTP handlers, directory structure,
and commands.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Approach](#2-approach)
3. [Directory Structure](#3-directory-structure)
4. [Test Commands](#4-test-commands)
5. [Unit vs Integration Tests](#5-unit-vs-integration-tests)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Unit test** — Tests a single function or service in isolation.
- **Integration test** — Tests an HTTP handler with a real router.

## 2. Approach

The framework uses Go's built-in `testing` package — no third-party
test framework is required. HTTP handlers are tested with
`net/http/httptest`.

## 3. Directory Structure

```text
tests/
├── unit/           # Service and helper tests
│   └── helpers_test.go
└── integration/    # HTTP handler tests
    └── health_test.go
```

Tests **MUST** use the `_test` suffix in the package name.

## 4. Test Commands

```bash
# Run all tests
go test ./tests/... -v

# Run with coverage
go test ./... -cover

# Run a specific test
go test ./tests/unit/ -run TestHashPassword -v
```

## 5. Unit vs Integration Tests

| Aspect | Unit | Integration |
|--------|------|-------------|
| Scope | Single function | Full HTTP request |
| Speed | Fast (ms) | Moderate |
| Dependencies | None / mocked | Router, middleware |
| Package | `unit_test` | `integration_test` |
| Tools | `testing` | `testing` + `httptest` |

## 6. Security Considerations

- Tests **SHOULD** use separate test databases or in-memory SQLite to
  avoid modifying production data.
- Sensitive test fixtures (API keys, passwords) **MUST NOT** be
  committed to source control.

## 7. References

- [Unit Tests](unit-tests.md)
- [Integration Tests](integration-tests.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
