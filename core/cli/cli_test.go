package cli

import (
	"bytes"
	"testing"
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
	for _, want := range []string{"RGo", "serve", "version"} {
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

// TC-05: NewApp returns booted application
func TestNewApp_ReturnsBootedApp(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_OUTPUT", "stdout")

	application := NewApp()
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
