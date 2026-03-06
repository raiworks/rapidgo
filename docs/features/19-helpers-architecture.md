# 🏗️ Architecture: Helpers

> **Feature**: `19` — Helpers
> **Discussion**: [`19-helpers-discussion.md`](19-helpers-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

Feature #19 adds 17 general-purpose helper functions to `app/helpers/`, covering passwords, random strings, string manipulation, numbers, time formatting, data conversion, and environment access. All functions are pure, stateless, and well-tested.

## File Structure

```
app/helpers/
├── password.go     # HashPassword, CheckPassword
├── random.go       # RandomString
├── string.go       # Slugify, Truncate, Contains, Title, Excerpt, StripHTML, Mask
├── number.go       # FormatBytes, Clamp
├── time.go         # TimeAgo, FormatDate
├── data.go         # StructToMap, MapKeys
├── env.go          # Env
└── helpers_test.go # All tests
```

### Files Created (7 + 1 test)
| File | Package | Lines (est.) |
|---|---|---|
| `app/helpers/password.go` | `helpers` | ~20 |
| `app/helpers/random.go` | `helpers` | ~15 |
| `app/helpers/string.go` | `helpers` | ~60 |
| `app/helpers/number.go` | `helpers` | ~30 |
| `app/helpers/time.go` | `helpers` | ~30 |
| `app/helpers/data.go` | `helpers` | ~25 |
| `app/helpers/env.go` | `helpers` | ~12 |

### Files Modified (0)
No existing files need modification.

---

## Component Design

### Password Helpers (`app/helpers/password.go`)

```go
package helpers

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain-text password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a hashed password with a plain-text candidate.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
```

### Random Helpers (`app/helpers/random.go`)

```go
package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomString generates a cryptographically random hex string of n bytes.
func RandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
```

### String Helpers (`app/helpers/string.go`)

```go
package helpers

import (
	"html"
	"regexp"
	"strings"
	"unicode"
)

// Slugify converts a string to a URL-friendly slug.
func Slugify(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

// Truncate shortens a string to max length, appending "..." if truncated.
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// Contains checks if needle exists in haystack (case-insensitive).
func Contains(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}

// Title converts a string to title case: "hello world" → "Hello World".
func Title(s string) string {
	return strings.Title(strings.ToLower(s))
}

// Excerpt returns the first n words from a string.
func Excerpt(s string, words int) string {
	parts := strings.Fields(s)
	if len(parts) <= words {
		return s
	}
	return strings.Join(parts[:words], " ") + "..."
}

// StripHTML removes all HTML tags from a string.
func StripHTML(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return html.UnescapeString(re.ReplaceAllString(s, ""))
}

// Mask masks a string showing only first/last n chars: "secret123" → "se*****23".
func Mask(s string, showFirst, showLast int) string {
	if len(s) <= showFirst+showLast {
		return s
	}
	masked := s[:showFirst] + strings.Repeat("*", len(s)-showFirst-showLast) + s[len(s)-showLast:]
	return masked
}
```

### Number Helpers (`app/helpers/number.go`)

```go
package helpers

import (
	"fmt"
	"math"
)

// FormatBytes converts bytes to human-readable: 1536 → "1.50 KB".
func FormatBytes(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	if bytes == 0 {
		return "0 B"
	}
	i := int(math.Log(float64(bytes)) / math.Log(1024))
	if i >= len(units) {
		i = len(units) - 1
	}
	return fmt.Sprintf("%.2f %s", float64(bytes)/math.Pow(1024, float64(i)), units[i])
}

// Clamp restricts a value between min and max.
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
```

### Time Helpers (`app/helpers/time.go`)

```go
package helpers

import (
	"fmt"
	"time"
)

// TimeAgo returns a human-readable relative time: "5 minutes ago".
func TimeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	default:
		return t.Format("Jan 2, 2006")
	}
}

// FormatDate formats time as "Jan 2, 2006 3:04 PM".
func FormatDate(t time.Time) string {
	return t.Format("Jan 2, 2006 3:04 PM")
}
```

### Data Helpers (`app/helpers/data.go`)

```go
package helpers

import "encoding/json"

// StructToMap converts a struct to map[string]interface{} via JSON.
func StructToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	return m, err
}

// MapKeys returns the keys of a map.
func MapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
```

### Config Helpers (`app/helpers/env.go`)

```go
package helpers

import "os"

// Env reads an env variable with a fallback default.
func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
```

---

## Dependencies

| Dependency | Type | Usage |
|---|---|---|
| `golang.org/x/crypto/bcrypt` | external | Password hashing |
| `crypto/rand` | stdlib | Random string generation |
| `encoding/hex` | stdlib | Hex encoding |
| `html` | stdlib | HTML unescaping |
| `regexp` | stdlib | HTML tag stripping |
| `encoding/json` | stdlib | Struct to map conversion |
| `math` | stdlib | FormatBytes logarithm |
| `os` | stdlib | Env variable access |

---

## Impact on Existing Code

| Component | Impact |
|---|---|
| `app/helpers/` | `.gitkeep` stays; 7 new files added |
| `go.mod` | `golang.org/x/crypto` becomes direct dependency |
| No other files modified | Helpers are standalone utilities |
