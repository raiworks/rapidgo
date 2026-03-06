package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalDriver stores files on the local filesystem.
type LocalDriver struct {
	BasePath string // e.g. "storage/uploads"
	BaseURL  string // e.g. "/uploads"
}

// safePath resolves the given path within BasePath and ensures it does not
// escape the base directory via path traversal (e.g. "../../etc/passwd").
func (d *LocalDriver) safePath(path string) (string, error) {
	absBase, err := filepath.Abs(d.BasePath)
	if err != nil {
		return "", err
	}
	full := filepath.Join(absBase, filepath.FromSlash(path))
	absTarget, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absTarget, absBase+string(filepath.Separator)) && absTarget != absBase {
		return "", errors.New("path traversal detected")
	}
	return absTarget, nil
}

// Put writes content to the given path under BasePath.
// Intermediate directories are created automatically.
func (d *LocalDriver) Put(path string, content io.Reader) (string, error) {
	fullPath, err := d.safePath(path)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, content); err != nil {
		return "", err
	}
	return path, nil
}

// Get opens the file at the given path and returns an io.ReadCloser.
func (d *LocalDriver) Get(path string) (io.ReadCloser, error) {
	fullPath, err := d.safePath(path)
	if err != nil {
		return nil, err
	}
	return os.Open(fullPath)
}

// Delete removes the file at the given path.
func (d *LocalDriver) Delete(path string) error {
	fullPath, err := d.safePath(path)
	if err != nil {
		return err
	}
	return os.Remove(fullPath)
}

// URL returns the public URL for the given path.
func (d *LocalDriver) URL(path string) string {
	return d.BaseURL + "/" + path
}
