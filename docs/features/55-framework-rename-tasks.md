# ✅ Tasks: Framework Rename / Rebrand

> **Feature**: `55` — Framework Rename / Rebrand
> **Architecture**: [`55-framework-rename-architecture.md`](55-framework-rename-architecture.md)
> **Branch**: `feature/55-framework-rename`
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/30 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [ ] Feature branch created from latest `main`
- [x] Dependent features are merged to `main` (all 41 complete)
- [ ] Test plan doc created
- [ ] Changelog doc created (empty)

---

## Phase A — Core Module & Go Source

> Module path, all imports, string literals in Go files.

- [ ] **A.1** — Update `go.mod` module path
  - [ ] `module github.com/RAiWorks/RGo` → `module github.com/RAiWorks/RapidGo`
- [ ] **A.2** — Replace all Go import paths across all `.go` files
  - [ ] `github.com/RAiWorks/RGo/` → `github.com/RAiWorks/RapidGo/` in every import statement
  - [ ] Verify no `.go` file still contains old import path
- [ ] **A.3** — Update CLI root command (`core/cli/root.go`)
  - [ ] `Use: "rgo"` → `Use: "rapidgo"`
  - [ ] `Short: "RGo —..."` → `Short: "RapidGo —..."`
  - [ ] `const Version = "0.1.0"` → `const Version = "0.2.0"`
  - [ ] Comment: `RGo application` → `RapidGo application`
- [ ] **A.4** — Update CLI version output (`core/cli/version.go`)
  - [ ] `"Print the RGo framework"` → `"Print the RapidGo framework"`
  - [ ] `"RGo Framework v%s\n"` → `"RapidGo Framework v%s\n"`
- [ ] **A.5** — Update CLI serve banner (`core/cli/serve.go`)
  - [ ] `"github.com/RAiWorks/RGo"` → `"github.com/RAiWorks/RapidGo"`
- [ ] **A.6** — Update home controller (`http/controllers/home_controller.go`)
  - [ ] `"Welcome to RGo"` → `"Welcome to RapidGo"`
- [ ] **A.7** — Update test assertions
  - [ ] `http/controllers/controllers_test.go` — `"Welcome to RGo"` → `"Welcome to RapidGo"`
  - [ ] `core/cli/cli_test.go` — `"RGo"` test assertion → `"RapidGo"`
  - [ ] `database/database_test.go` — `"rgo_dev"` → `"rapidgo_dev"`
  - [ ] `http/responses/response_test.go` — `"name": "RGo"` → `"name": "RapidGo"` (sample data, rename for consistency)
- [ ] **A.8** — Update database defaults (`database/connection.go`)
  - [ ] `config.Env("DB_NAME", "rgo_dev")` → `config.Env("DB_NAME", "rapidgo_dev")`
- [ ] 📍 **Checkpoint A** — `go build ./...` compiles with zero errors

---

## Phase B — Configuration & Infrastructure

> `.env`, Makefile, Dockerfile, Caddyfile, LICENSE.

- [ ] **B.1** — Update `.env` defaults
  - [ ] `# RGo Framework` → `# RapidGo Framework`
  - [ ] `APP_NAME=RGo` → `APP_NAME=RapidGo`
  - [ ] `DB_NAME=rgo_dev` → `DB_NAME=rapidgo_dev`
  - [ ] `CACHE_PREFIX=rgo_` → `CACHE_PREFIX=rapidgo_`
  - [ ] `MAIL_FROM_NAME=RGo` → `MAIL_FROM_NAME=RapidGo`
- [ ] **B.2** — Update `Makefile`
  - [ ] `go build -o bin/rgo` → `go build -o bin/rapidgo`
- [ ] **B.3** — Update `Dockerfile`
  - [ ] Update any comment referencing RGo
  - [ ] Binary stays as `server` (no change needed)
- [ ] **B.4** — Update `Caddyfile`
  - [ ] `# RGo Framework` → `# RapidGo Framework`
- [ ] **B.5** — ~~Update `.gitignore` if binary name is referenced~~ **NO-OP** — `.gitignore` uses `bin/` (whole directory), not `bin/rgo`
- [ ] **B.6** — Create `LICENSE` file
  - [ ] MIT License with `Copyright (c) 2026 RAi Works (https://rai.works)`
