# Feature #41 — Code Generation (CLI): Plan

## Tasks

1. Create `core/cli/make_scaffold.go` with `scaffold()` helper + 4 commands + templates.
2. Register commands in `root.go`.
3. Write tests: each command generates correct file with correct content.
4. Full regression + go vet.
5. Commit, merge, push.

## Test Plan

| TC | Description | Expected |
|----|-------------|----------|
| 01 | make:controller generates file | File exists with correct package + struct |
| 02 | make:model generates file | File exists with BaseModel embed |
| 03 | make:service generates file | File with constructor |
| 04 | make:provider generates file | File with Register + Boot |
| 05 | Duplicate name prevents overwrite | Returns error |
