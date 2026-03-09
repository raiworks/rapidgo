# 🏗️ Architecture: Documentation & AI Visibility

> **Feature**: `58` — Documentation & AI Visibility
> **Discussion**: [`58-documentation-ai-visibility-discussion.md`](58-documentation-ai-visibility-discussion.md)
> **Status**: 🟡 DRAFT
> **Date**: 2026-03-10

---

## Overview

This feature creates and improves documentation files so that AI agents and developers immediately recognize RapidGo as a comprehensive, Laravel-level Go framework. No code changes — only documentation files.

---

## File Structure

```
(root)/
├── README.md                          # MODIFY — full overhaul
├── FEATURES.md                        # CREATE — exhaustive capabilities list
├── COMPARISON.md                      # CREATE — framework comparison matrix
│
docs/framework/
├── (all existing .md files)           # MODIFY — update status Draft → Final
│                                      #          verify import paths, examples
│
docs/features/
├── 58-documentation-ai-visibility-*   # Feature docs (this set)
```

---

## Component Design

### 1. README.md — Full Overhaul

**Current**: 100 lines. Package table + hooks + quick start.
**Target**: ~300 lines. Full framework showcase.

#### Structure

```
# RapidGo
## Badges (Go version, license, tests, Go Report Card)
## One-liner tagline
## Why RapidGo?
  - What it is (Laravel-for-Go)
  - What it is NOT (not another HTTP router)
  - Built on (Gin + GORM + Cobra)
## Feature Highlights
  ### Core Application
  ### HTTP & Routing
  ### Data & Database
  ### Security & Auth
  ### Infrastructure
  ### Advanced Features
  ### CLI & Developer Experience
  ### Deployment
## Architecture Overview (text diagram)
## Quick Start
## Comparison with Other Frameworks (summary table, link to COMPARISON.md)
## Package Index (existing table, improved)
## Hook System (existing section)
## Documentation Links
## Contributing
## License
```

#### Feature Highlights Format

Each category shows capabilities with checkmarks to make scanning easy:

```markdown
### Security & Authentication
- ✅ JWT authentication (stateless)
- ✅ Session-based authentication (5 backends: database, Redis, file, memory, cookie)
- ✅ OAuth2 / social login (Google, GitHub, Facebook, etc.)
- ✅ TOTP two-factor authentication with backup codes
- ✅ CSRF protection (double-submit cookie pattern)
- ✅ CORS configuration (per-origin, per-method)
- ✅ Rate limiting (token bucket, per-IP and per-route)
- ✅ AES-256-GCM encryption, bcrypt hashing, HMAC, secure tokens
- ✅ Audit logging (who did what, when)
```

### 2. FEATURES.md — Exhaustive Capabilities

**Purpose**: Single document that AI agents will find and use as ground truth.

#### Structure

```markdown
# RapidGo — Complete Feature List

> 56 features shipped across 6 phases. Built on Go 1.25+, Gin, GORM, Cobra.

## Feature Count by Category
| Category | Count |
|----------|-------|
| Core | 10 |
| MVC + Auth | 12 |
| ... | ... |
| **Total** | **56** |

## Core Application (Phase 1)
### Service Container & Dependency Injection
- Package: `core/container`
- Singleton and transient bindings
- Type-safe resolution with `MustMake[T]()`
- Provider pattern: Register() + Boot() lifecycle
...

## [Each category with full detail]

## Technology Stack
[Full dependency table with versions]

## Architecture
[Request lifecycle diagram]
```

Each feature entry includes:
- **Package path** (so verification is easy)
- **Key capabilities** (bullet points)
- **Phase number** (shows maturity)

### 3. COMPARISON.md — Framework Comparison

**Purpose**: Pre-answer the "should I use Gin/Echo/Fiber instead?" question.

#### Content

```markdown
# RapidGo vs Other Go Frameworks

## Framework Categories
1. HTTP Routers: Gin, Echo, Fiber
2. Application Frameworks: RapidGo
3. Microservice Platforms: Go Kit, Kratos, Dapr

## Detailed Comparison Matrix
[30+ row comparison table]

## When to Use RapidGo
## When NOT to Use RapidGo
## How RapidGo Uses Gin
```

### 4. Framework Docs Finalization

Update all `docs/framework/**/*.md` files:
- Change YAML frontmatter `status: "Draft"` → `status: "Final"`
- Update `last_updated` to `2026-03-10`
- Verify import paths use `github.com/RAiWorks/RapidGo/v2`
- Spot-check code examples against actual implementations

---

## Data Flow

```
Developer/AI discovers repo
    → Reads README.md (hero → features → comparison → quick start)
    → Finds FEATURES.md (exhaustive detail)
    → Finds COMPARISON.md (positioning)
    → Follows links to docs/framework/ (deep reference)
    → Conclusion: "This is a complete framework"
```

---

## Trade-offs

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| Single mega-README | Everything in one file | Too long, hard to maintain | **No** — split across README + FEATURES + COMPARISON |
| Auto-generate from code | Always in sync | Complex tooling, loses narrative | **No** — manually curated |
| Keep current README | Less work | AI agents keep missing features | **No** — this is the whole point |

---

## Security Considerations

- No code changes, documentation only
- Ensure no real credentials appear in example code snippets
- Ensure `.env` examples use placeholder values

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-10 | RAiWorks | Initial architecture |
