# 💬 Discussion: Helpers

> **Feature**: `19` — Helpers
> **Depends on**: #01 (Project Setup)
> **Status**: 🟢 RESOLVED

---

## Context

The framework needs general-purpose utility functions for common web application tasks. The blueprint defines two sections — **Helpers** (password, random, string basics) and **Built-in String & Data Helpers** (extended string, number, time, data, config utilities). Both belong in `app/helpers/`.

## Blueprint Reference

**Helpers** (lines 2617–2674):
- `HashPassword`, `CheckPassword` (bcrypt)
- `RandomString` (crypto/rand hex)
- `Slugify`, `Truncate`

**Built-in String & Data Helpers** (lines 2805–2945):
- String: `Contains`, `Title`, `Excerpt`, `StripHTML`, `Mask`
- Number: `FormatBytes`, `Clamp`
- Time: `TimeAgo`, `FormatDate`
- Data: `StructToMap`, `MapKeys`
- Config: `Env`

## Scope Decision

Feature #19 implements **all functions from both blueprint sections** in `app/helpers/`, organized into logical files:

1. `app/helpers/password.go` — `HashPassword`, `CheckPassword`
2. `app/helpers/random.go` — `RandomString`
3. `app/helpers/string.go` — `Slugify`, `Truncate`, `Contains`, `Title`, `Excerpt`, `StripHTML`, `Mask`
4. `app/helpers/number.go` — `FormatBytes`, `Clamp`
5. `app/helpers/time.go` — `TimeAgo`, `FormatDate`
6. `app/helpers/data.go` — `StructToMap`, `MapKeys`
7. `app/helpers/env.go` — `Env`

**Out of scope**: `core/crypto/` utilities (belong to Feature #22 — Crypto & Security), Pagination helpers (separate feature).

## Key Decisions

| # | Decision | Rationale |
|---|---|---|
| 1 | Split into 7 files by domain | Better organization than one massive file; matches framework doc categories |
| 2 | `golang.org/x/crypto/bcrypt` for passwords | Blueprint specifies it; already indirect dependency |
| 3 | `strings.Title` deprecation | Blueprint uses `strings.Title`; we'll use `cases.Title` from `golang.org/x/text` or keep `strings.Title` with a note |
| 4 | All functions are pure and stateless | Blueprint guideline — easy to test |
| 5 | `StripHTML` is not an XSS sanitizer | Security note: use proper template escaping for XSS prevention |
