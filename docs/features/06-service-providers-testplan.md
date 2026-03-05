# 🧪 Test Plan: Service Providers

> **Feature**: `06` — Service Providers
> **Tasks**: [`06-service-providers-tasks.md`](06-service-providers-tasks.md)
> **Date**: 2026-03-06

---

## Acceptance Criteria

- [ ] `ConfigProvider` implements `container.Provider` interface
- [ ] `LoggerProvider` implements `container.Provider` interface
- [ ] `ConfigProvider.Register()` calls `config.Load()` — env vars become available
- [ ] `LoggerProvider.Boot()` calls `logger.Setup()` — slog configured
- [ ] `cmd/main.go` uses `app.New()` → `Register()` → `Boot()` pattern
- [ ] Application output is identical to before (banner + log line)
- [ ] `go vet ./...` reports no issues
- [ ] All tests pass with `go test ./app/providers/...`

---

## Test Cases

### TC-01: ConfigProvider implements Provider interface

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Create `ConfigProvider{}` → 2. Assign to `container.Provider` variable |
| **Expected Result** | Compiles without error — interface is satisfied |
| **Status** | ⬜ Not Run |
| **Notes** | Compile-time check via `var _ container.Provider = (*ConfigProvider)(nil)` |

### TC-02: LoggerProvider implements Provider interface

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Create `LoggerProvider{}` → 2. Assign to `container.Provider` variable |
| **Expected Result** | Compiles without error — interface is satisfied |
| **Status** | ⬜ Not Run |
| **Notes** | Compile-time check via `var _ container.Provider = (*LoggerProvider)(nil)` |

### TC-03: ConfigProvider.Register loads config

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `.env` file exists with `APP_NAME=RGo` |
| **Steps** | 1. Create container → 2. Call `ConfigProvider.Register(c)` → 3. Read `config.Env("APP_NAME", "")` |
| **Expected Result** | Returns `"RGo"` — config is loaded |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-04: LoggerProvider.Boot sets up logger

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Config loaded (via ConfigProvider.Register) |
| **Steps** | 1. Load config → 2. Call `LoggerProvider.Boot(c)` → 3. Call `slog.Info("test")` |
| **Expected Result** | No panic, logger is configured |
| **Status** | ⬜ Not Run |
| **Notes** | Integration test — logger reads LOG_LEVEL from env |

### TC-05: Full App bootstrap with both providers

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `.env` file exists |
| **Steps** | 1. `app.New()` → 2. Register ConfigProvider → 3. Register LoggerProvider → 4. `app.Boot()` → 5. Read `config.Env("APP_NAME", "")` |
| **Expected Result** | Config loaded, logger initialized, no panics |
| **Status** | ⬜ Not Run |
| **Notes** | End-to-end bootstrap test |

### TC-06: ConfigProvider.Boot is no-op

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | Container exists |
| **Steps** | 1. Call `ConfigProvider.Boot(c)` |
| **Expected Result** | No panic, no side effects |
| **Status** | ⬜ Not Run |
| **Notes** | Verifies empty Boot() doesn't break |

### TC-07: LoggerProvider.Register is no-op

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | Container exists |
| **Steps** | 1. Call `LoggerProvider.Register(c)` |
| **Expected Result** | No panic, no side effects |
| **Status** | ⬜ Not Run |
| **Notes** | Verifies empty Register() doesn't break |

### TC-08: Provider registration order — Config before Logger

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `.env` exists |
| **Steps** | 1. `app.New()` → 2. Register ConfigProvider then LoggerProvider → 3. `Boot()` → 4. Verify `slog.Info` works |
| **Expected Result** | Logger reads config values correctly because config was loaded during ConfigProvider.Register() |
| **Status** | ⬜ Not Run |
| **Notes** | Validates the Register→Register→Boot→Boot lifecycle |

---

## Edge Cases

| # | Scenario | Expected Behavior |
|---|---|---|
| 1 | ConfigProvider.Boot() called | No-op, no error |
| 2 | LoggerProvider.Register() called | No-op, no error |
| 3 | Config loaded before logger | Logger reads correct LOG_LEVEL etc. |

## Test Summary

| Category | Total | Pass | Fail | Skip |
|---|---|---|---|---|
| Happy Path | 6 | — | — | — |
| Edge Cases | 2 | — | — | — |
| **Total** | **8** | — | — | — |

**Result**: ⬜ NOT RUN
