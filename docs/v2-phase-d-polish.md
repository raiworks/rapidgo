# RapidGo v2 — Phase D: Polish

> **Phase**: D — Polish and Release  
> **Steps**: D1 (rapidgo new command) + D2 (documentation)  
> **Branches**: `feature/v2-09-rapidgo-new-cmd`, `feature/v2-10-library-readme`  
> **Pre-requisite**: Phase C complete (library + starter repos working independently)  
> **Post-condition**: v2.0.0 tagged and released  

---

## Step D1: `rapidgo new` CLI Command

### Branch

`feature/v2-09-rapidgo-new-cmd` (from `v2`)

### Objective

Add a `rapidgo new myapp` command that scaffolds a new project from the starter template. This downloads the starter repo, replaces the module name, and runs `go mod tidy`.

### Files Changed

| Action | File |
|--------|------|
| CREATE | `core/cli/new.go` |
| CREATE | `core/cli/new_test.go` |
| MODIFY | `core/cli/root.go` | Add `newCmd` to root command |

### Implementation: `core/cli/new.go`

```go
package cli

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const starterRepo = "https://github.com/RAiWorks/RapidGo-starter/archive/refs/heads/main.zip"

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new RapidGo project from the starter template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Validate project name
		if strings.ContainsAny(name, `/\:*?"<>|`) {
			return fmt.Errorf("invalid project name: %s", name)
		}

		// Check target directory doesn't exist
		if _, err := os.Stat(name); err == nil {
			return fmt.Errorf("directory %q already exists", name)
		}

		fmt.Printf("Creating new RapidGo project: %s\n", name)

		// 1. Download starter template
		fmt.Println("  Downloading starter template...")
		zipPath, err := downloadStarter()
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
		defer os.Remove(zipPath)

		// 2. Extract to project directory
		fmt.Println("  Extracting template...")
		if err := extractZip(zipPath, name); err != nil {
			return fmt.Errorf("extract failed: %w", err)
		}

		// 3. Replace module name
		fmt.Println("  Configuring module name...")
		if err := replaceModuleName(name); err != nil {
			return fmt.Errorf("module rename failed: %w", err)
		}

		// 4. Run go mod tidy
		fmt.Println("  Running go mod tidy...")
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = name
		tidyCmd.Stdout = os.Stdout
		tidyCmd.Stderr = os.Stderr
		if err := tidyCmd.Run(); err != nil {
			return fmt.Errorf("go mod tidy failed: %w", err)
		}

		fmt.Println("=================================")
		fmt.Printf("  Project %s created!\n", name)
		fmt.Println("=================================")
		fmt.Printf("  cd %s\n", name)
		fmt.Println("  cp .env.example .env")
		fmt.Println("  go run cmd/main.go serve")
		fmt.Println("=================================")

		return nil
	},
}

// downloadStarter downloads the starter template zip from GitHub.
func downloadStarter() (string, error) {
	resp, err := http.Get(starterRepo)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d from GitHub", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "rapidgo-starter-*.zip")
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", err
	}
	tmpFile.Close()
	return tmpFile.Name(), nil
}

