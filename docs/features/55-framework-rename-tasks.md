# ✅ Tasks: Framework Rename / Rebrand

> **Feature**: `55` — Framework Rename / Rebrand
> **Architecture**: [`55-framework-rename-architecture.md`](55-framework-rename-architecture.md)
> **Branch**: `feature/55-framework-rename`
> **Status**: � COMPLETE
> **Progress**: 30/30 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [x] Feature branch created from latest `main`
- [x] Dependent features are merged to `main` (all 41 complete)
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — Core Module & Go Source

> Module path, all imports, string literals in Go files.

- [x] **A.1** — Update `go.mod` module path
  - [x] `module github.com/RAiWorks/RGo` → `module github.com/RAiWorks/RapidGo`
- [x] **A.2** — Replace all Go import paths across all `.go` files
  - [x] `github.com/RAiWorks/RGo/` → `github.com/RAiWorks/RapidGo/` in every import statement
  - [x] Verify no `.go` file still contains old import path
- [x] **A.3** — Update CLI root command (`core/cli/root.go`)
  - [x] `Use: "rgo"` → `Use: "rapidgo"`
  - [x] `Short: "RGo —..."` → `Short: "RapidGo —..."`
  - [x] `const Version = "0.1.0"` → `const Version = "0.2.0"`
  - [x] Comment: `RGo application` → `RapidGo application`
- [x] **A.4** — Update CLI version output (`core/cli/version.go`)
  - [x] `"Print the RGo framework"` → `"Print the RapidGo framework"`
  - [x] `"RGo Framework v%s\n"` → `"RapidGo Framework v%s\n"`
- [x] **A.5** — Update CLI serve banner (`core/cli/serve.go`)
  - [x] `"github.com/RAiWorks/RGo"` → `"github.com/RAiWorks/RapidGo"`
- [x] **A.6** — Update home controller (`http/controllers/home_controller.go`)
  - [x] `"Welcome to RGo"` → `"Welcome to RapidGo"`
- [x] **A.7** — Update test assertions
  - [x] `http/controllers/controllers_test.go` — `"Welcome to RGo"` → `"Welcome to RapidGo"`
  - [x] `core/cli/cli_test.go` — `"RGo"` test assertion → `"RapidGo"`
  - [x] `database/database_test.go` — `"rgo_dev"` → `"rapidgo_dev"`
  - [x] `http/responses/response_test.go` — `"name": "RGo"` → `"name": "RapidGo"` (sample data, rename for consistency)
- [x] **A.8** — Update database defaults (`database/connection.go`)
  - [x] `config.Env("DB_NAME", "rgo_dev")` → `config.Env("DB_NAME", "rapidgo_dev")`
- [x] 📍 **Checkpoint A** — `go build ./...` compiles with zero errors

---

## Phase B — Configuration & Infrastructure

> `.env`, Makefile, Dockerfile, Caddyfile, LICENSE.

- [x] **B.1** — Update `.env` defaults
  - [x] `# RGo Framework` → `# RapidGo Framework`
  - [x] `APP_NAME=RGo` → `APP_NAME=RapidGo`
  - [x] `DB_NAME=rgo_dev` → `DB_NAME=rapidgo_dev`
  - [x] `CACHE_PREFIX=rgo_` → `CACHE_PREFIX=rapidgo_`
  - [x] `MAIL_FROM_NAME=RGo` → `MAIL_FROM_NAME=RapidGo`
- [x] **B.2** — Update `Makefile`
  - [x] `go build -o bin/rgo` → `go build -o bin/rapidgo`
- [x] **B.3** — Update `Dockerfile`
  - [x] ~~Update any comment referencing RGo~~ No RGo references found — no change needed
  - [x] Binary stays as `server` (no change needed)
- [x] **B.4** — Update `Caddyfile`
  - [x] `# RGo Framework` → `# RapidGo Framework`
- [x] **B.5** — ~~Update `.gitignore` if binary name is referenced~~ **NO-OP** — `.gitignore` uses `bin/` (whole directory), not `bin/rgo`
- [x] **B.6** — Create `LICENSE` file
  - [x] MIT License with `Copyright (c) 2026 RAi Works (https://rai.works)`
- [x] **B.7** — Run `go mod tidy` to regenerate `go.sum`
- [x] 📍 **Checkpoint B** — `go build -o bin/rapidgo ./cmd` succeeds, `./bin/rapidgo version` outputs `RapidGo Framework v0.2.0`