- [ ] **B.7** — Run `go mod tidy` to regenerate `go.sum`
- [ ] 📍 **Checkpoint B** — `go build -o bin/rapidgo ./cmd` succeeds, `./bin/rapidgo version` outputs `RapidGo Framework v0.2.0`

---

## Phase C — Templates & Views

> Welcome page, any other templates.

- [ ] **C.1** — Update `resources/views/home.html`
  - [ ] Page title: `RGo Framework` → `RapidGo Framework`
  - [ ] Heading/branding text: all `RGo` → `RapidGo`
  - [ ] GitHub link: `github.com/RAiWorks/RGo` → `github.com/RAiWorks/RapidGo`
  - [ ] Footer: `RGo Framework` → `RapidGo Framework`
- [ ] 📍 **Checkpoint C** — Start server, visit localhost:8080, verify page shows "RapidGo" everywhere

---

## Phase D — Documentation

> README, project docs, feature docs, framework docs.

- [ ] **D.1** — Update `README.md`
  - [ ] All `RGo` → `RapidGo`
  - [ ] Module path, repo URL, badge links
  - [ ] Add MIT license badge
- [ ] **D.2** — Update `docs/project-context.md`
  - [ ] Project name, repo URL, module path, all text references
- [ ] **D.3** — Update `docs/project-roadmap.md`
  - [ ] Header references
- [ ] **D.4** — Update `docs/framework/service-mode-architecture.md`
  - [ ] All `RGo` → `RapidGo`
  - [ ] All `rgo` CLI references → `rapidgo`
  - [ ] All `RGO_MODE` → `RAPIDGO_MODE`
- [ ] **D.5** — Update `docs/features/56-service-mode-discussion.md`
  - [ ] All `RGo` → `RapidGo`
  - [ ] All `rgo` CLI references → `rapidgo`
  - [ ] All `RGO_MODE` → `RAPIDGO_MODE`
- [ ] **D.6** — Update all feature docs (`docs/features/01-* through 41-*`)
  - [ ] All `RGo` → `RapidGo` in content
  - [ ] All `github.com/RAiWorks/RGo` → `github.com/RAiWorks/RapidGo` in code blocks
  - [ ] All `rgo` CLI references → `rapidgo`
- [ ] **D.7** — Update `docs/framework/**/*.md` (all framework reference docs)
  - [ ] All `RGo` → `RapidGo`
- [ ] 📍 **Checkpoint D** — Grep for stale `RGo` references: `grep -r "RGo" docs/ | grep -v RapidGo` returns zero results

---

## Phase E — Testing & Validation

> Full test suite, CLI verification, stale reference check.

- [ ] **E.1** — Run full test suite: `go test ./... -count=1`
- [ ] **E.2** — Verify CLI commands
  - [ ] `./bin/rapidgo version` → `RapidGo Framework v0.2.0`
  - [ ] `./bin/rapidgo serve` → banner shows `RapidGo Framework` and `github.com/RAiWorks/RapidGo`
- [ ] **E.3** — Verify no stale references in Go source
  - [ ] `grep -r "github.com/RAiWorks/RGo[^a-zA-Z]" --include="*.go"` returns zero matches
  - [ ] `grep -r '"rgo"' --include="*.go"` returns zero matches (except RGO_TEST_* in helpers_test.go)
- [ ] **E.4** — Verify no stale references in docs
  - [ ] `grep -r "RGo" docs/ | grep -v RapidGo` returns zero matches
- [ ] **E.5** — Verify welcome page in browser
  - [ ] Build and run: `./bin/rapidgo serve`
  - [ ] Visit `http://localhost:8080` — shows "RapidGo"
  - [ ] GitHub link points to `github.com/RAiWorks/RapidGo`
- [ ] 📍 **Checkpoint E** — All tests pass, all validations green

---

## Phase F — Documentation & Cleanup

> Final review, changelog, self-review.

- [ ] **F.1** — Update changelog doc with final summary
- [ ] **F.2** — Self-review all diffs
- [ ] **F.3** — Verify `go.sum` is clean (`go mod tidy` run)
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
- [ ] Create review doc → `55-framework-rename-review.md`
- [ ] Rename GitHub repository (`RAiWorks/RGo` → `RAiWorks/RapidGo`)
