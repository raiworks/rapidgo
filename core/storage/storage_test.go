package storage

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// helper: create a LocalDriver rooted in a temp dir.
func newTestDriver(t *testing.T) *LocalDriver {
	t.Helper()
	return &LocalDriver{
		BasePath: t.TempDir(),
		BaseURL:  "/uploads",
	}
}

// TC-01: Put writes file to disk
func TestLocalDriver_Put_WritesFile(t *testing.T) {
	d := newTestDriver(t)
	path, err := d.Put("hello.txt", strings.NewReader("hello world"))
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
	if path != "hello.txt" {
		t.Fatalf("expected path 'hello.txt', got %q", path)
	}
	data, err := os.ReadFile(filepath.Join(d.BasePath, "hello.txt"))
	if err != nil {
		t.Fatalf("file not found on disk: %v", err)
	}
	if string(data) != "hello world" {
		t.Fatalf("content mismatch: got %q", string(data))
	}
}

// TC-02: Put creates intermediate directories
func TestLocalDriver_Put_CreatesDirectories(t *testing.T) {
	d := newTestDriver(t)
	_, err := d.Put("a/b/c/deep.txt", strings.NewReader("nested"))
	if err != nil {
		t.Fatalf("Put with nested path failed: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(d.BasePath, "a", "b", "c", "deep.txt"))
	if err != nil {
		t.Fatalf("nested file not found: %v", err)
	}
	if string(data) != "nested" {
		t.Fatalf("content mismatch: got %q", string(data))
	}
}

// TC-03: Get returns file content
func TestLocalDriver_Get_ReturnsContent(t *testing.T) {
	d := newTestDriver(t)
	_, _ = d.Put("read.txt", strings.NewReader("read me"))

	rc, err := d.Get("read.txt")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer rc.Close()
	data, _ := io.ReadAll(rc)
	if string(data) != "read me" {
		t.Fatalf("expected 'read me', got %q", string(data))
	}
}

// TC-04: Get returns error for missing file
func TestLocalDriver_Get_MissingFile(t *testing.T) {
	d := newTestDriver(t)
	_, err := d.Get("nope.txt")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

// TC-05: Delete removes file
func TestLocalDriver_Delete_RemovesFile(t *testing.T) {
	d := newTestDriver(t)
	_, _ = d.Put("del.txt", strings.NewReader("bye"))

	err := d.Delete("del.txt")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(d.BasePath, "del.txt")); !os.IsNotExist(err) {
		t.Fatal("file still exists after Delete")
	}
}

// TC-06: Delete returns error for missing file
func TestLocalDriver_Delete_MissingFile(t *testing.T) {
	d := newTestDriver(t)
	err := d.Delete("gone.txt")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

// TC-07: URL returns BaseURL + path
func TestLocalDriver_URL(t *testing.T) {
	d := newTestDriver(t)
	url := d.URL("images/photo.jpg")
	if url != "/uploads/images/photo.jpg" {
		t.Fatalf("expected '/uploads/images/photo.jpg', got %q", url)
	}
}

// TC-08: Path traversal in Put is rejected
func TestLocalDriver_Put_PathTraversal(t *testing.T) {
	d := newTestDriver(t)
	_, err := d.Put("../../etc/passwd", strings.NewReader("evil"))
	if err == nil {
		t.Fatal("expected error for path traversal in Put")
	}
}

// TC-09: Path traversal in Get is rejected
func TestLocalDriver_Get_PathTraversal(t *testing.T) {
	d := newTestDriver(t)
	_, err := d.Get("../../etc/passwd")
	if err == nil {
		t.Fatal("expected error for path traversal in Get")
	}
}

// TC-10: Path traversal in Delete is rejected
func TestLocalDriver_Delete_PathTraversal(t *testing.T) {
	d := newTestDriver(t)
	err := d.Delete("../../etc/passwd")
	if err == nil {
		t.Fatal("expected error for path traversal in Delete")
	}
}

// TC-11: NewDriver returns LocalDriver by default
func TestNewDriver_DefaultLocal(t *testing.T) {
	t.Setenv("STORAGE_DRIVER", "")
	t.Setenv("STORAGE_LOCAL_PATH", t.TempDir())

	drv, err := NewDriver()
	if err != nil {
		t.Fatalf("NewDriver failed: %v", err)
	}
	if drv == nil {
		t.Fatal("expected non-nil driver")
	}
	if _, ok := drv.(*LocalDriver); !ok {
		t.Fatalf("expected *LocalDriver, got %T", drv)
	}
}

// TC-12: NewDriver returns error for unknown driver
func TestNewDriver_UnknownDriver(t *testing.T) {
	t.Setenv("STORAGE_DRIVER", "unknown")

	_, err := NewDriver()
	if err == nil {
		t.Fatal("expected error for unknown driver")
	}
}
