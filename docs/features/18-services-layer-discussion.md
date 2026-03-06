# 💬 Discussion: Services Layer

> **Feature**: `18` — Services Layer
> **Depends on**: #05 (Service Container), #11 (Models)
> **Status**: 🟢 RESOLVED

---

## Context

Features #05 and #11 built the service container and model layer. Controllers (#15) handle HTTP concerns and return responses (#16). The framework needs a services layer to contain business logic, sitting between controllers and models in the delegation chain: Controller → Service → Model.

## Blueprint Reference

The blueprint's **Services Layer** section (lines 2528–2600) prescribes:
- Services in `app/services/`
- `UserService` struct with `DB *gorm.DB`
- `NewUserService(db)` constructor
- CRUD methods: `GetByID`, `Create`, `Update`, `Delete`
- Controller integration example using `responses.Error` and `responses.Created`

The framework doc `docs/framework/infrastructure/services-layer.md` mirrors this exactly and adds a `make:service` CLI command note (out of scope — belongs to CLI scaffolding, a later feature).

## Scope Decision

Feature #18 implements the **UserService example** as the canonical service pattern. Specifically:

1. **UserService** in `app/services/user_service.go` — `GetByID`, `Create`, `Update`, `Delete` methods
2. **Tests** in `app/services/user_service_test.go` — unit tests using SQLite in-memory

**Out of scope**: Controller integration (would create import cycles between `controllers` and `services`), `make:service` CLI command (later feature), password hashing (belongs to Crypto & Security #22), request validation (belongs to Validation, a later feature).

## Key Decisions

| # | Decision | Rationale |
|---|---|---|
| 1 | `UserService` as example pattern | Blueprint showcases it explicitly; users follow the pattern for other services |
| 2 | `Password` stored as-is (with comment) | Blueprint code has `// hash before saving in real code`; actual hashing belongs to #22 |
| 3 | Use `github.com/RAiWorks/RGo/database/models` import | Correct module path for the framework's own models package |
| 4 | Tests use SQLite in-memory | Consistent with database test patterns from #09, #11, #14 |
| 5 | No container registration | Services are instantiated directly with `NewUserService(db)` — blueprint doesn't register them in the container |
