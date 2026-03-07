# 🧪 Test Plan: Queue Workers / Background Jobs

> **Feature**: `42` — Queue Workers / Background Jobs
> **Architecture**: [`42-queue-workers-architecture.md`](42-queue-workers-architecture.md)
> **Coverage Target**: All queue components — handler registry, dispatcher, all 4 drivers, worker pool

---

## Pre-Conditions

- [ ] All existing test packages pass before implementation begins
- [ ] No stale references or compilation errors in codebase

---

## Test Cases

### TC-01: Handler Registry

| Field | Value |
|---|---|
| **ID** | TC-01 |
| **Title** | Handler registry: register, resolve, reset |
| **Type** | Unit |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. `RegisterHandler("email", handler)`. 2. `ResolveHandler("email")` → non-nil. 3. `ResolveHandler("unknown")` → nil. 4. `ResetHandlers()`. 5. `ResolveHandler("email")` → nil. |
| **Expected** | Handlers correctly registered, resolved, and cleared |

---

### TC-02: Dispatcher — Dispatch

| Field | Value |
|---|---|
| **ID** | TC-02 |
| **Title** | Dispatcher.Dispatch marshals payload and calls driver.Push |
| **Type** | Unit |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Register handler for "test_job". 2. Create Dispatcher with MemoryDriver. 3. `Dispatch(ctx, "default", "test_job", map[string]string{"key": "value"})`. 4. Pop from memory driver. 5. Verify job fields (Queue, Type, Payload JSON, Attempts=0, AvailableAt <= now). |
| **Expected** | Job stored with correct JSON payload and metadata |

---

### TC-03: Dispatcher — DispatchDelayed

| Field | Value |
|---|---|
| **ID** | TC-03 |
| **Title** | DispatchDelayed sets future AvailableAt |
| **Type** | Unit |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Register handler. 2. DispatchDelayed with 5-minute delay. 3. Pop immediately → nil (not yet available). 4. Verify stored job has AvailableAt ~5min in future. |
| **Expected** | Job not available until delay passes |

---

### TC-04: MemoryDriver — Full Lifecycle

| Field | Value |
|---|---|
| **ID** | TC-04 |
| **Title** | MemoryDriver: Push → Pop → Delete cycle |
| **Type** | Unit |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. `Push` job. 2. `Size` → 1. 3. `Pop` → job returned with ReservedAt set. 4. `Size` → 0 (reserved). 5. `Delete` → job removed. |
| **Expected** | Complete lifecycle works correctly |

---

### TC-05: MemoryDriver — Release & Retry

| Field | Value |
|---|---|
| **ID** | TC-05 |
| **Title** | MemoryDriver: Release puts job back for retry |
| **Type** | Unit |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Push job. 2. Pop. 3. Release with 0s delay. 4. Pop again → same job with Attempts incremented. |
| **Expected** | Job re-queued with Attempts+1, ReservedAt cleared, new AvailableAt |

---

### TC-06: MemoryDriver — Fail

| Field | Value |
|---|---|
| **ID** | TC-06 |
| **Title** | MemoryDriver: Fail moves job to failed storage |
| **Type** | Unit |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Push job. 2. Pop. 3. Fail(job, error). 4. Size → 0 (removed from active). |
| **Expected** | Job removed from active queue, stored in failed jobs |

---

### TC-07: MemoryDriver — Delayed Job

| Field | Value |
|---|---|
| **ID** | TC-07 |
| **Title** | MemoryDriver: Pop respects AvailableAt for delayed jobs |
| **Type** | Unit |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Push job with AvailableAt 1 hour in future. 2. Pop → nil (not yet available). 3. Size → 1 (still pending). |
| **Expected** | Delayed job not returned by Pop until AvailableAt passes |

---

### TC-08: SyncDriver — Immediate Execution

| Field | Value |
|---|---|
| **ID** | TC-08 |
| **Title** | SyncDriver executes handler immediately on Push |
| **Type** | Unit |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Register handler that sets a flag. 2. Create SyncDriver. 3. Push job. 4. Verify flag is set (handler executed). |
| **Expected** | Handler executed synchronously during Push |

---

### TC-09: SyncDriver — Returns Handler Error

