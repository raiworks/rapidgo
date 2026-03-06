package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TC-01: make:controller generates a controller file.
func TestScaffold_Controller(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "http", "controllers")
	var buf bytes.Buffer
	err := scaffold("Controller", "OrderController", controllerTpl, dir, &buf)
	if err != nil {
		t.Fatalf("scaffold controller: %v", err)
	}

	path := filepath.Join(dir, "order_controller.go")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "type OrderController struct{}") {
		t.Fatal("expected struct declaration in controller")
	}
	if !strings.Contains(content, "package controllers") {
		t.Fatal("expected package controllers")
	}
	if !strings.Contains(buf.String(), "Controller created:") {
		t.Fatalf("expected confirmation message, got: %s", buf.String())
	}
}

// TC-02: make:model generates a model file.
func TestScaffold_Model(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "database", "models")
	var buf bytes.Buffer
	err := scaffold("Model", "Product", modelTpl, dir, &buf)
	if err != nil {
		t.Fatalf("scaffold model: %v", err)
	}

	path := filepath.Join(dir, "product.go")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "type Product struct") {
		t.Fatal("expected struct declaration in model")
	}
	if !strings.Contains(content, "BaseModel") {
		t.Fatal("expected BaseModel embed")
	}
}

// TC-03: make:service generates a service file.
func TestScaffold_Service(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "app", "services")
	var buf bytes.Buffer
	err := scaffold("Service", "PaymentService", serviceTpl, dir, &buf)
	if err != nil {
		t.Fatalf("scaffold service: %v", err)
	}

	path := filepath.Join(dir, "payment_service.go")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "func NewPaymentService(db *gorm.DB)") {
		t.Fatal("expected constructor in service")
	}
}

// TC-04: make:provider generates a provider file.
func TestScaffold_Provider(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "app", "providers")
	var buf bytes.Buffer
	err := scaffold("Provider", "CacheProvider", providerTpl, dir, &buf)
	if err != nil {
		t.Fatalf("scaffold provider: %v", err)
	}

	path := filepath.Join(dir, "cache_provider.go")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "func (p *CacheProvider) Register") {
		t.Fatal("expected Register method in provider")
	}
	if !strings.Contains(content, "func (p *CacheProvider) Boot") {
		t.Fatal("expected Boot method in provider")
	}
}

// TC-05: Duplicate scaffold prevents overwrite.
func TestScaffold_DuplicatePreventsOverwrite(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "http", "controllers")
	var buf bytes.Buffer

	// First call succeeds.
	if err := scaffold("Controller", "DupCtrl", controllerTpl, dir, &buf); err != nil {
		t.Fatalf("first scaffold: %v", err)
	}

	// Second call with same name should fail.
	err := scaffold("Controller", "DupCtrl", controllerTpl, dir, &buf)
	if err == nil {
		t.Fatal("expected error on duplicate scaffold")
	}
	if !strings.Contains(err.Error(), "file already exists") {
		t.Fatalf("expected 'file already exists' error, got: %v", err)
	}
}
