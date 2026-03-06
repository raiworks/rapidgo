# Feature #41 — Code Generation (CLI): Design

## File

`core/cli/make_scaffold.go`

## Commands

| Command | Output Path | Template |
|---------|-------------|----------|
| `make:controller [Name]` | `http/controllers/<name>.go` | Controller with Index, Show, Store, Update, Destroy |
| `make:model [Name]` | `database/models/<name>.go` | Model struct with BaseModel embed |
| `make:service [Name]` | `app/services/<name>.go` | Service struct with DB field + constructor |
| `make:provider [Name]` | `app/providers/<name>.go` | Provider with Register + Boot stubs |

## Shared Helper

```go
func scaffold(kind, name, tpl, dir string, out io.Writer) error
```

- Creates directory if needed (`os.MkdirAll`).
- Uses `os.O_CREATE|os.O_EXCL` to prevent overwriting.
- Executes template with `{{.Name}}`.
- Prints confirmation message.

## Registration

All four commands added to `rootCmd` in `init()` in `root.go`.

## File Naming

`toSnakeCase(name) + ".go"` — same conversion used by `make:migration`.
