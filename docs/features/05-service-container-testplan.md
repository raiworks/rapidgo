# 🧪 Test Plan: Service Container

> **Feature**: `05` — Service Container
> **Tasks**: [`05-service-container-tasks.md`](05-service-container-tasks.md)
> **Date**: 2026-03-05

---

## Acceptance Criteria

- [ ] `Container` supports `Bind()`, `Singleton()`, `Instance()` registration
- [ ] `Make()` resolves instances first, then bindings, panics on missing
- [ ] `MustMake[T]()` resolves and casts with generics
- [ ] `Has()` returns true for registered services, false for unregistered
- [ ] `Singleton()` creates instance only once across multiple `Make()` calls
- [ ] `Bind()` creates a new instance on each `Make()` call
- [ ] Container is thread-safe under concurrent access
- [ ] `Provider` interface has `Register()` and `Boot()` methods
- [ ] `App` struct orchestrates provider lifecycle correctly
- [ ] All tests pass with `go test ./core/container/... ./core/app/...`
- [ ] `go vet ./...` reports no issues

---

## Test Cases

### TC-01: Bind registers transient factory — new instance each time

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Empty container |
| **Steps** | 1. `Bind("counter", factory)` where factory returns incrementing int → 2. `Make("counter")` twice |
| **Expected Result** | Each call returns a different value (factory called each time) |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-02: Singleton creates instance only once

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Empty container |
| **Steps** | 1. `Singleton("db", factory)` with call counter → 2. `Make("db")` three times |
| **Expected Result** | Factory called exactly once; all three calls return same instance |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-03: Instance registers pre-created object

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Empty container |
| **Steps** | 1. Create object → 2. `Instance("config", obj)` → 3. `Make("config")` |
| **Expected Result** | Returns the exact same object (pointer equality) |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-04: Make panics on unregistered service

| Property | Value |
|---|---|
| **Category** | Error |
| **Precondition** | Empty container |
| **Steps** | 1. Call `Make("nonexistent")` |
| **Expected Result** | Panics with message containing `"service not found: nonexistent"` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-05: MustMake resolves with correct type

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Container with string instance registered |
| **Steps** | 1. `Instance("greeting", "hello")` → 2. `MustMake[string](c, "greeting")` |
| **Expected Result** | Returns `"hello"` as `string` type |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-06: MustMake panics on type mismatch

| Property | Value |
|---|---|
| **Category** | Error |
| **Precondition** | Container with int instance registered |
| **Steps** | 1. `Instance("num", 42)` → 2. `MustMake[string](c, "num")` |
| **Expected Result** | Panics (interface conversion panic) |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-07: Has returns true for bound service

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Container with binding |
| **Steps** | 1. `Bind("svc", factory)` → 2. `Has("svc")` |
| **Expected Result** | Returns `true` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-08: Has returns true for instance

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Container with instance |
| **Steps** | 1. `Instance("cfg", obj)` → 2. `Has("cfg")` |
| **Expected Result** | Returns `true` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-09: Has returns false for unregistered service

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Empty container |
| **Steps** | 1. `Has("nonexistent")` |
| **Expected Result** | Returns `false` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-10: Instance takes priority over binding in Make

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | Container with same-name binding and instance |
| **Steps** | 1. `Bind("svc", factory)` → 2. `Instance("svc", obj)` → 3. `Make("svc")` |
| **Expected Result** | Returns the instance (not the factory result) — instances are checked first |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-11: Bind overwrites previous binding

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | Container with existing binding |
| **Steps** | 1. `Bind("svc", factoryA)` → 2. `Bind("svc", factoryB)` → 3. `Make("svc")` |
| **Expected Result** | Returns result from factoryB (last-write-wins) |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-12: Concurrent Make on singleton is safe

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | Container with singleton bound |
| **Steps** | 1. `Singleton("svc", factory)` → 2. Launch 100 goroutines calling `Make("svc")` concurrently |
| **Expected Result** | No race conditions, all goroutines get the same instance |
| **Status** | ⬜ Not Run |
| **Notes** | Run with `-race` flag |

### TC-13: App.Register calls provider.Register immediately

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | New App |
| **Steps** | 1. Create mock provider → 2. `app.Register(provider)` → 3. Check provider's Register was called |
| **Expected Result** | Provider's `Register()` called exactly once, service is available via `app.Make()` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-14: App.Boot calls all providers' Boot in registration order

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | App with 2 registered providers |
| **Steps** | 1. Register providerA then providerB → 2. `app.Boot()` → 3. Check boot order |
| **Expected Result** | providerA.Boot called before providerB.Boot |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-15: App.Make resolves service from container

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | App with registered provider |
| **Steps** | 1. Register provider that binds "greeting" → 2. `app.Make("greeting")` |
| **Expected Result** | Returns the service registered by the provider |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-16: Instance overwrites previous instance with same name

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | Container with existing instance |
| **Steps** | 1. `Instance("svc", objA)` → 2. `Instance("svc", objB)` → 3. `Make("svc")` |
| **Expected Result** | Returns `objB` (last-write-wins) |
| **Status** | ⬜ Not Run |
| **Notes** | Enables test mocking by replacing services |

### TC-17: Bind factory resolves another service from container

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Container with two services where one depends on the other |
| **Steps** | 1. `Instance("config", cfg)` → 2. `Bind("svc", func(c) { return NewSvc(MustMake[Config](c, "config")) })` → 3. `Make("svc")` |
| **Expected Result** | Factory receives the container, resolves "config", returns service built with that dependency |
| **Status** | ⬜ Not Run |
| **Notes** | Core DI pattern — factory uses container to resolve dependencies |

---

## Edge Cases

| # | Scenario | Expected Behavior |
|---|---|---|
| 1 | Instance overrides binding with same name | Instance returned (checked first in Make) |
| 2 | Bind called twice with same name | Last binding wins |
| 3 | 100 concurrent Make calls on singleton | Thread-safe, same instance returned |
| 4 | Provider with empty Boot() | No error, no-op |
| 5 | Instance called twice with same name | Last instance wins |
| 6 | Factory resolves other services from container | Dependency chain works correctly |

## Security Tests

| # | Test | Expected |
|---|---|---|
| 1 | Container operations under `-race` detector | No race conditions detected |

---

## Test Summary

| Category | Total | Pass | Fail | Skip |
|---|---|---|---|---|
| Happy Path | 10 | — | — | — |
| Error Cases | 2 | — | — | — |
| Edge Cases | 5 | — | — | — |
| **Total** | **17** | — | — | — |

**Result**: ⬜ NOT RUN
