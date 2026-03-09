# 🧪 Test Plan: Documentation & AI Visibility

> **Feature**: `58` — Documentation & AI Visibility
> **Tasks**: [`58-documentation-ai-visibility-tasks.md`](58-documentation-ai-visibility-tasks.md)
> **Date**: 2026-03-10

---

## Test Strategy

This is a documentation-only feature. Tests are verification checks, not code tests.

---

## Test Cases

| # | Test | Method | Expected Result |
|---|------|--------|-----------------|
| TC-01 | README contains all 8 feature categories | Manual review | All 8 categories present with checkmarks |
| TC-02 | README mentions all 56 features | Count checkmarks | ≥56 feature items listed |
| TC-03 | README has comparison table | Visual check | Table with RapidGo vs Gin/Echo/Fiber/Go Kit |
| TC-04 | README has architecture diagram | Visual check | Text-based request lifecycle diagram |
| TC-05 | FEATURES.md exists at repo root | File check | File exists, non-empty |
| TC-06 | FEATURES.md lists all 56 features | Count entries | 56 features with package paths |
| TC-07 | FEATURES.md package paths are valid | Cross-check with codebase | Every listed package exists |
| TC-08 | COMPARISON.md exists at repo root | File check | File exists, non-empty |
| TC-09 | COMPARISON.md has 30+ row matrix | Count rows | ≥30 feature comparison rows |
| TC-10 | No framework docs have "Draft" status | Grep for `status: "Draft"` | 0 matches in `docs/framework/` |
| TC-11 | All framework docs use v2 import paths | Grep for old paths | 0 matches for non-v2 paths |
| TC-12 | README links all work | Click test | No broken links |
| TC-13 | AI readability check | Ask AI to summarize repo | AI identifies RapidGo as full-featured framework |

---

## Acceptance Criteria

1. ✅ README.md showcases all 56 features in organized categories
2. ✅ FEATURES.md provides exhaustive detail with verifiable package paths
3. ✅ COMPARISON.md definitively answers "why not just use Gin?"
4. ✅ All framework docs at "Final" status
5. ✅ No stale import paths or code examples
