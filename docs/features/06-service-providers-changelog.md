# 📝 Changelog: Service Providers

> **Feature**: `06` — Service Providers
> **Branch**: `feature/06-service-providers`
> **Started**: 2026-03-06
> **Completed**: 2026-03-06

---

## Log

- **2026-03-06** — Created `app/providers/config_provider.go`: `ConfigProvider` with `Register()` → `config.Load()`
- **2026-03-06** — Created `app/providers/logger_provider.go`: `LoggerProvider` with `Boot()` → `logger.Setup()`
- **2026-03-06** — Updated `cmd/main.go`: App bootstrap pattern (`app.New()` → Register → Boot)
- **2026-03-06** — Created `app/providers/providers_test.go`: 8 tests (6 runtime + 2 compile-time interface checks)
- **2026-03-06** — All tests pass across entire project, `go vet` clean
- **2026-03-06** — `go run cmd/main.go` output identical to before — banner + slog.Info unchanged

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| Tests use `t.Setenv` instead of `.env` file | Tests relied on `.env` being present | Tests set env vars via `t.Setenv` for isolation | Test runner CWD is `app/providers/`, not project root — `.env` not found. `t.Setenv` is idiomatic Go test pattern. |

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
| Use `t.Setenv` for test env vars | Config tests in `core/config/` use same pattern — consistent approach | 2026-03-06 |
