# 🧪 Test Plan: Redis Pub/Sub Package

> **Feature**: `67` — Redis Pub/Sub Package
> **Tasks**: [`67-redis-pubsub-tasks.md`](67-redis-pubsub-tasks.md)
> **Date**: 2026-04-20

---

## Acceptance Criteria

- [ ] `NewRedisPublisher(client)` returns a working Publisher that sends messages via Redis PUBLISH
- [ ] `NewRedisSubscriber(client, opts)` returns a Subscriber that receives messages and dispatches to handler
- [ ] Subscriber auto-reconnects on disconnect with exponential backoff
- [ ] Subscriber returns nil on context cancellation (clean shutdown)
- [ ] All tests pass hermetically using miniredis (no real Redis required)

---

## Test Cases

### TC-01: Publish and Receive

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | miniredis running, publisher and subscriber created |
| **Steps** | 1. Start subscriber in goroutine on channel "test" → 2. Publish "hello" to "test" → 3. Assert handler receives channel="test", payload="hello" |
| **Expected Result** | Handler called once with correct channel and payload |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-02: Multiple Channels

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | miniredis running |
| **Steps** | 1. Subscribe to ["ch1", "ch2"] → 2. Publish to "ch1" → 3. Publish to "ch2" → 4. Assert handler called for both |
| **Expected Result** | Handler called twice, once per channel with correct payloads |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-03: Context Cancellation (Clean Shutdown)

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Subscriber running |
| **Steps** | 1. Start subscriber with cancellable context → 2. Cancel context → 3. Assert Subscribe() returns nil |
| **Expected Result** | Subscribe returns nil (not an error) |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-04: Reconnect After Disconnect

| Property | Value |
|---|---|
| **Category** | Error / Resilience |
| **Precondition** | miniredis running, subscriber connected |
| **Steps** | 1. Subscribe in goroutine → 2. Close miniredis → 3. Restart miniredis on same addr → 4. Publish message → 5. Assert handler receives it |
| **Expected Result** | Subscriber reconnects and handler receives the post-reconnect message |
| **Status** | ⬜ Not Run |
| **Notes** | May need short sleep to allow reconnect backoff |

### TC-05: Publisher Error on Failed Publish

| Property | Value |
|---|---|
| **Category** | Error |
| **Precondition** | miniredis running |
| **Steps** | 1. Create publisher → 2. Close miniredis → 3. Call Publish() → 4. Assert error is returned |
| **Expected Result** | Publish returns a non-nil error |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-06: SubscriberOptions Defaults

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | — |
| **Steps** | 1. Call NewRedisSubscriber with zero-value SubscriberOptions → 2. Verify defaults are applied (500ms/30s/30s) |
| **Expected Result** | Options filled with documented defaults |
| **Status** | ⬜ Not Run |
| **Notes** | Verify via internal state or behavior |

---

## Edge Cases

| # | Scenario | Expected Behavior |
|---|---|---|
| 1 | Publish to channel with no subscribers | No error — Redis silently drops (PUBLISH returns 0) |
| 2 | Subscribe with empty channels slice | Subscribe returns immediately or errors gracefully |
| 3 | Handler panics | Should not crash the subscriber loop — recover and log |
| 4 | Nil logger in SubscriberOptions | Falls back to slog.Default() |

## Security Tests

No security-sensitive behavior in this feature. Auth is handled by the `*redis.Client` passed in. Channel names and payloads are the caller's responsibility.

## Performance Considerations

No performance-critical paths in this feature. Pub/sub throughput is bounded by Redis and network, not by this wrapper.

---

## Test Summary

| Category | Total | Pass | Fail | Skip |
|---|---|---|---|---|
| Happy Path | 3 | — | — | — |
| Error Cases | 2 | — | — | — |
| Edge Cases | 1 | — | — | — |
| **Total** | 6 | — | — | — |

**Result**: ⬜ NOT RUN
