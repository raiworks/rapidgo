# Feature #40 — Testing Infrastructure: Plan

## Tasks

1. Create `testing/testutil/testutil.go` with utility functions.
2. Write self-tests.
3. Run full regression + go vet.
4. Commit, merge to main, push.

## Test Plan

| TC | Description | Expected |
|----|-------------|----------|
| 01 | NewTestRouter returns usable router | Non-nil, can register routes |
| 02 | NewTestDB returns working GORM DB | Can create tables, insert rows |
| 03 | DoRequest returns correct response | Status and body match handler |
| 04 | AssertStatus passes on match | No test failure |
| 05 | AssertJSONKey validates key/value | Correct key extracted |
