# 🧪 Test Plan: Redis DB Selection & Multi-Client Helper

> **Feature**: `66` — Redis DB Selection & Multi-Client Helper
> **Tasks**: [`66-redis-db-multi-client-tasks.md`](66-redis-db-multi-client-tasks.md)
> **Date**: 2026-04-20

---

## Acceptance Criteria

- [ ] `REDIS_DB` env var is honored by cache and session Redis clients
- [ ] `NewRedisClient(nil)` returns a client on the env-configured DB (default 0)
- [ ] `NewRedisClient(&db)` returns a client on the overridden DB
- [ ] Pool size and timeout env vars are respected
- [ ] Existing tests pass with zero regressions
- [ ] Starter provider updated with named-client pattern

---

## Test Cases

### TC-01: Default DB (no REDIS_DB set)

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `REDIS_DB` env not set |
| **Steps** | 1. Call `NewRedisClient(nil)` → 2. Check `client.Options().DB` |
| **Expected Result** | `DB == 0` |
| **Status** | ⬜ Not Run |

### TC-02: REDIS_DB env respected

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `REDIS_DB=3` |
| **Steps** | 1. Set env `REDIS_DB=3` → 2. Call `NewRedisClient(nil)` → 3. Check `client.Options().DB` |
| **Expected Result** | `DB == 3` |
| **Status** | ⬜ Not Run |

### TC-03: DB override takes precedence

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `REDIS_DB=3`, override `5` |
| **Steps** | 1. Set env `REDIS_DB=3` → 2. Call `NewRedisClient(&5)` → 3. Check `client.Options().DB` |
| **Expected Result** | `DB == 5` |
| **Status** | ⬜ Not Run |

### TC-04: Invalid REDIS_DB (non-integer)

| Property | Value |
|---|---|
| **Category** | Error Case |
| **Precondition** | `REDIS_DB=abc` |
| **Steps** | 1. Set env `REDIS_DB=abc` → 2. Call `NewRedisClient(nil)` → 3. Check `client.Options().DB` |
| **Expected Result** | `DB == 0` (fallback to default) |
| **Status** | ⬜ Not Run |

### TC-05: Pool size from env

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `REDIS_POOL_SIZE=20` |
| **Steps** | 1. Set env → 2. Call `NewRedisClient(nil)` → 3. Check `client.Options().PoolSize` |
| **Expected Result** | `PoolSize == 20` |
| **Status** | ⬜ Not Run |

### TC-06: Timeout from env

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `REDIS_DIAL_TIMEOUT=10s` |
| **Steps** | 1. Set env → 2. Call `NewRedisClient(nil)` → 3. Check `client.Options().DialTimeout` |
| **Expected Result** | `DialTimeout == 10s` |
| **Status** | ⬜ Not Run |

### TC-07: Existing cache tests pass

| Property | Value |
|---|---|
| **Category** | Regression |
| **Precondition** | None |
| **Steps** | `go test ./core/cache/...` |
| **Expected Result** | All existing tests pass |
| **Status** | ⬜ Not Run |

### TC-08: Existing session tests pass

| Property | Value |
|---|---|
| **Category** | Regression |
| **Precondition** | None |
| **Steps** | `go test ./core/session/...` |
| **Expected Result** | All existing tests pass |
| **Status** | ⬜ Not Run |

### TC-09: Full regression

| Property | Value |
|---|---|
| **Category** | Regression |
| **Precondition** | None |
| **Steps** | `go test ./...` |
| **Expected Result** | All tests pass |
| **Status** | ⬜ Not Run |

---

## Edge Cases

| # | Scenario | Expected Behavior |
|---|---|---|
| 1 | `REDIS_DB=-1` | Falls back to default 0 |
| 2 | `REDIS_DB=16` | Uses value as-is (Redis will reject if server doesn't support it) |
| 3 | `REDIS_POOL_SIZE=0` | Falls back to default 10 |
| 4 | `REDIS_DIAL_TIMEOUT=invalid` | Falls back to default 5s |
| 5 | Both env and override set | Override wins |

## Security Tests

| # | Test | Expected |
|---|---|---|
| 1 | `REDIS_PASSWORD` not logged | Password never appears in log output |

---

## Test Summary

| Category | Total | Pass | Fail | Skip |
|---|---|---|---|---|
| Happy Path | 5 | — | — | — |
| Error Cases | 1 | — | — | — |
| Regression | 3 | — | — | — |
| **Total** | 9 | — | — | — |

**Result**: ⬜ NOT RUN
