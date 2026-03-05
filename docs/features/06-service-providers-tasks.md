# ✅ Tasks: Service Providers

> **Feature**: `06` — Service Providers
> **Architecture**: [`06-service-providers-architecture.md`](06-service-providers-architecture.md)
> **Branch**: `feature/06-service-providers`
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/14 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [ ] Feature branch created from latest `main`
- [x] Dependent features are merged to `main`
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — Config Provider

> ConfigProvider that loads .env via config.Load().

- [ ] **A.1** — Create `app/providers/config_provider.go` with `ConfigProvider` struct
- [ ] **A.2** — Implement `Register()` — calls `config.Load()`
- [ ] **A.3** — Implement `Boot()` — no-op
- [ ] 📍 **Checkpoint A** — ConfigProvider compiles, `go vet` clean

---

## Phase B — Logger Provider

> LoggerProvider that sets up slog in Boot phase.

- [ ] **B.1** — Create `app/providers/logger_provider.go` with `LoggerProvider` struct
- [ ] **B.2** — Implement `Register()` — no-op
- [ ] **B.3** — Implement `Boot()` — calls `logger.Setup()`
- [ ] 📍 **Checkpoint B** — LoggerProvider compiles, `go vet` clean

---

## Phase C — Update main.go

> Switch from direct calls to App bootstrap pattern.

- [ ] **C.1** — Update `cmd/main.go` to use `app.New()`, register providers, call `Boot()`
- [ ] **C.2** — Remove direct `config.Load()` and `logger.Setup()` calls
- [ ] **C.3** — Verify `go run cmd/main.go` produces identical output
- [ ] 📍 **Checkpoint C** — App boots correctly, banner and log output unchanged

---

## Phase D — Testing

> Test providers and bootstrap lifecycle.

- [ ] **D.1** — Create `app/providers/providers_test.go` with test cases
- [ ] **D.2** — Run `go test ./app/providers/...` — all tests pass
- [ ] **D.3** — Run `go vet ./...` — no issues
- [ ] 📍 **Checkpoint D** — All tests pass, zero vet warnings

---

## Phase E — Documentation & Cleanup

> Changelog, roadmap, self-review.

- [ ] **E.1** — Update changelog doc with implementation summary
- [ ] **E.2** — Self-review all diffs — code is clean, idiomatic Go
- [ ] 📍 **Checkpoint E** — Clean code, complete docs, ready to ship

---

## Ship 🚀

- [ ] All phases complete
- [ ] All checkpoints verified
- [ ] Final commit with descriptive message
- [ ] Merge to `main`
- [ ] Push `main`
- [ ] **Keep the feature branch** — do not delete
- [ ] Update project roadmap progress
- [ ] Create review doc → `06-service-providers-review.md`
