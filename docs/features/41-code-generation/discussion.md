# Feature #41 — Code Generation (CLI): Discussion

## Overview

`make:controller`, `make:model`, `make:service`, `make:provider` CLI commands that scaffold boilerplate Go files from templates, following the existing `make:migration` pattern.

## Blueprint Reference

Four generators:
- `make:controller [name]` → `http/controllers/<name>.go`
- `make:model [name]` → `database/models/<name>.go`
- `make:service [name]` → `app/services/<name>.go`
- `make:provider [name]` → `app/providers/<name>.go`

Shared `generate()` helper using `text/template`.

## Dependencies

- `#10` — CLI framework (Cobra, already integrated).
- Existing `make:migration` command establishes the pattern.
