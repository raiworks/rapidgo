# Feature #34 — Events / Hooks System: Test Plan

## Test cases

| TC | Description | Method | Expected |
|----|-------------|--------|----------|
| TC-01 | Listen + DispatchSync invokes handler | Listen, DispatchSync | Handler called with payload |
| TC-02 | Multiple listeners all invoked | Listen x3, DispatchSync | All 3 handlers called |
| TC-03 | Dispatch fires handlers asynchronously | Listen, Dispatch, wait | Handler called in goroutine |
| TC-04 | Dispatching unknown event is no-op | Dispatch("nope", nil) | No panic |
| TC-05 | Has returns true for registered event | Listen("x", ...) | Has("x") == true |
| TC-06 | Has returns false for unregistered event | — | Has("y") == false |
| TC-07 | Payload passed correctly | Listen, DispatchSync(_, "data") | Handler receives "data" |
| TC-08 | Concurrent Listen + Dispatch is safe | Parallel goroutines | No race/panic |

## Notes

- TC-03 uses `sync.WaitGroup` to wait for async handler completion.
- TC-08 runs multiple goroutines calling Listen and Dispatch simultaneously.
- Tests verify handler invocation via counters protected by `sync.Mutex` or `atomic`.
