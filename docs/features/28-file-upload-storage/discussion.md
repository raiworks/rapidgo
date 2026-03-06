# Feature #28 — File Upload & Storage: Discussion

## What problem does this solve?

Applications need to store user-uploaded files (images, documents, assets) on disk or cloud storage. A driver-based storage abstraction lets the framework swap between local disk and S3-compatible services without changing application code.

## Why now?

Core infrastructure (#02, #05), helpers (#05), and response helpers are shipped. Storage is the next foundational service — mail (#29) and other features may depend on file references.

## What does the blueprint specify?

- **`Driver` interface** — `Put(path, io.Reader)`, `Get(path)`, `Delete(path)`, `URL(path)`.
- **`LocalDriver`** — stores files on local disk under a configurable base path.
- **`S3Driver`** — uses `aws-sdk-go-v2` for S3-compatible storage.
- **Env vars**: `STORAGE_DRIVER`, `STORAGE_LOCAL_PATH`, plus AWS vars for S3.
- **Upload controller** — example showing `FormFile` → `storageDriver.Put()`.

## Design decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Package location | `core/storage/` | Follows existing `core/*` pattern |
| Local driver | Implement fully | Stdlib only, no external deps |
| S3 driver | Defer to future | `aws-sdk-go-v2` is a massive dependency tree; not needed for core framework |
| Upload controller | Out of scope | App-level code, not core framework; blueprint shows it as a usage example |
| `NewDriver()` factory | Implement | Reads `STORAGE_DRIVER` env var, returns appropriate driver |
| Path traversal guard | Add to LocalDriver | Prevent `../../etc/passwd` attacks — not in blueprint but essential for security |

## What is out of scope?

- S3 driver (future — requires `aws-sdk-go-v2`).
- Upload controller (app-level, not framework).
- File size limits / MIME type validation (middleware concern, not storage).
- Image processing / thumbnails.
