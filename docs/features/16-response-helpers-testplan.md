# 🧪 Test Plan: Response Helpers

> **Feature**: `16` — Response Helpers
> **Architecture**: [`16-response-helpers-architecture.md`](16-response-helpers-architecture.md)
> **Status**: ⬜ NOT RUN
> **Result**: ⬜ NOT RUN

---

## Test File

`http/responses/response_test.go`

All tests use `httptest.NewRecorder()` with Gin's test mode — no server required.

---

## Test Cases

### TC-01: `TestSuccess_Returns200WithData`
**What**: `Success` returns 200 with `success: true` and data.
**How**: Create test context. Call `Success(c, map)`. Parse JSON. Assert status 200, `success == true`, data present.
**Pass**: Status 200, correct envelope.

### TC-02: `TestCreated_Returns201WithData`
**What**: `Created` returns 201 with `success: true` and data.
**How**: Create test context. Call `Created(c, map)`. Parse JSON. Assert status 201, `success == true`, data present.
**Pass**: Status 201, correct envelope.

### TC-03: `TestError_Returns404WithMessage`
**What**: `Error` returns the given status with `success: false` and error message.
**How**: Create test context. Call `Error(c, 404, "not found")`. Parse JSON. Assert status 404, `success == false`, `error == "not found"`.
**Pass**: Status 404, correct error envelope.

### TC-04: `TestError_Returns422WithMessage`
**What**: `Error` works with different status codes.
**How**: Call `Error(c, 422, "validation failed")`. Assert status 422.
**Pass**: Status 422, correct error envelope.

### TC-05: `TestPaginated_ReturnsDataWithMeta`
**What**: `Paginated` returns 200 with data and pagination meta.
**How**: Call `Paginated(c, items, 1, 10, 25)`. Parse JSON. Assert status 200, `meta.page == 1`, `meta.total_pages == 3`.
**Pass**: Meta fields correct, `totalPages = ceil(25/10) = 3`.

### TC-06: `TestPaginated_ExactDivision`
**What**: `Paginated` computes `totalPages` correctly when total is evenly divisible.
**How**: Call `Paginated(c, items, 1, 10, 30)`. Assert `meta.total_pages == 3`.
**Pass**: `totalPages = 30/10 = 3` (no remainder).

### TC-07: `TestSuccess_OmitsErrorField`
**What**: Success response does not include `error` field in JSON.
**How**: Call `Success(c, data)`. Check raw JSON body does not contain `"error"`.
**Pass**: No `error` key in JSON output.

### TC-08: `TestError_OmitsDataField`
**What**: Error response does not include `data` field in JSON.
**How**: Call `Error(c, 400, "bad")`. Check raw JSON body does not contain `"data"`.
**Pass**: No `data` key in JSON output.

---

## Test Summary

| ID | Test Name | Type | Scope |
|---|---|---|---|
| TC-01 | `TestSuccess_Returns200WithData` | Unit | Success helper |
| TC-02 | `TestCreated_Returns201WithData` | Unit | Created helper |
| TC-03 | `TestError_Returns404WithMessage` | Unit | Error helper (404) |
| TC-04 | `TestError_Returns422WithMessage` | Unit | Error helper (422) |
| TC-05 | `TestPaginated_ReturnsDataWithMeta` | Unit | Paginated with remainder |
| TC-06 | `TestPaginated_ExactDivision` | Unit | Paginated exact division |
| TC-07 | `TestSuccess_OmitsErrorField` | Unit | Omitempty behavior |
| TC-08 | `TestError_OmitsDataField` | Unit | Omitempty behavior |

**Total**: 8 test cases
**Expected new test count**: 164 + 8 = 172
