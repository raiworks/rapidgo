# 🧪 Test Plan: Project Setup & Structure

> **Feature**: `01` — Project Setup & Structure
> **Tasks**: [`01-project-setup-tasks.md`](01-project-setup-tasks.md)
> **Date**: 2026-03-05

---

## Acceptance Criteria

The feature is DONE when ALL of these are true:

- [ ] Go module initialized with path `github.com/RAiWorks/RGo` and `go 1.21`
- [ ] All 43 directories exist in the correct hierarchy
- [ ] `cmd/main.go` compiles and runs without errors
- [ ] Running `go run ./cmd/...` prints the RGo startup banner
- [ ] `go vet ./...` reports zero issues
- [ ] `.env` exists with all placeholder configuration groups
- [ ] `.gitignore` excludes binaries, `.env.local`, IDE files, storage artifacts
- [ ] `Makefile` works: `make build`, `make run`, `make clean`
- [ ] `README.md` exists with links to documentation
- [ ] No third-party dependencies in `go.mod` (only standard library)

---

## Test Cases

### TC-01: Go Module Initialization

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Go 1.21+ installed, project root is clean |
| **Steps** | 1. Check `go.mod` exists → 2. Verify `module github.com/RAiWorks/RGo` → 3. Verify `go 1.21` directive |
| **Expected Result** | `go.mod` contains correct module path and Go version, no `require` blocks |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-02: Directory Structure Completeness

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | All Phase B tasks completed |
| **Steps** | 1. List all directories recursively → 2. Compare against architecture doc's 43-directory list → 3. Verify every leaf dir has `.gitkeep` or a `.go` file |
| **Expected Result** | All 43 directories present, no missing, no extras |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-03: Application Compiles

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `cmd/main.go` created |
| **Steps** | 1. Run `go build ./cmd/...` → 2. Check exit code is 0 → 3. Verify `bin/rgo` binary exists (via Makefile) |
| **Expected Result** | Zero compilation errors, binary produced |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-04: Application Runs and Prints Banner

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | TC-03 passes |
| **Steps** | 1. Run `go run ./cmd/...` → 2. Capture stdout |
| **Expected Result** | Output contains "RGo Framework" and "github.com/RAiWorks/RGo" |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-05: Go Vet Clean

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | All Go files created |
| **Steps** | 1. Run `go vet ./...` → 2. Check exit code |
| **Expected Result** | Exit code 0, no warnings or errors |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-06: Makefile Targets

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `Makefile` created |
| **Steps** | 1. Run `make build` → verify `bin/rgo` → 2. Run `make run` → verify banner → 3. Run `make clean` → verify `bin/` removed |
| **Expected Result** | All three targets execute successfully |
| **Status** | ⬜ Not Run |
| **Notes** | Requires `make` installed; if not available, verify equivalent manual commands |

### TC-07: `.env` Configuration File

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `.env` created |
| **Steps** | 1. Open `.env` → 2. Verify sections: Application, Database, Session, Cache, Redis, JWT, Mail, Logging, Storage, Server → 3. Verify all values are placeholders (no real secrets) |
| **Expected Result** | All 10 configuration groups present with safe placeholder values |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-08: `.gitignore` Coverage

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `.gitignore` created |
| **Steps** | 1. Verify `bin/` is ignored → 2. Verify `.env.local` is ignored → 3. Verify `.idea/` and `.vscode/` ignored → 4. Verify `storage/logs/*.log` ignored → 5. Verify `.gitkeep` files are NOT ignored |
| **Expected Result** | All patterns present, `.gitkeep` exception works |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-09: No Third-Party Dependencies

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | Module initialized |
| **Steps** | 1. Open `go.mod` → 2. Check for `require` block |
| **Expected Result** | No `require` block exists — only standard library used |
| **Status** | ⬜ Not Run |
| **Notes** | `go.sum` should not exist at this stage |

### TC-10: README Links Valid

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `README.md` created |
| **Steps** | 1. Verify link to `docs/project-context.md` → 2. Verify link to `docs/project-roadmap.md` → 3. Verify link to `docs/mastery.md` → 4. Verify link to `docs/framework/README.md` |
| **Expected Result** | All 4 linked files exist at the referenced paths |
| **Status** | ⬜ Not Run |
| **Notes** | — |

---

## Edge Cases

| # | Scenario | Expected Behavior |
|---|---|---|
| 1 | `go build` with no dependencies | Compiles with only standard library — no errors |
| 2 | Empty directories tracked by Git | `.gitkeep` files ensure Git tracks every directory |
| 3 | Running on Windows vs Linux/macOS | Makefile uses portable commands; `rm -rf` works on both via Git Bash / WSL |
| 4 | `.env` loaded when no `.env.local` exists | Framework should function with `.env` defaults alone |

## Security Tests

| # | Test | Expected |
|---|---|---|
| 1 | `.env.local` is gitignored | `git status` never shows `.env.local` even if created |
| 2 | `storage/` artifacts gitignored | Log files, cache, sessions, uploads never committed |
| 3 | No secrets in `.env` | All values are placeholder/default — no real credentials |

## Performance Considerations

N/A — this feature creates project structure only. No runtime performance characteristics.

---

## Test Summary

<!-- Fill AFTER running all tests -->

| Category | Total | Pass | Fail | Skip |
|---|---|---|---|---|
| Happy Path | 8 | — | — | — |
| Edge Cases | 4 | — | — | — |
| Security | 3 | — | — | — |
| **Total** | **15** | — | — | — |

**Result**: ⬜ NOT RUN

> To be executed during **Phase E** of [`01-project-setup-tasks.md`](01-project-setup-tasks.md).
