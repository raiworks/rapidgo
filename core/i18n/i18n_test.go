package i18n

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func writeJSON(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return p
}

// TC-01: LoadFile loads JSON translations.
func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	p := writeJSON(t, dir, "en.json", `{"hello":"Hello!"}`)

	tr := NewTranslator("en")
	if err := tr.LoadFile("en", p); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	got := tr.Get("en", "hello")
	if got != "Hello!" {
		t.Fatalf("Get = %q, want %q", got, "Hello!")
	}
}

// TC-02: Get missing key returns raw key.
func TestGetMissingKey(t *testing.T) {
	tr := NewTranslator("en")
	got := tr.Get("en", "nope")
	if got != "nope" {
		t.Fatalf("Get = %q, want %q", got, "nope")
	}
}

// TC-03: Get falls back to fallback locale.
func TestGetFallback(t *testing.T) {
	dir := t.TempDir()
	p := writeJSON(t, dir, "en.json", `{"greeting":"Hi"}`)

	tr := NewTranslator("en")
	_ = tr.LoadFile("en", p)

	got := tr.Get("fr", "greeting")
	if got != "Hi" {
		t.Fatalf("Get(fr) = %q, want fallback %q", got, "Hi")
	}
}

// TC-04: Get with template args interpolates.
func TestGetInterpolation(t *testing.T) {
	dir := t.TempDir()
	p := writeJSON(t, dir, "en.json", `{"welcome":"Welcome, {{.Name}}!"}`)

	tr := NewTranslator("en")
	_ = tr.LoadFile("en", p)

	got := tr.Get("en", "welcome", map[string]string{"Name": "Carlos"})
	if got != "Welcome, Carlos!" {
		t.Fatalf("Get = %q, want %q", got, "Welcome, Carlos!")
	}
}

// TC-05: Get with no args returns plain message.
func TestGetPlain(t *testing.T) {
	dir := t.TempDir()
	p := writeJSON(t, dir, "en.json", `{"errors.not_found":"Resource not found"}`)

	tr := NewTranslator("en")
	_ = tr.LoadFile("en", p)

	got := tr.Get("en", "errors.not_found")
	if got != "Resource not found" {
		t.Fatalf("Get = %q, want %q", got, "Resource not found")
	}
}

// TC-06: LoadFile returns error for missing file.
func TestLoadFileMissing(t *testing.T) {
	tr := NewTranslator("en")
	err := tr.LoadFile("en", "nonexistent.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

// TC-07: LoadFile returns error for invalid JSON.
func TestLoadFileInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := writeJSON(t, dir, "bad.json", `{not valid json}`)

	tr := NewTranslator("en")
	err := tr.LoadFile("en", p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// TC-08: LoadDir loads all JSON files.
func TestLoadDir(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, dir, "en.json", `{"hi":"Hello"}`)
	writeJSON(t, dir, "es.json", `{"hi":"Hola"}`)

	tr := NewTranslator("en")
	if err := tr.LoadDir(dir); err != nil {
		t.Fatalf("LoadDir: %v", err)
	}

	if got := tr.Get("en", "hi"); got != "Hello" {
		t.Fatalf("en hi = %q", got)
	}
	if got := tr.Get("es", "hi"); got != "Hola" {
		t.Fatalf("es hi = %q", got)
	}
}

// TC-09: LoadDir skips non-JSON files.
func TestLoadDirSkipsNonJSON(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, dir, "en.json", `{"k":"v"}`)
	writeJSON(t, dir, "notes.txt", `not json at all`)

	tr := NewTranslator("en")
	if err := tr.LoadDir(dir); err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	if got := tr.Get("en", "k"); got != "v" {
		t.Fatalf("Get = %q", got)
	}
}

// TC-10: Concurrent Get is safe.
func TestConcurrentGet(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, dir, "en.json", `{"key":"value"}`)

	tr := NewTranslator("en")
	_ = tr.LoadFile("en", filepath.Join(dir, "en.json"))

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tr.Get("en", "key")
		}()
	}
	wg.Wait()
}
