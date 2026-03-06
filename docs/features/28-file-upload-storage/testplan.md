# Feature #28 — File Upload & Storage: Test Plan

## Test cases

| TC | Description | Method | Expected |
|----|-------------|--------|----------|
| TC-01 | Put writes file to disk | Put | File exists, content matches |
| TC-02 | Put creates intermediate directories | Put | Nested path succeeds |
| TC-03 | Get returns file content | Get | Content matches what was Put |
| TC-04 | Get returns error for missing file | Get | Non-nil error |
| TC-05 | Delete removes file | Delete | File no longer exists |
| TC-06 | Delete returns error for missing file | Delete | Non-nil error |
| TC-07 | URL returns BaseURL + path | URL | Correct URL string |
| TC-08 | Path traversal in Put is rejected | Put("../evil") | Error returned |
| TC-09 | Path traversal in Get is rejected | Get("../evil") | Error returned |
| TC-10 | Path traversal in Delete is rejected | Delete("../evil") | Error returned |
| TC-11 | NewDriver returns LocalDriver by default | NewDriver() | Non-nil Driver |
| TC-12 | NewDriver returns error for unknown driver | Set env="unknown" | Error returned |

## Notes

- Tests use `t.TempDir()` for isolated temp directories.
- Path traversal tests verify security boundary enforcement.
- TC-11 and TC-12 use `t.Setenv()` for env manipulation.