// extractZip extracts a GitHub archive zip to the target directory.
// GitHub archives contain a top-level directory (e.g., "RapidGo-starter-main/")
// which is stripped during extraction.
func extractZip(zipPath, targetDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Find the common prefix (GitHub adds "RepoName-branch/")
	prefix := ""
	if len(r.File) > 0 {
		prefix = strings.SplitN(r.File[0].Name, "/", 2)[0] + "/"
	}

	for _, f := range r.File {
		// Strip the GitHub prefix directory
		relPath := strings.TrimPrefix(f.Name, prefix)
		if relPath == "" {
			continue
		}

		destPath := filepath.Join(targetDir, relPath)

		// Ensure path doesn't escape target directory (zip slip protection)
		if !strings.HasPrefix(filepath.Clean(destPath), filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in zip: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(destPath, 0755)
			continue
		}

		os.MkdirAll(filepath.Dir(destPath), 0755)
		outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// replaceModuleName replaces "github.com/RAiWorks/RapidGo-starter" with
// the project name in go.mod and all .go files.
func replaceModuleName(projectDir string) error {
	oldModule := "github.com/RAiWorks/RapidGo-starter"
	newModule := projectDir // simple name like "myapp"

	return filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		name := filepath.Base(path)

		// Only process .go files and go.mod
		if ext != ".go" && name != "go.mod" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		newContent := strings.ReplaceAll(string(content), oldModule, newModule)
		if newContent != string(content) {
			return os.WriteFile(path, []byte(newContent), info.Mode())
		}

		return nil
	})
}
```

### Modification: `core/cli/root.go`

Add `newCmd` to the init block:

```go
func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(migrateRollbackCmd)
	rootCmd.AddCommand(migrateStatusCmd)
	rootCmd.AddCommand(makeMigrationCmd)
	rootCmd.AddCommand(dbSeedCmd)
	rootCmd.AddCommand(makeControllerCmd)
	rootCmd.AddCommand(makeModelCmd)
	rootCmd.AddCommand(makeServiceCmd)
	rootCmd.AddCommand(makeProviderCmd)
	rootCmd.AddCommand(workCmd)
	rootCmd.AddCommand(scheduleRunCmd)
	rootCmd.AddCommand(makeAdminCmd)
	rootCmd.AddCommand(newCmd)               // ← ADD
}
```

### Test Plan: D1

| # | Test | Verifies |
|---|------|----------|
| T01 | `rapidgo new` with no args | Shows error + usage |
| T02 | `rapidgo new myapp` | Downloads, extracts, renames module, runs `go mod tidy` |
| T03 | `rapidgo new existing-dir` | Fails with "already exists" error |
| T04 | `rapidgo new "bad/name"` | Fails with "invalid project name" |
| T05 | Scaffolded project `go build ./...` | Compiles |
| T06 | Scaffolded project `go run cmd/main.go version` | Prints version |
| T07 | Scaffolded project has no `RapidGo-starter` references | `grep -r "RapidGo-starter" myapp/` returns nothing |
| T08 | Zip slip protection | Malicious zip paths are rejected |

### Verification

```bash
go build ./...
go test ./core/cli/ -run TestNew -v

# Integration test:
go run cmd/main.go new testproject
cd testproject
go build ./...
go run cmd/main.go version
grep -r "RapidGo-starter" .  # should return nothing
cd ..
rm -rf testproject
```

---

## Step D2: Documentation (READMEs)

### Branch

`feature/v2-10-library-readme` (from `v2`, independent of D1)

### Objective

Write focused READMEs for both repos. The library README explains how to `go get` and use the packages. The starter README explains how to clone, configure, and develop.

### Files Changed

| Action | File | Repo |
|--------|------|------|
| REWRITE | `README.md` | Library (RapidGo) |
| CREATE | `README.md` | Starter (RapidGo-starter) |

### Library README Outline

```markdown
# RapidGo

A batteries-included Go web framework with Laravel-style developer experience.

## Install

    go get github.com/RAiWorks/RapidGo

## Quick Start

    rapidgo new myapp
    cd myapp
    cp .env.example .env
    go run cmd/main.go serve

## Package Index

| Package | Import | Purpose |
|---------|--------|---------|
| `core/app` | Application lifecycle | ... |
| `core/auth` | JWT authentication | ... |
| `core/cache` | File + Redis caching | ... |
| ... | ... | ... |

## Creating a Project

