# 🧪 Test Plan: CSRF Protection

> **Feature**: `24` — CSRF Protection
> **Total Test Cases**: 11

---

## Token Generation

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-01 | GET generates token | GET request, no session token | Token set in `csrf_token` context key, 64-char hex |
| TC-10 | Token persists in session | Same session data across two calls | Same token returned |

## Safe Methods (Skip Validation)

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-06 | HEAD skips validation | HEAD request, no token submitted | 200 OK |
| TC-07 | OPTIONS skips validation | OPTIONS request, no token submitted | 200 OK |

## State-Changing Methods (Validate Token)

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-02 | POST valid form token | POST with `_csrf_token` form field matching | 200 OK |
| TC-03 | POST valid header token | POST with `X-CSRF-Token` header matching | 200 OK |
| TC-04 | POST missing token | POST with no token | 403 `{"error": "CSRF token mismatch"}` |
| TC-05 | POST wrong token | POST with incorrect token | 403 `{"error": "CSRF token mismatch"}` |
| TC-08 | PUT valid token | PUT with valid header token | 200 OK |
| TC-09 | DELETE missing token | DELETE with no token | 403 |

## Alias

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-11 | "csrf" alias resolvable | `ResolveAlias("csrf")` | Non-nil handler |

---

## Acceptance Criteria

1. All 11 tests pass
2. Full regression (`go test ./... -count=1`) — 0 failures
3. `go vet ./...` — clean
4. No new dependencies in `go.mod`
