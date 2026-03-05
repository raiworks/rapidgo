# 💬 Discussion: Configuration System

> **Feature**: `02` — Configuration System
> **Status**: 🟢 COMPLETE
> **Branch**: `feature/02-configuration`
> **Depends On**: #01 (Project Setup & Structure ✅)
> **Date Started**: 2026-03-05
> **Date Completed**: 2026-03-05

---

## Summary

Implement the framework's configuration system — load environment variables from `.env` files using godotenv, provide typed accessor helpers (`Env()`, `EnvInt()`, `EnvBool()`), and add environment detection functions (`IsProduction()`, `IsDevelopment()`, `IsTesting()`, `IsDebug()`). This is the first feature to introduce third-party dependencies and the first `core/` package with real code.

---

## Functional Requirements

- As a framework developer, I want `.env` files loaded automatically at startup so that all config keys are available via `os.Getenv()`
- As a framework developer, I want an `Env(key, fallback)` helper so that missing keys return sensible defaults instead of empty strings
- As a framework developer, I want typed helpers (`EnvInt`, `EnvBool`) so that I don't manually parse strings for every config value
- As a framework developer, I want `IsProduction()`, `IsDevelopment()`, `IsTesting()`, and `IsDebug()` helpers so that the framework can branch behavior by environment
- As a framework developer, I want `config.Load()` to be the very first call in `main()` so that all subsequent initialization has access to config values
- As a framework developer, I want graceful fallback when no `.env` file exists so that production deployments can inject env vars via Docker/Kubernetes

## Current State / Reference

### What Exists

- `.env` file with 10 config groups and placeholder values (Feature #01)
- `core/config/` directory with `.gitkeep` (Feature #01)
- Blueprint defines the config system approach (godotenv + Env helper)
- Framework reference doc (`docs/framework/core/configuration.md`) fully specifies the design

### What Works Well

- The `.env` file already contains all the config keys the framework will need
- Blueprint provides clear, simple code patterns — no over-engineering
- godotenv is a mature, minimal library with no transitive dependencies

### What Needs Improvement

N/A — greenfield code in `core/config/`.

## Proposed Approach

1. **Add godotenv dependency** — `go get github.com/joho/godotenv`
2. **Create `core/config/config.go`** — `Load()` function that calls `godotenv.Load()`
3. **Create `core/config/env.go`** — `Env()`, `EnvInt()`, `EnvBool()` typed accessor helpers
4. **Create `core/config/environment.go`** — `AppEnv()`, `IsProduction()`, `IsDevelopment()`, `IsTesting()`, `IsDebug()`
5. **Create `core/config/config_test.go`** — Unit tests for all helpers
6. **Update `cmd/main.go`** — Call `config.Load()` as the first line, use `Env()` to display app name and port in the banner

**Why godotenv and not Viper?** The blueprint recommends godotenv for `.env`-only projects and Viper for YAML/TOML/JSON. Our framework uses `.env` as the primary config mechanism. Viper can be added later if YAML config files are needed, but godotenv keeps the dependency tree minimal.

## Edge Cases & Risks

- [x] **No `.env` file** — `godotenv.Load()` returns an error but we log it and continue. System environment variables still work. This supports Docker/K8s deployments.
- [x] **Empty values vs missing keys** — `os.Getenv()` returns `""` for both missing and empty. The `Env()` fallback triggers on `""`, which means you can't set a value to empty string. This is an acceptable trade-off (same as Laravel's `env()` helper).
- [x] **`EnvInt` with non-numeric values** — Returns the fallback default and does NOT panic. Silent fallback is intentional for config robustness.
- [x] **`EnvBool` truthy values** — Only `"true"` and `"1"` are truthy. Everything else (including `"yes"`, `"on"`) returns false. Simple and predictable.
- [x] **Thread safety** — `os.Getenv()` is safe for concurrent reads. `godotenv.Load()` calls `os.Setenv()` which is also safe. No mutex needed.
- [x] **`.env` file encoding** — godotenv handles UTF-8. No special handling needed.

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Feature #01 — Project Setup & Structure | Feature | ✅ Done |
| `github.com/joho/godotenv` | External | ✅ Available (`go get`) |

## Open Questions

All resolved:

- [x] **godotenv or Viper?** → godotenv for now. Minimal, handles `.env` only. Viper deferred to future if needed.
- [x] **Should `Load()` panic on missing `.env`?** → No. Log and continue. Production uses system env vars.
- [x] **Should we support `.env.local` override loading?** → Not in this feature. `.env` is sufficient. `.env.local` can be added later if needed.
- [x] **Where do typed helpers live?** → Same package `core/config/`. Split across files for clarity: `config.go` (Load), `env.go` (helpers), `environment.go` (detection).

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-05 | Use godotenv, not Viper | Simpler, minimal deps, `.env`-only is sufficient |
| 2026-03-05 | Split into 3 files | `config.go`, `env.go`, `environment.go` — clear separation |
| 2026-03-05 | `EnvBool` only treats `"true"` and `"1"` as truthy | Simple, predictable, matches Go conventions |
| 2026-03-05 | `EnvInt` silently falls back on parse error | Config shouldn't panic — robustness over strictness |
| 2026-03-05 | Log missing `.env` as info, don't panic | Supports container deployments without `.env` files |

## Discussion Complete ✅

**Summary**: Feature #02 creates the `core/config` package with godotenv-based `.env` loading, typed accessors (`Env`, `EnvInt`, `EnvBool`), and environment detection helpers. First third-party dependency.
**Completed**: 2026-03-05
**Next**: Create architecture doc → `02-configuration-architecture.md`
