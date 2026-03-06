# ✅ Tasks: Response Helpers

> **Feature**: `16` — Response Helpers
> **Architecture**: [`16-response-helpers-architecture.md`](16-response-helpers-architecture.md)
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/3 phases complete

---

## Phase A — Response Helpers

> Create `http/responses/response.go` with types and helpers.

- [ ] **A.1** — Create `http/responses/response.go`: `APIResponse`, `Meta`, `Success`, `Created`, `Error`, `Paginated`
- [ ] **A.2** — `go build ./http/responses/...` clean
- [ ] **A.3** — `go vet ./http/responses/...` clean
- [ ] 📍 **Checkpoint A** — Response helpers compile

## Phase B — Testing

> Create tests for all response helpers.

- [ ] **B.1** — Create `http/responses/response_test.go` with test cases from testplan
- [ ] **B.2** — `go test ./http/responses/... -v` — all pass
- [ ] **B.3** — `go test ./... -count=1` — full regression, 0 failures
- [ ] 📍 **Checkpoint B** — All new tests pass, no regressions

## Phase C — Changelog & Self-Review

- [ ] **C.1** — Update `16-response-helpers-changelog.md` with build log and deviations
- [ ] **C.2** — Cross-check: verify code matches architecture doc
- [ ] 📍 **Checkpoint C** — Changelog complete, architecture consistent
