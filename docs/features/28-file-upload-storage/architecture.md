# Feature #28 — File Upload & Storage: Architecture

## Component overview

```
.env
    STORAGE_DRIVER=local
    STORAGE_LOCAL_PATH=storage/uploads
        │
        ▼
core/storage/storage.go      Driver interface + NewDriver() factory
core/storage/local.go        LocalDriver (Put/Get/Delete/URL)
        │
        ▼
app/providers/               (no provider needed — instantiated on demand)
```

## New files

| File | Purpose |
|------|---------|
| `core/storage/storage.go` | `Driver` interface, `NewDriver()` factory function |
| `core/storage/local.go` | `LocalDriver` struct implementing `Driver` for local disk |

## Driver interface

```go
type Driver interface {
    Put(path string, content io.Reader) (string, error)
    Get(path string) (io.ReadCloser, error)
    Delete(path string) error
    URL(path string) string
}
```

## LocalDriver

| Method | Behaviour |
|--------|-----------|
| `Put(path, content)` | Creates directories, writes file, returns stored path |
| `Get(path)` | Opens file, returns `io.ReadCloser` |
| `Delete(path)` | Removes file |
| `URL(path)` | Returns `BaseURL + "/" + path` |

### Security: Path traversal guard

All LocalDriver methods validate that the resolved path stays within `BasePath`. Any `../` traversal attempt returns an error. This is not in the blueprint but is critical for security.

## Environment variables

| Var | Default | Purpose |
|-----|---------|---------|
| `STORAGE_DRIVER` | `"local"` | Driver selection |
| `STORAGE_LOCAL_PATH` | `"storage/uploads"` | Base directory for local storage |

## Factory function

`NewDriver()` reads `STORAGE_DRIVER` from env:
- `"local"` (default) → returns `&LocalDriver{BasePath, BaseURL}`
- Unknown driver → returns error
