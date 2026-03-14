package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/raiworks/rapidgo/v2/core/app"
	"github.com/raiworks/rapidgo/v2/core/container"
	"github.com/raiworks/rapidgo/v2/core/service"
)

// TC-01: Version constant is set
func TestVersion_IsSet(t *testing.T) {
	if Version == "" {
		t.Fatal("expected Version to be non-empty")
	}
}

// TC-02: Version command output
func TestVersionCmd_Output(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"version"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte(Version)) {
		t.Fatalf("expected output to contain %q, got %q", Version, output)
	}
}

// TC-03: Root command shows help (no args)
func TestRootCmd_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, want := range []string{"RapidGo", "serve", "version"} {
		if !bytes.Contains([]byte(output), []byte(want)) {
			t.Errorf("expected output to contain %q, got %q", want, output)
		}
	}
}

// TC-04: Serve command has port flag
func TestServeCmd_HasPortFlag(t *testing.T) {
	f := serveCmd.Flags().Lookup("port")
	if f == nil {
		t.Fatal("expected 'port' flag on serve command")
	}
	if f.Shorthand != "p" {
		t.Fatalf("expected port flag shorthand 'p', got %q", f.Shorthand)
	}
}

// TC-05: Serve command has mode flag
func TestServeCmd_HasModeFlag(t *testing.T) {
	f := serveCmd.Flags().Lookup("mode")
	if f == nil {
		t.Fatal("expected 'mode' flag on serve command")
	}
	if f.Shorthand != "m" {
		t.Fatalf("expected mode flag shorthand 'm', got %q", f.Shorthand)
	}
}

// TC-06: Serve command usage includes mode flag
func TestServeCmd_HelpIncludesModeFlag(t *testing.T) {
	usage := serveCmd.UsageString()

	if !bytes.Contains([]byte(usage), []byte("--mode")) {
		t.Fatal("expected serve usage to contain '--mode'")
	}
	if !bytes.Contains([]byte(usage), []byte("-m")) {
		t.Fatal("expected serve usage to contain '-m' shorthand")
	}
}

// TC-07: Invalid mode string causes ParseMode error
func TestParseMode_InvalidReturnsError(t *testing.T) {
	_, err := service.ParseMode("invalid")
	if err == nil {
		t.Fatal("expected error for invalid mode string")
	}
}

// TC-05: NewApp returns booted application
func TestNewApp_ReturnsBootedApp(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_OUTPUT", "stdout")

	// Set a bootstrap function that registers test bindings
	original := bootstrapFn
	t.Cleanup(func() { bootstrapFn = original })

	SetBootstrap(func(a *app.App, mode service.Mode) {
		a.Container.Bind("router", func(_ *container.Container) interface{} { return "fake-router" })
		a.Container.Bind("db", func(_ *container.Container) interface{} { return "fake-db" })
	})

	application := NewApp(service.ModeAll)
	if application == nil {
		t.Fatal("expected non-nil application")
	}
	if application.Container == nil {
		t.Fatal("expected non-nil container")
	}
	if !application.Container.Has("router") {
		t.Fatal("expected 'router' binding to be registered")
	}
	if !application.Container.Has("db") {
		t.Fatal("expected 'db' binding to be registered")
	}
}

// TC-06: Serve command registered on root
func TestRootCmd_HasServeCommand(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "serve" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected 'serve' command to be registered on root")
	}
}

// TC-07: Version command registered on root
func TestRootCmd_HasVersionCommand(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "version" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected 'version' command to be registered on root")
	}
}

// TC-10 (Feature #12): toSnakeCase converts PascalCase to snake_case
func TestToSnakeCase(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"CreateUsersTable", "create_users_table"},
		{"addEmailIndex", "add_email_index"},
		{"simple", "simple"},
		{"ABCTest", "a_b_c_test"},
	}

	for _, tc := range cases {
		got := toSnakeCase(tc.input)
		if got != tc.want {
			t.Errorf("toSnakeCase(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// TC-M01: make:module creates 4 files in modules/<snake_name>/
func TestMakeModuleCmd_Creates4Files(t *testing.T) {
	tmp := t.TempDir()
	original, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(original)

	makeModuleCmd.SetArgs([]string{"Product"})
	var buf bytes.Buffer
	makeModuleCmd.SetOut(&buf)

	if err := makeModuleCmd.RunE(makeModuleCmd, []string{"Product"}); err != nil {
		t.Fatalf("make:module failed: %v", err)
	}

	expected := []string{
		filepath.Join("modules", "product", "models.go"),
		filepath.Join("modules", "product", "service.go"),
		filepath.Join("modules", "product", "controller.go"),
		filepath.Join("modules", "product", "routes.go"),
	}

	for _, path := range expected {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", path)
		}
	}

	output := buf.String()
	if !strings.Contains(output, "Models created") {
		t.Error("expected output to contain 'Models created'")
	}
	if !strings.Contains(output, "Routes created") {
		t.Error("expected output to contain 'Routes created'")
	}
}
