# 🧪 Test Plan: Services Layer

> **Feature**: `18` — Services Layer
> **Architecture**: [`18-services-layer-architecture.md`](18-services-layer-architecture.md)
> **Status**: ⬜ NOT RUN
> **Result**: ⬜ NOT RUN

---

## Test File

`app/services/user_service_test.go`

All tests use SQLite in-memory with `AutoMigrate` for the `User` model.

---

## Test Cases

### TC-01: `TestNewUserService_ReturnsService`
**What**: Constructor returns a valid `UserService` with the DB set.
**How**: Create in-memory DB. Call `NewUserService(db)`. Assert non-nil, DB field set.
**Pass**: Service created with DB.

### TC-02: `TestCreate_ReturnsUser`
**What**: `Create` inserts a new user and returns it.
**How**: Call `svc.Create("Alice", "alice@example.com", "pass")`. Assert returned user has ID > 0, correct fields.
**Pass**: User persisted with correct data.

### TC-03: `TestCreate_DuplicateEmail_ReturnsError`
**What**: `Create` rejects duplicate email.
**How**: Create user with email. Call `Create` again with same email. Assert error contains "email already exists".
**Pass**: Error returned, no duplicate inserted.

### TC-04: `TestGetByID_ReturnsUser`
**What**: `GetByID` retrieves an existing user.
**How**: Create user. Call `svc.GetByID(user.ID)`. Assert fields match.
**Pass**: Correct user returned.

### TC-05: `TestGetByID_NotFound_ReturnsError`
**What**: `GetByID` returns error for non-existent ID.
**How**: Call `svc.GetByID(9999)`. Assert error is not nil.
**Pass**: Error returned.

### TC-06: `TestUpdate_UpdatesFields`
**What**: `Update` modifies specified fields.
**How**: Create user. Call `svc.Update(id, map{"name": "Bob"})`. Assert returned user has name "Bob".
**Pass**: User updated in database.

### TC-07: `TestUpdate_NotFound_ReturnsError`
**What**: `Update` returns error for non-existent ID.
**How**: Call `svc.Update(9999, map{"name": "X"})`. Assert error.
**Pass**: Error returned.

### TC-08: `TestDelete_RemovesUser`
**What**: `Delete` removes a user from the database.
**How**: Create user. Call `svc.Delete(id)`. Call `svc.GetByID(id)`. Assert GetByID returns error.
**Pass**: User no longer exists.

---

## Test Summary

| ID | Test Name | Type | Scope |
|---|---|---|---|
| TC-01 | `TestNewUserService_ReturnsService` | Unit | Constructor |
| TC-02 | `TestCreate_ReturnsUser` | Integration | Create user |
| TC-03 | `TestCreate_DuplicateEmail_ReturnsError` | Integration | Duplicate email |
| TC-04 | `TestGetByID_ReturnsUser` | Integration | Get by ID |
| TC-05 | `TestGetByID_NotFound_ReturnsError` | Integration | Not found |
| TC-06 | `TestUpdate_UpdatesFields` | Integration | Update fields |
| TC-07 | `TestUpdate_NotFound_ReturnsError` | Integration | Update not found |
| TC-08 | `TestDelete_RemovesUser` | Integration | Delete user |

**Total**: 8 test cases
**Expected new test count**: 177 + 8 = 185