| Field | Value |
|---|---|
| **ID** | TC-09 |
| **Title** | SyncDriver returns handler error from Push |
| **Type** | Unit |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Register handler that returns error. 2. Push job. 3. Verify Push returns the handler's error. |
| **Expected** | Error propagated back to caller |

---

### TC-10: SyncDriver — Unknown Handler

| Field | Value |
|---|---|
| **ID** | TC-10 |
| **Title** | SyncDriver returns error for unknown job type |
| **Type** | Unit |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Push job with unregistered type name. 2. Verify error returned. |
| **Expected** | Error indicates handler not found |

---

### TC-11: DatabaseDriver — Push/Pop/Delete

| Field | Value |
|---|---|
| **ID** | TC-11 |
| **Title** | DatabaseDriver full lifecycle with SQLite |
| **Type** | Integration |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Create SQLite in-memory DB. 2. Auto-migrate job models. 3. Push job. 4. Pop → job returned. 5. Delete → row removed. |
| **Expected** | Database operations work correctly with SQLite |

---

### TC-12: DatabaseDriver — Release & Retry

| Field | Value |
|---|---|
| **ID** | TC-12 |
| **Title** | DatabaseDriver release re-queues job for retry |
| **Type** | Integration |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Push job. 2. Pop. 3. Release with delay. 4. Check DB: reserved_at=NULL, attempts incremented, available_at in future. |
| **Expected** | Job updated correctly in database |

---

### TC-13: DatabaseDriver — Fail

| Field | Value |
|---|---|
| **ID** | TC-13 |
| **Title** | DatabaseDriver fail stores in failed_jobs table |
| **Type** | Integration |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Push job. 2. Pop. 3. Fail(job, error). 4. Verify job deleted from jobs table. 5. Verify row exists in failed_jobs with error message. |
| **Expected** | Job moved from jobs to failed_jobs |

---

### TC-14: RedisDriver — Push/Pop/Delete

| Field | Value |
|---|---|
| **ID** | TC-14 |
| **Title** | RedisDriver full lifecycle with miniredis |
| **Type** | Integration |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Start miniredis. 2. Create RedisDriver. 3. Push job. 4. Pop → job returned. 5. Size → 0. |
| **Expected** | Redis list operations work correctly |

---

### TC-15: RedisDriver — Release

| Field | Value |
|---|---|
| **ID** | TC-15 |
| **Title** | RedisDriver release pushes job back to list |
| **Type** | Integration |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Push job. 2. Pop. 3. Release. 4. Pop again → job with incremented Attempts. |
| **Expected** | Job back in list with updated metadata |

---

### TC-16: RedisDriver — Fail

| Field | Value |
|---|---|
| **ID** | TC-16 |
| **Title** | RedisDriver fail stores in failed list |
| **Type** | Integration |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Push job. 2. Pop. 3. Fail. 4. Check failed list length → 1. |
| **Expected** | Failed job stored in separate Redis list |

---

### TC-17: RedisDriver — Size

| Field | Value |
|---|---|
| **ID** | TC-17 |
| **Title** | RedisDriver Size returns list length |
| **Type** | Unit |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Push 3 jobs. 2. Size → 3. 3. Pop 1. 4. Size → 2. |
| **Expected** | Accurate count of pending jobs |

---

### TC-18: Worker — Processes Jobs

| Field | Value |
|---|---|
| **ID** | TC-18 |
| **Title** | Worker processes jobs from memory driver |
| **Type** | Integration |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Register handler that records calls. 2. Dispatch 3 jobs. 3. Start worker with ctx. 4. Wait for processing. 5. Cancel ctx. 6. Verify all 3 jobs processed. |
| **Expected** | All dispatched jobs processed by worker |

---

### TC-19: Worker — Retry on Failure

| Field | Value |
|---|---|
| **ID** | TC-19 |
| **Title** | Worker retries failed job up to max attempts |
| **Type** | Integration |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Register handler that fails first 2 times, succeeds on 3rd. 2. Dispatch job with MaxAttempts=3. 3. Run worker. 4. Verify handler called 3 times. 5. Verify job deleted (succeeded). |
| **Expected** | Job retried until success or max attempts |

---

### TC-20: Worker — Fail After Max Attempts

