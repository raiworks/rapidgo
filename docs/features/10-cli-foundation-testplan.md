# 🧪 Test Plan: CLI Foundation

> **Feature**: `10` — CLI Foundation
> **Tasks**: [`10-cli-foundation-tasks.md`](10-cli-foundation-tasks.md)
> **Date**: 2026-03-06

---

## Acceptance Criteria

- [ ] `Execute()` runs without error when invoked with valid subcommands
- [ ] Root command (`rgo`) displays help text when run with no args
- [ ] `serve` command starts the HTTP server with default port from `APP_PORT`
- [ ] `serve --port 9090` overrides the env port
- [ ] `version` command prints correct version string
- [ ] `NewApp()` returns a fully booted `*app.Application`
- [ ] `NewApp()` application has all 5 providers registered
- [ ] All tests pass with `go test ./core/cli/...`
- [ ] All tests pass with `go test ./...` (full regression)
- [ ] `go vet ./...` reports no issues

---

## Test Cases

### TC-01: Version constant is set

**File**: `core/cli/cli_test.go`
**Function**: `TestVersion_IsSet`

| Step | Action | Expected |
|---|---|---|
| 1 | Reference `cli.Version` | Non-empty string |
| 2 | Assert `Version != ""` | Pass |

---

### TC-02: Version command output

**File**: `core/cli/cli_test.go`
**Function**: `TestVersionCmd_Output`

| Step | Action | Expected |
|---|---|---|
| 1 | Create root command, set output buffer | Buffer captures output |
| 2 | Execute with args `["version"]` | No error |
| 3 | Assert output contains `Version` value | Match |

---

### TC-03: Root command shows help (no args)

**File**: `core/cli/cli_test.go`
**Function**: `TestRootCmd_Help`

| Step | Action | Expected |
|---|---|---|
| 1 | Create root command, set output buffer | Buffer captures output |
| 2 | Execute with no args | No error |
| 3 | Assert output contains "RGo" | Match |
| 4 | Assert output contains "serve" | Match |
| 5 | Assert output contains "version" | Match |

---

### TC-04: Serve command has port flag

**File**: `core/cli/cli_test.go`
**Function**: `TestServeCmd_HasPortFlag`

| Step | Action | Expected |
|---|---|---|
| 1 | Look up "port" flag on `serveCmd` | Non-nil flag |
| 2 | Assert shorthand is "p" | Match |

---

### TC-05: NewApp returns booted application

**File**: `core/cli/cli_test.go`
**Function**: `TestNewApp_ReturnsBootedApp`

| Step | Action | Expected |
|---|---|---|
| 1 | Set `APP_ENV=testing` via `t.Setenv` | Env set |
| 2 | Call `NewApp()` | Returns `*app.Application` |
| 3 | Assert result is not nil | Pass |
| 4 | Assert `app.Container` is not nil | Pass |
| 5 | Assert `app.Container.Has("router")` | `true` |
| 6 | Assert `app.Container.Has("db")` | `true` |

---

### TC-06: Serve command registered on root

**File**: `core/cli/cli_test.go`
**Function**: `TestRootCmd_HasServeCommand`

| Step | Action | Expected |
|---|---|---|
| 1 | Iterate `rootCmd.Commands()` | List of subcommands |
| 2 | Find command with `Use == "serve"` | Found |

---

### TC-07: Version command registered on root

**File**: `core/cli/cli_test.go`
**Function**: `TestRootCmd_HasVersionCommand`

| Step | Action | Expected |
|---|---|---|
| 1 | Iterate `rootCmd.Commands()` | List of subcommands |
| 2 | Find command with `Use == "version"` | Found |

---

## Test Summary

| TC | Function | Type | Status |
|---|---|---|---|
| TC-01 | `TestVersion_IsSet` | Unit | ⬜ |
| TC-02 | `TestVersionCmd_Output` | Unit | ⬜ |
| TC-03 | `TestRootCmd_Help` | Unit | ⬜ |
| TC-04 | `TestServeCmd_HasPortFlag` | Unit | ⬜ |
| TC-05 | `TestNewApp_ReturnsBootedApp` | Integration | ⬜ |
| TC-06 | `TestRootCmd_HasServeCommand` | Unit | ⬜ |
| TC-07 | `TestRootCmd_HasVersionCommand` | Unit | ⬜ |

**Total**: 7 test cases
