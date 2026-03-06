package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/spf13/cobra"
)

var makeMigrationCmd = &cobra.Command{
	Use:   "make:migration [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		timestamp := time.Now().Format("20060102150405")
		version := timestamp + "_" + toSnakeCase(name)
		filename := version + ".go"
		path := filepath.Join("database", "migrations", filename)

		if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()

		t := template.Must(template.New("migration").Parse(migrationTpl))
		if err := t.Execute(f, map[string]string{"Version": version}); err != nil {
			return fmt.Errorf("failed to write template: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Migration created: %s\n", path)
		return nil
	},
}

// toSnakeCase converts PascalCase/camelCase to snake_case.
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

var migrationTpl = `package migrations

import "gorm.io/gorm"

func init() {
	Register(Migration{
		Version: "{{.Version}}",
		Up: func(db *gorm.DB) error {
			// TODO: implement migration
			return nil
		},
		Down: func(db *gorm.DB) error {
			// TODO: implement rollback
			return nil
		},
	})
}
`
