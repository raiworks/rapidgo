# 🧪 Test Plan: Helpers

> **Feature**: `19` — Helpers
> **Architecture**: [`19-helpers-architecture.md`](19-helpers-architecture.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Test File

`app/helpers/helpers_test.go` — package `helpers`

---

## Test Cases

### Password Helpers

| # | Test | Description | Assert |
|---|---|---|---|
| T01 | `TestHashPassword` | Hash a plain-text password | No error, hash ≠ plain |
| T02 | `TestCheckPassword_Valid` | Compare correct password | Returns `true` |
| T03 | `TestCheckPassword_Invalid` | Compare wrong password | Returns `false` |

### Random Helpers

| # | Test | Description | Assert |
|---|---|---|---|
| T04 | `TestRandomString_Length` | Generate random string of 16 bytes | Hex len = 32 |
| T05 | `TestRandomString_Unique` | Generate two strings | They differ |

### String Helpers

| # | Test | Description | Assert |
|---|---|---|---|
| T06 | `TestSlugify` | "Hello World!" → "hello-world" | Match expected |
| T07 | `TestTruncate_Short` | String shorter than max | Returns unchanged |
| T08 | `TestTruncate_Long` | String longer than max | Returns truncated + "..." |
| T09 | `TestContains` | Case-insensitive match | Returns `true` |
| T10 | `TestContains_NoMatch` | No match | Returns `false` |
| T11 | `TestTitle` | "hello world" → "Hello World" | Match expected |
| T12 | `TestExcerpt` | First N words | Truncated + "..." |
| T13 | `TestStripHTML` | Remove tags from HTML | Plain text |
| T14 | `TestMask` | "secret123" with first=2, last=2 → "se*****23" | Match expected |

### Number Helpers

| # | Test | Description | Assert |
|---|---|---|---|
| T15 | `TestFormatBytes_Zero` | 0 bytes | "0 B" |
| T16 | `TestFormatBytes_KB` | 1536 bytes | "1.50 KB" |
| T17 | `TestClamp_InRange` | Value within range | Returns value |
| T18 | `TestClamp_BelowMin` | Value below min | Returns min |
| T19 | `TestClamp_AboveMax` | Value above max | Returns max |

### Time Helpers

| # | Test | Description | Assert |
|---|---|---|---|
| T20 | `TestTimeAgo_JustNow` | Time < 1 minute ago | "just now" |
| T21 | `TestTimeAgo_MinutesAgo` | Time 5 minutes ago | "5 minutes ago" |
| T22 | `TestFormatDate` | Fixed time | Matches expected format |

### Data Helpers

| # | Test | Description | Assert |
|---|---|---|---|
| T23 | `TestStructToMap` | Convert struct to map | Keys/values match |
| T24 | `TestMapKeys` | Extract keys from map | All keys present |

### Env Helpers

| # | Test | Description | Assert |
|---|---|---|---|
| T25 | `TestEnv_Set` | Env var exists | Returns value |
| T26 | `TestEnv_Fallback` | Env var missing | Returns fallback |

---

## Coverage Summary

| Package | Functions | Tests | Target |
|---|---|---|---|
| `app/helpers` | 17 | 26 | 100% |

**Total tests after feature**: 185 (existing) + 26 = **211**