| Field | Value |
|---|---|
| **ID** | TC-20 |
| **Title** | Worker moves job to failed after max attempts exceeded |
| **Type** | Integration |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Register handler that always fails. 2. Dispatch job with MaxAttempts=2. 3. Run worker. 4. Verify job moved to failed storage. |
| **Expected** | After MaxAttempts failures, Fail() called instead of Release() |

---

### TC-21: Worker — Panic Recovery

| Field | Value |
|---|---|
| **ID** | TC-21 |
| **Title** | Worker recovers from handler panic |
| **Type** | Unit |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Register handler that panics. 2. Dispatch job. 3. Start worker. 4. Verify worker continues running (doesn't crash). 5. Verify job marked as failed. |
| **Expected** | Panic caught, job failed, worker continues |

---

### TC-22: Worker — Graceful Shutdown

| Field | Value |
|---|---|
| **ID** | TC-22 |
| **Title** | Worker stops cleanly on context cancellation |
| **Type** | Integration |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Start worker. 2. Cancel context. 3. Verify Run returns nil. 4. Verify no goroutine leak (rely on WaitGroup). |
| **Expected** | Clean shutdown, nil error returned |

---

### TC-23: Worker — Job Timeout

| Field | Value |
|---|---|
| **ID** | TC-23 |
| **Title** | Worker enforces job timeout |
| **Type** | Integration |
| **File** | `core/queue/queue_test.go` |
| **Steps** | 1. Register handler that blocks for 10 seconds. 2. Set timeout to 100ms. 3. Dispatch job. 4. Run worker. 5. Verify job failed with context deadline exceeded. |
| **Expected** | Job terminated and failed when timeout exceeded |

---

### TC-24: Full Regression

| Field | Value |
|---|---|
| **ID** | TC-24 |
| **Title** | All existing tests pass with no regressions |
| **Type** | Regression |
| **File** | All packages |
| **Steps** | `go test ./... -count=1` |
| **Expected** | All 31+ packages pass. Zero failures. |

---

## Test Summary

| TC | Title | Type | Status |
|---|---|---|---|
| TC-01 | Handler registry | Unit | ⬜ |
| TC-02 | Dispatcher.Dispatch | Unit | ⬜ |
| TC-03 | Dispatcher.DispatchDelayed | Unit | ⬜ |
| TC-04 | MemoryDriver lifecycle | Unit | ⬜ |
| TC-05 | MemoryDriver release/retry | Unit | ⬜ |
| TC-06 | MemoryDriver fail | Unit | ⬜ |
| TC-07 | MemoryDriver delayed job | Unit | ⬜ |
| TC-08 | SyncDriver immediate exec | Unit | ⬜ |
| TC-09 | SyncDriver error propagation | Unit | ⬜ |
| TC-10 | SyncDriver unknown handler | Unit | ⬜ |
| TC-11 | DatabaseDriver lifecycle | Integration | ⬜ |
| TC-12 | DatabaseDriver release/retry | Integration | ⬜ |
| TC-13 | DatabaseDriver fail | Integration | ⬜ |
| TC-14 | RedisDriver lifecycle | Integration | ⬜ |
| TC-15 | RedisDriver release | Integration | ⬜ |
| TC-16 | RedisDriver fail | Integration | ⬜ |
| TC-17 | RedisDriver size | Unit | ⬜ |
| TC-18 | Worker processes jobs | Integration | ⬜ |
| TC-19 | Worker retry on failure | Integration | ⬜ |
| TC-20 | Worker fail after max attempts | Integration | ⬜ |
| TC-21 | Worker panic recovery | Unit | ⬜ |
| TC-22 | Worker graceful shutdown | Integration | ⬜ |
| TC-23 | Worker job timeout | Integration | ⬜ |
| TC-24 | Full regression | Regression | ⬜ |

---

## Known Limitations

| # | Limitation | Mitigation |
|---|---|---|
| 1 | SQLite doesn't support `SKIP LOCKED` | Fall back to non-SKIP LOCKED query; single worker recommended for SQLite |
| 2 | Redis driver doesn't support delayed jobs natively | Delayed jobs check AvailableAt on Pop; Redis BRPOP not used for delayed |
| 3 | No job priority system | Use separate named queues for priority separation |
| 4 | No dead letter queue | Failed jobs stored in failed_jobs table/list; manual retry only |