See [RapidGo-starter](https://github.com/RAiWorks/RapidGo-starter).

## Hook System

The `core/cli` package provides 6 hooks for wiring application code:
- `SetBootstrap()` — register service providers
- `SetRoutes()` — register routes
- `SetJobRegistrar()` — register job handlers
- `SetScheduleRegistrar()` — register scheduled tasks
- `SetModelRegistry()` — provide models for AutoMigrate
- `SetSeeder()` — run database seeders

## License

MIT
```

### Starter README Outline

```markdown
# RapidGo Starter

A scaffold project for the [RapidGo](https://github.com/RAiWorks/RapidGo) framework.

## Getting Started

### Option 1: CLI (recommended)

    go install github.com/RAiWorks/RapidGo/cmd/rapidgo@latest
    rapidgo new myapp

### Option 2: Clone

    git clone https://github.com/RAiWorks/RapidGo-starter myapp
    cd myapp
    # Update module name in go.mod and all .go files
    go mod tidy

### Configure

    cp .env.example .env
    # Edit .env with your database credentials

### Run

    go run cmd/main.go serve
    go run cmd/main.go migrate
    go run cmd/main.go db:seed

## Project Structure

    cmd/main.go           ← Entry point with hook wiring
    app/providers/        ← Service providers
    app/helpers/          ← Utility functions
    app/services/         ← Business logic
    app/jobs/             ← Queue job handlers
    app/schedule/         ← Scheduled tasks
    routes/               ← Route definitions
    http/controllers/     ← Request handlers
    database/models/      ← GORM models
    database/migrations/  ← Database migrations
    database/seeders/     ← Seed data
    resources/            ← Views, translations, static files

## Hook Wiring

See `cmd/main.go` for how hooks connect your app to the framework.

## License

MIT
```

### Verification

- [ ] Library README renders correctly on GitHub
- [ ] Starter README renders correctly on GitHub
- [ ] All links work
- [ ] Package index matches actual packages
- [ ] Commands in README actually work

---

## Release Process

After Phase D is complete:

### 1. Final Verification

```bash
# Library
cd RapidGo
go build ./...
go test ./... -count=1
go vet ./...
grep -rn "RAiWorks/RapidGo/app\|RAiWorks/RapidGo/routes\|RAiWorks/RapidGo/http\|RAiWorks/RapidGo/plugins" core/
# Expected: no output

# Starter
cd ../RapidGo-starter
go build ./...
go test ./... -count=1
go vet ./...
go run cmd/main.go version
```

### 2. Tag and Release

```bash
# Library
cd RapidGo
git checkout v2
git tag -a v2.0.0 -m "v2.0.0 — Importable library split"
git push origin v2.0.0

# Set v2 as default branch on GitHub (Settings → Default branch → v2)

# Starter
cd ../RapidGo-starter
git tag -a v1.0.0 -m "v1.0.0 — Initial starter template for RapidGo v2"
git push origin v1.0.0
```

### 3. Verify Public Access

```bash
# Test go get works
mkdir /tmp/test-import
cd /tmp/test-import
go mod init test
go get github.com/RAiWorks/RapidGo@v2.0.0
# Should succeed

# Test rapidgo new works
go install github.com/RAiWorks/RapidGo/cmd/rapidgo@v2.0.0
rapidgo new testapp
cd testapp
go build ./...
```

### 4. Create GitHub Releases

Create releases on both repos with:
- Changelog summary
- Migration guide (for users upgrading from v1 → v2)
- Links between the two repos

---

## Phase D Checklist

| # | Check | Status |
|---|-------|--------|
| 1 | `rapidgo new myapp` creates working project | - [ ] |
| 2 | Scaffolded project builds and runs all commands | - [ ] |
| 3 | Library README complete with package index | - [ ] |
| 4 | Starter README complete with getting-started guide | - [ ] |
| 5 | v2.0.0 tagged on library | - [ ] |
| 6 | v1.0.0 tagged on starter | - [ ] |
| 7 | `go get github.com/RAiWorks/RapidGo@v2.0.0` works | - [ ] |
| 8 | v2 set as default branch on GitHub | - [ ] |
| 9 | GitHub releases created | - [ ] |
