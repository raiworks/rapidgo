# 📝 Changelog: Helpers

> **Feature**: `19` — Helpers
> **Date**: 2026-03-06

---

## Added

- `app/helpers/password.go` — `HashPassword`, `CheckPassword` (bcrypt)
- `app/helpers/random.go` — `RandomString` (crypto/rand hex)
- `app/helpers/string.go` — `Slugify`, `Truncate`, `Contains`, `Title`, `Excerpt`, `StripHTML`, `Mask`
- `app/helpers/number.go` — `FormatBytes`, `Clamp`
- `app/helpers/time.go` — `TimeAgo`, `FormatDate`
- `app/helpers/data.go` — `StructToMap`, `MapKeys`
- `app/helpers/env.go` — `Env`
- `app/helpers/helpers_test.go` — 26 tests covering all 17 functions

## Changed

- `go.mod` — `golang.org/x/crypto` promoted from indirect to direct dependency

## Deviations

_None documented yet. Will be updated during BUILD phase._
