# 🧪 Test Plan: Configuration System

> **Feature**: `02` — Configuration System
> **Tasks**: [`02-configuration-tasks.md`](02-configuration-tasks.md)
> **Status**: ✅ Complete
> **Test Cases**: 14 (9 unit + 3 integration + 2 edge)

---

## Unit Tests — `core/config/config_test.go`

These tests run via `go test ./core/config/... -v`.

### TC-01: Load() with .env file present

| Field | Detail |
|---|---|
| **Precondition** | `.env` file exists in project root |
| **Action** | Call `config.Load()` |
| **Expected** | No error output, `.env` values available via `os.Getenv()` |
| **Pass Criteria** | Function completes without panic, env vars are set |

### TC-02: Load() without .env file

| Field | Detail |
|---|---|
| **Precondition** | No `.env` file in working directory |
| **Action** | Call `config.Load()` in a temp directory with no `.env` |
| **Expected** | Log message "No .env file found, using system environment" — no panic |
| **Pass Criteria** | Function returns gracefully, system env vars still accessible |

### TC-03: Env() with key present

| Field | Detail |
|---|---|
| **Precondition** | `os.Setenv("TEST_KEY", "hello")` |
| **Action** | Call `config.Env("TEST_KEY", "default")` |
| **Expected** | Returns `"hello"` |
| **Pass Criteria** | Return value equals `"hello"`, not fallback |

### TC-04: Env() with key absent (fallback)

| Field | Detail |
|---|---|
| **Precondition** | `TEST_MISSING` is not set |
| **Action** | Call `config.Env("TEST_MISSING", "fallback_value")` |
| **Expected** | Returns `"fallback_value"` |
| **Pass Criteria** | Return value equals `"fallback_value"` |

### TC-05: EnvInt() with valid integer

| Field | Detail |
|---|---|
| **Precondition** | `os.Setenv("TEST_INT", "42")` |
| **Action** | Call `config.EnvInt("TEST_INT", 0)` |
| **Expected** | Returns `42` |
| **Pass Criteria** | Return value equals `42` |

### TC-06: EnvInt() with invalid string

| Field | Detail |
|---|---|
| **Precondition** | `os.Setenv("TEST_INT_BAD", "not_a_number")` |
| **Action** | Call `config.EnvInt("TEST_INT_BAD", 99)` |
| **Expected** | Returns `99` (fallback) |
| **Pass Criteria** | Return value equals `99` |

### TC-07: EnvBool() truthy values

| Field | Detail |
|---|---|
| **Precondition** | `os.Setenv("TEST_BOOL_T", "true")`, `os.Setenv("TEST_BOOL_1", "1")` |
| **Action** | Call `config.EnvBool("TEST_BOOL_T", false)` and `config.EnvBool("TEST_BOOL_1", false)` |
| **Expected** | Both return `true` |
| **Pass Criteria** | Both calls return `true` |

### TC-08: EnvBool() falsy values

| Field | Detail |
|---|---|
| **Precondition** | `os.Setenv("TEST_BOOL_F", "false")`, `os.Setenv("TEST_BOOL_0", "0")` |
| **Action** | Call `config.EnvBool("TEST_BOOL_F", true)` and `config.EnvBool("TEST_BOOL_0", true)` |
| **Expected** | Both return `false` (not truthy) |
| **Pass Criteria** | Both calls return `false` |

### TC-09: Environment detection functions

| Field | Detail |
|---|---|
| **Precondition** | `os.Setenv("APP_ENV", "production")` |
| **Action** | Call `IsProduction()`, `IsDevelopment()`, `IsTesting()` |
| **Expected** | `IsProduction()` → `true`, others → `false` |
| **Pass Criteria** | Only `IsProduction()` returns `true` |

---

## Integration Tests — Manual Verification

### TC-10: go build succeeds

| Field | Detail |
|---|---|
| **Action** | Run `go build -o bin/rgo.exe ./cmd/` |
| **Expected** | Binary produced at `bin/rgo.exe` with zero errors |
| **Pass Criteria** | Exit code 0, binary exists |

### TC-11: go run displays banner with config values

| Field | Detail |
|---|---|
| **Precondition** | `.env` file has `APP_NAME=RGo`, `APP_PORT=8080`, `APP_ENV=development`, `APP_DEBUG=true` |
| **Action** | Run `go run ./cmd/` |
| **Expected** | Banner shows "RGo Framework", "Environment: development", "Port: 8080", "Debug: true" |
| **Pass Criteria** | All 4 config values appear correctly in output |

### TC-12: go vet passes

| Field | Detail |
|---|---|
| **Action** | Run `go vet ./...` |
| **Expected** | Zero warnings across all packages |
| **Pass Criteria** | Exit code 0, no output |

---

## Edge Cases

### TC-13: Env() with empty string value

| Field | Detail |
|---|---|
| **Precondition** | `os.Setenv("TEST_EMPTY", "")` |
| **Action** | Call `config.Env("TEST_EMPTY", "fallback")` |
| **Expected** | Returns `"fallback"` (empty string treated as unset) |
| **Pass Criteria** | Return value equals `"fallback"` |

### TC-14: EnvBool() with empty (fallback)

| Field | Detail |
|---|---|
| **Precondition** | `TEST_BOOL_EMPTY` not set |
| **Action** | Call `config.EnvBool("TEST_BOOL_EMPTY", true)` |
| **Expected** | Returns `true` (fallback) |
| **Pass Criteria** | Return value equals `true` |

---

## Execution Notes

- All unit tests (TC-01 through TC-09, TC-13, TC-14) run via `go test ./core/config/... -v`
- Each test must use `t.Setenv()` (Go 1.17+) which automatically restores the original value after the test
- Integration tests (TC-10, TC-11, TC-12) run manually in terminal
