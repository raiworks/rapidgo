# 🧪 Test Plan: Controllers

> **Feature**: `15` — Controllers
> **Architecture**: [`15-controllers-architecture.md`](15-controllers-architecture.md)
> **Status**: ⬜ NOT RUN
> **Result**: ⬜ NOT RUN

---

## Test File

`http/controllers/controllers_test.go`

All tests use `httptest.NewRecorder()` with Gin's test mode — no server required.

---

## Test Cases

### TC-01: `TestHome_ReturnsWelcome`
**What**: `Home` returns a JSON welcome message with status 200.
**How**: Create a Gin test context with `GET /`. Call `Home(c)`. Assert status 200 and response body contains `"Welcome to RGo"`.
**Pass**: Status 200, JSON body with `message` field.

### TC-02: `TestPostController_Index`
**What**: `PostController.Index` returns status 200 with index message.
**How**: Create a Gin test context with `GET /posts`. Call `Index(c)`. Assert status 200.
**Pass**: Status 200, JSON body with `"PostController index"`.

### TC-03: `TestPostController_Store`
**What**: `PostController.Store` returns status 201.
**How**: Create a Gin test context with `POST /posts`. Call `Store(c)`. Assert status 201.
**Pass**: Status 201, JSON body with `"PostController store"`.

### TC-04: `TestPostController_Show`
**What**: `PostController.Show` returns the id from URL param.
**How**: Create a Gin test context with `GET /posts/42`, set param `id=42`. Call `Show(c)`. Assert status 200 and response contains `"id": "42"`.
**Pass**: Status 200, JSON body includes the id.

### TC-05: `TestPostController_Update`
**What**: `PostController.Update` returns status 200 with id.
**How**: Create a Gin test context with `PUT /posts/7`, set param `id=7`. Call `Update(c)`. Assert status 200 and response contains `"id": "7"`.
**Pass**: Status 200, JSON body includes the id.

### TC-06: `TestPostController_Destroy`
**What**: `PostController.Destroy` returns status 200.
**How**: Create a Gin test context with `DELETE /posts/1`. Call `Destroy(c)`. Assert status 200.
**Pass**: Status 200, JSON body with `"PostController destroy"`.

### TC-07: `TestPostController_ImplementsResourceController`
**What**: `PostController` satisfies the `ResourceController` interface at compile time.
**How**: Compile-time assertion: `var _ router.ResourceController = (*PostController)(nil)`.
**Pass**: Compiles without error.

### TC-08: `TestRoutes_HomeRegistered`
**What**: `GET /` is registered and returns 200.
**How**: Create a Router, call `RegisterWeb(r)`, send `GET /` via `httptest`. Assert status 200.
**Pass**: Status 200 response.

### TC-09: `TestRoutes_APIPostsRegistered`
**What**: `GET /api/posts` is registered and returns 200.
**How**: Create a Router, call `RegisterAPI(r)`, send `GET /api/posts` via `httptest`. Assert status 200.
**Pass**: Status 200 response.

---

## Test Summary

| ID | Test Name | Type | Scope |
|---|---|---|---|
| TC-01 | `TestHome_ReturnsWelcome` | Unit | Home controller |
| TC-02 | `TestPostController_Index` | Unit | Index action |
| TC-03 | `TestPostController_Store` | Unit | Store action (201) |
| TC-04 | `TestPostController_Show` | Unit | Show with param |
| TC-05 | `TestPostController_Update` | Unit | Update with param |
| TC-06 | `TestPostController_Destroy` | Unit | Destroy action |
| TC-07 | `TestPostController_ImplementsResourceController` | Compile | Interface compliance |
| TC-08 | `TestRoutes_HomeRegistered` | Integration | Web routes |
| TC-09 | `TestRoutes_APIPostsRegistered` | Integration | API routes |

**Total**: 9 test cases
**Expected new test count**: 155 + 9 = 164