---

## Phase C — Templates & Views

> Welcome page, any other templates.

- [x] **C.1** — Update `resources/views/home.html`
  - [x] Page title: `RGo Framework` → `RapidGo Framework`
  - [x] Heading/branding text: all `RGo` → `RapidGo`
  - [x] GitHub link: `github.com/RAiWorks/RGo` → `github.com/RAiWorks/RapidGo`
  - [x] Footer: `RGo Framework` → `RapidGo Framework`
- [x] 📍 **Checkpoint C** — Start server, visit localhost:8080, verify page shows "RapidGo" everywhere

---

## Phase D — Documentation

> README, project docs, feature docs, framework docs.

- [x] **D.1** — Update `README.md`
  - [x] All `RGo` → `RapidGo`
  - [x] Module path, repo URL, badge links
  - [x] Add MIT license badge
- [x] **D.2** — Update `docs/project-context.md`
  - [x] Project name, repo URL, module path, all text references
- [x] **D.3** — Update `docs/project-roadmap.md`
  - [x] Header references
- [x] **D.4** — Update `docs/framework/service-mode-architecture.md`
  - [x] All `RGo` → `RapidGo`
  - [x] All `rgo` CLI references → `rapidgo`
  - [x] All `RGO_MODE` → `RAPIDGO_MODE`
- [x] **D.5** — Update `docs/features/56-service-mode-discussion.md`
  - [x] All `RGo` → `RapidGo`
  - [x] All `rgo` CLI references → `rapidgo`
  - [x] All `RGO_MODE` → `RAPIDGO_MODE`
- [x] **D.6** — Update all feature docs (`docs/features/01-* through 41-*`)
  - [x] All `RGo` → `RapidGo` in content
  - [x] All `github.com/RAiWorks/RGo` → `github.com/RAiWorks/RapidGo` in code blocks
  - [x] All `rgo` CLI references → `rapidgo`
- [x] **D.7** — Update `docs/framework/**/*.md` (all framework reference docs)
  - [x] All `RGo` → `RapidGo`
- [x] 📍 **Checkpoint D** — Grep for stale `RGo` references: `grep -r "RGo" docs/ | grep -v RapidGo` returns zero results

---

## Phase E — Testing & Validation

> Full test suite, CLI verification, stale reference check.

- [x] **E.1** — Run full test suite: `go test ./... -count=1`
- [x] **E.2** — Verify CLI commands
  - [x] `./bin/rapidgo version` → `RapidGo Framework v0.2.0`
  - [x] `./bin/rapidgo serve` → banner shows `RapidGo Framework` and `github.com/RAiWorks/RapidGo`
- [x] **E.3** — Verify no stale references in Go source
  - [x] `grep -r "github.com/RAiWorks/RGo[^a-zA-Z]" --include="*.go"` returns zero matches
  - [x] `grep -r '"rgo"' --include="*.go"` returns zero matches (except RGO_TEST_* in helpers_test.go)
- [x] **E.4** — Verify no stale references in docs
  - [x] `grep -r "RGo" docs/ | grep -v RapidGo` returns zero matches
- [x] **E.5** — Verify welcome page in browser
  - [x] Build and run: `./bin/rapidgo serve`
  - [x] Visit `http://localhost:8080` — shows "RapidGo"
  - [x] GitHub link points to `github.com/RAiWorks/RapidGo`
- [x] 📍 **Checkpoint E** — All tests pass, all validations green

---

## Phase F — Documentation & Cleanup

> Final review, changelog, self-review.

- [x] **F.1** — Update changelog doc with final summary
- [x] **F.2** — Self-review all diffs
- [x] **F.3** — Verify `go.sum` is clean (`go mod tidy` run)
- [x] 📍 **Checkpoint F** — Clean code, complete docs, ready to ship

---

## Ship 🚀

- [x] All phases complete
- [x] All checkpoints verified
- [x] Final commit with descriptive message
- [x] Push to feature branch
- [x] Merge to `main`
- [x] Push `main`
- [x] **Keep the feature branch** — do not delete
- [x] Create review doc → `55-framework-rename-review.md`
- [x] Rename GitHub repository (`RAiWorks/RGo` → `RAiWorks/RapidGo`)
