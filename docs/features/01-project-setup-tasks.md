# ✅ Tasks: Project Setup & Structure

> **Feature**: `01` — Project Setup & Structure
> **Architecture**: [`01-project-setup-architecture.md`](01-project-setup-architecture.md)
> **Branch**: `feature/01-project-setup`
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/30 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [ ] Feature branch created from latest `main`
- [ ] Dependent features are merged to `main` (N/A — no dependencies)
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — Go Module Initialization

> Initialize the Go module and verify the toolchain.

- [ ] **A.1** — Verify Go version is 1.21+ (`go version`)
- [ ] **A.2** — Run `go mod init github.com/RAiWorks/RGo`
- [ ] **A.3** — Set Go version in `go.mod` to `go 1.21`
- [ ] 📍 **Checkpoint A** — `go.mod` exists with correct module path and Go version

---

## Phase B — Directory Structure

> Create the full directory tree with `.gitkeep` placeholders.

- [ ] **B.1** — Create `cmd/` directory
- [ ] **B.2** — Create `core/` directory tree (16 subdirectories):
  - [ ] `core/app/`
  - [ ] `core/container/`
  - [ ] `core/router/`
  - [ ] `core/middleware/`
  - [ ] `core/config/`
  - [ ] `core/logger/`
  - [ ] `core/errors/`
  - [ ] `core/session/`
  - [ ] `core/validation/`
  - [ ] `core/crypto/`
  - [ ] `core/cache/`
  - [ ] `core/mail/`
  - [ ] `core/events/`
  - [ ] `core/i18n/`
  - [ ] `core/server/`
  - [ ] `core/websocket/`
- [ ] **B.3** — Create `database/` directory tree:
  - [ ] `database/migrations/`
  - [ ] `database/seeders/`
  - [ ] `database/models/`
  - [ ] `database/querybuilder/`
- [ ] **B.4** — Create `app/` directory tree:
  - [ ] `app/providers/`
  - [ ] `app/services/`
  - [ ] `app/helpers/`
- [ ] **B.5** — Create `http/` directory tree:
  - [ ] `http/controllers/`
  - [ ] `http/requests/`
  - [ ] `http/responses/`
- [ ] **B.6** — Create `routes/` directory
- [ ] **B.7** — Create `resources/` directory tree:
  - [ ] `resources/views/`
  - [ ] `resources/lang/`
  - [ ] `resources/static/`
- [ ] **B.8** — Create `storage/` directory tree:
  - [ ] `storage/uploads/`
  - [ ] `storage/cache/`
  - [ ] `storage/sessions/`
  - [ ] `storage/logs/`
- [ ] **B.9** — Create `tests/` directory tree:
  - [ ] `tests/unit/`
  - [ ] `tests/integration/`
- [ ] **B.10** — Add `.gitkeep` to every leaf directory that has no Go source files
- [ ] 📍 **Checkpoint B** — All 43 directories exist, all leaf directories have `.gitkeep` or a Go file

---

## Phase C — Entry Point & Placeholder Files

> Create `main.go` and placeholder Go files.

- [ ] **C.1** — Create `cmd/main.go` with startup banner (as defined in architecture doc)
- [ ] **C.2** — Create `database/connection.go` with `package database` declaration
- [ ] **C.3** — Create `routes/web.go` with `package routes` declaration
- [ ] **C.4** — Create `routes/api.go` with `package routes` declaration
- [ ] 📍 **Checkpoint C** — `go build ./cmd/...` succeeds, `go run ./cmd/...` prints banner

---

## Phase D — Project Configuration Files

> Create `.env`, `Makefile`, `.gitignore`, and `README.md`.

- [ ] **D.1** — Create `.env` with all placeholder configuration values (grouped by subsystem)
- [ ] **D.2** — Create `Makefile` with targets: `build`, `run`, `test`, `clean`, `fmt`, `vet`, `lint`
- [ ] **D.3** — Update `.gitignore` with Go, environment, IDE, storage, and OS rules
- [ ] **D.4** — Create project `README.md` with overview and links to docs
- [ ] 📍 **Checkpoint D** — `make build` produces `bin/rgo`, `make run` prints banner, `make clean` removes `bin/`

---

## Phase E — Testing & Verification

> Execute the test plan, verify all acceptance criteria.

- [ ] **E.1** — Run test plan: `go build ./cmd/...` compiles without errors
- [ ] **E.2** — Run test plan: `go run ./cmd/...` prints startup banner
- [ ] **E.3** — Run test plan: `go vet ./...` reports no issues
- [ ] **E.4** — Run test plan: All 43 directories exist in correct hierarchy
- [ ] **E.5** — Run test plan: `.env` is parseable, `.gitignore` covers required patterns
- [ ] **E.6** — Run test plan: `make build` / `make run` / `make clean` all work
- [ ] 📍 **Checkpoint E** — All acceptance criteria met, test summary filled in testplan doc

---

## Phase F — Documentation & Cleanup

> Finalize documentation and self-review.

- [ ] **F.1** — Update changelog doc with implementation summary
- [ ] **F.2** — Update project roadmap — mark Feature #01 as ✅
- [ ] **F.3** — Self-review all diffs (every file created in this feature)
- [ ] 📍 **Checkpoint F** — Clean code, complete docs, ready to ship

---

## Ship 🚀

- [ ] All phases complete
- [ ] All checkpoints verified
- [ ] Final commit with descriptive message
- [ ] Push to feature branch
- [ ] Merge to `main`
- [ ] Push `main`
- [ ] **Keep the feature branch** — do not delete
- [ ] Create review doc → `01-project-setup-review.md`
