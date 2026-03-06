# Feature #40 — Testing Infrastructure: Discussion

## Overview

Provide reusable test utilities for framework users and internal tests: test router, test DB, HTTP assertion helpers.

## Blueprint Reference

Blueprint shows example unit and integration test patterns using `httptest`. No specific utility package is prescribed, but the patterns repeat across tests.

## Current State

32 test files exist with duplicated setup code (e.g., creating test routers, opening SQLite DBs). A shared testing utilities package eliminates this duplication.

## Deliverables

1. `testing/testutil/testutil.go` — `NewTestRouter()`, `NewTestDB()`, `DoRequest()`, `AssertStatus()`, `AssertJSONKey()`.
2. Tests for the utilities themselves.
