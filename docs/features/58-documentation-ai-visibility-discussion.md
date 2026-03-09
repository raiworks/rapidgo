# 💬 Discussion: Documentation & AI Visibility

> **Feature**: `58` — Documentation & AI Visibility
> **Status**: 🟡 IN PROGRESS
> **Branch**: `feature/58-documentation-ai-visibility`
> **Depends On**: All features #01–#57 (framework complete)
> **Date Started**: 2026-03-10
> **Date Completed**: —

---

## Summary

RapidGo has 56 shipped features covering everything from DI container to Prometheus metrics, queue workers, TOTP 2FA, OAuth2, GraphQL, plugin system, audit logging, and more. However, external reviewers (including AI agents like ChatGPT, Gemini, Claude) consistently miss these capabilities because the documentation doesn't make them visible at the right touchpoints.

This feature overhauls all documentation that AI agents and developers encounter when evaluating the framework: README.md, GitHub repo metadata, framework reference docs, and in-code documentation. The goal is that any reviewer — human or AI — immediately understands RapidGo is a full-featured Laravel-for-Go framework, not a simple HTTP router wrapper.

---

## The Problem

Three separate AI reviews (ChatGPT, Gemini) evaluated RapidGo and concluded:
- "Missing DI container" — **we have one** (Feature #05)
- "Missing event system" — **we have one** (Feature #34)
- "Missing queue workers" — **we have them** (Feature #42)
- "Missing CLI generators" — **we have them** (Feature #41)
- "Missing Prometheus metrics" — **we have it** (Feature #54)
- "Closer to a simple HTTP router" — **we're a full application framework**

The root cause is **documentation visibility**, not missing features. When an AI agent (or developer) checks the repo, they see:
1. **README.md** — Lists package names but doesn't describe capabilities or depth
2. **No feature highlights** at the repo root or in GitHub description
3. **No comparison table** showing RapidGo vs Gin/Echo/Fiber/Laravel
4. **Framework docs** exist in `docs/framework/` but are not linked prominently
5. **No FEATURES.md** or capabilities overview at the repo root
6. **Seeders doc** described an outdated pattern (now fixed)
7. **Several framework docs** still at "Draft" status

This matters because in 2026, the first thing people do when evaluating a framework is ask an AI agent about it. If the AI says "it's missing X" when X exists, adoption suffers.

---

## Functional Requirements

- As a developer evaluating RapidGo, I want to see ALL capabilities within 30 seconds of opening the repo
- As an AI agent analyzing the repository, I want structured metadata that clearly lists every feature category with status
- As a developer comparing frameworks, I want a comparison table showing RapidGo vs alternatives
- As a contributor, I want all framework docs to be "Final" status, not "Draft"
- As a user reading the README, I want to understand this is a Laravel-level framework, not a Gin wrapper

---

## Current State / Reference

### What Exists
- **README.md** — Package index table, hook system, quick start. Functional but sparse.
- **docs/framework/** — 59 RFC-style documents covering everything. Thorough but not discoverable from README.
- **docs/project-context.md** — Complete capabilities table. Internal doc, not visible to external reviewers.
- **GitHub repo description** — Unknown current state.

### What Works Well
- The `docs/framework/` documentation set is comprehensive and well-structured
- RFC-style format with frontmatter, terminology, security considerations
- Architecture diagrams exist for core systems
- Framework reference README has a complete documentation map

### What Needs Improvement
- README.md doesn't showcase the framework's depth
- No FEATURES.md at repo root
- No comparison with other Go frameworks
- No "Why RapidGo?" section
- Framework docs have "Draft" status markers that should be "Final"
- No badges (Go version, license, test status, etc.)
- Missing GitHub topics/tags for discoverability
- No structured AI-readable metadata (GitHub repo description, topics)

---

## Proposed Approach

### 1. README.md Overhaul
Transform from a package index into a comprehensive framework showcase:
- Hero section with tagline and badges
- "Why RapidGo?" section explaining the Laravel-for-Go positioning
- Feature highlights organized by category (not just package names)
- Comparison table: RapidGo vs Gin vs Echo vs Fiber vs Go Kit
- Architecture diagram (text-based)
- Quick start with meaningful example
- Links to all documentation sections

### 2. FEATURES.md
New file at repo root — exhaustive, structured list of every capability:
- Organized by category (Core, HTTP, Data, Security, Infrastructure, Advanced, CLI, Deployment)
- Each feature with description, package location, and "since version" marker
- Total feature count prominently displayed
- AI agents will find this file and use it as ground truth

### 3. Framework Docs Finalization
- Update all "Draft" status docs to "Final"
- Ensure all docs reference correct import paths (`github.com/RAiWorks/RapidGo/v2`)
- Verify code examples compile and match actual implementation
- Add cross-references between related docs

### 4. GitHub Metadata
- Repository description: comprehensive one-liner
- Topics: go, golang, web-framework, gin, gorm, laravel, rest-api, graphql, websocket, queue, scheduler, oauth2, totp, prometheus, etc.

---

## Edge Cases & Risks

- [ ] README.md getting too long — balance between comprehensive and readable
- [ ] Feature list becoming stale if new features are added without updating
- [ ] Code examples in docs that don't compile against v2 import paths
- [ ] Framework docs with stale references to pre-rename (RGo) paths

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| All features #01–#57 | Feature | ✅ Done |
| Seeders doc fix | Docs | ✅ Done (2026-03-10) |

---

## Open Questions

- [ ] Should FEATURES.md be auto-generated from code/docs or manually maintained?
- [ ] Should we add a CONTRIBUTING.md as part of this feature?
- [ ] Should we include performance benchmarks in the README or a separate doc?

---

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-10 | Fix seeders doc first | It was showing outdated function-based pattern instead of actual interface-based registry |
| 2026-03-10 | Two separate features (#58 docs, #59 website) | Docs overhaul is independent of building a website — different scope, different deliverables |
| 2026-03-10 | Target AI agents as primary audience for repo docs | In 2026, developers ask AI to evaluate frameworks before trying them |
