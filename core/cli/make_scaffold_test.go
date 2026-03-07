package cli

import (
	"bytes"
	"errors"
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

// --- Admin Scaffold Tests ---

// TC-06: make:admin generates admin controller file.
func TestAdminScaffold_Controller(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "http", "controllers", "admin", "article_controller.go")
	var buf bytes.Buffer
	err := adminScaffold("Admin controller", "Article", adminControllerTpl, path, &buf)
	if err != nil {
		t.Fatalf("admin scaffold controller: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "type ArticleController struct{}") {
		t.Fatal("expected ArticleController struct declaration")
	}
	if !strings.Contains(content, "package admin") {
		t.Fatal("expected package admin")
	}
	if !strings.Contains(content, `"admin/article/index.html"`) {
		t.Fatal("expected template path with snake_case resource")
	}
	if !strings.Contains(buf.String(), "Admin controller created:") {
		t.Fatalf("expected confirmation message, got: %s", buf.String())
	}
}

// TC-07: make:admin generates 4 view templates.
func TestAdminScaffold_Views(t *testing.T) {
	base := t.TempDir()
	views := []string{"index.html", "show.html", "create.html", "edit.html"}
	tpls := []string{adminIndexTpl, adminShowTpl, adminCreateTpl, adminEditTpl}

	for i, view := range views {
		path := filepath.Join(base, "resources", "views", "admin", "product", view)
		var buf bytes.Buffer
		err := adminScaffold("Admin view", "Product", tpls[i], path, &buf)
		if err != nil {
			t.Fatalf("admin scaffold %s: %v", view, err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", view, err)
		}
		content := string(data)
		if !strings.Contains(content, "{{ .title }}") {
			t.Fatalf("%s: expected {{ .title }} Gin template variable", view)
		}
	}
}

// TC-08: make:admin generates layout template.
func TestAdminScaffold_Layout(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "resources", "views", "admin", "layout.html")
	var buf bytes.Buffer
	err := adminScaffold("Admin layout", "Post", adminLayoutTpl, path, &buf)
	if err != nil {
		t.Fatalf("admin scaffold layout: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read layout: %v", err)
	}
	if !strings.Contains(string(data), "Admin Panel") {
		t.Fatal("expected 'Admin Panel' in layout")
	}
}

// TC-09: make:admin generates dashboard template.
func TestAdminScaffold_Dashboard(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "resources", "views", "admin", "dashboard.html")
	var buf bytes.Buffer
	err := adminScaffold("Admin dashboard", "Post", adminDashboardTpl, path, &buf)
	if err != nil {
		t.Fatalf("admin scaffold dashboard: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read dashboard: %v", err)
	}
	if !strings.Contains(string(data), "Dashboard") {
		t.Fatal("expected 'Dashboard' in dashboard template")
	}
}

// TC-10: make:admin skips layout if already exists.
func TestAdminScaffold_SkipsExistingLayout(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "resources", "views", "admin", "layout.html")

	// Create layout first with custom content.
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		t.Fatal(err)
	}
	original := []byte("custom layout")
	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatal(err)
	}

	// adminScaffold should fail with os.ErrExist.
	var buf bytes.Buffer
	err := adminScaffold("Admin layout", "Post", adminLayoutTpl, path, &buf)
	if err == nil {
		t.Fatal("expected error when layout exists")
	}
	if !errors.Is(err, os.ErrExist) {
		t.Fatalf("expected os.ErrExist, got: %v", err)
	}

	// Original content should be preserved.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "custom layout" {
		t.Fatal("layout file was overwritten")
	}
}

// TC-11: make:admin prevents controller overwrite.
func TestAdminScaffold_PreventsDuplicateController(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "http", "controllers", "admin", "user_controller.go")
	var buf bytes.Buffer

	// First call succeeds.
	if err := adminScaffold("Admin controller", "User", adminControllerTpl, path, &buf); err != nil {
		t.Fatalf("first admin scaffold: %v", err)
	}

	// Second call should fail.
	err := adminScaffold("Admin controller", "User", adminControllerTpl, path, &buf)
	if err == nil {
		t.Fatal("expected error on duplicate admin scaffold")
	}
	if !strings.Contains(err.Error(), "file already exists") {
		t.Fatalf("expected 'file already exists' error, got: %v", err)
	}
}

// TC-12: adminScaffold uses custom delimiters correctly.
func TestAdminScaffold_CustomDelimiters(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "test_output.go")
	var buf bytes.Buffer

	// Template uses [[ ]] for scaffold-time and {{ }} for Gin-time.
	tpl := `name=[[.Name]] resource=[[.Resource]] gin={{ .title }}`
	err := adminScaffold("Test", "MyPost", tpl, path, &buf)
	if err != nil {
		t.Fatalf("admin scaffold: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	content := string(data)

	// [[ ]] should be substituted.
	if !strings.Contains(content, "name=MyPost") {
		t.Fatalf("expected Name substitution, got: %s", content)
	}
	if !strings.Contains(content, "resource=my_post") {
		t.Fatalf("expected Resource substitution, got: %s", content)
	}
	// {{ }} should pass through literally.
	if !strings.Contains(content, "gin={{ .title }}") {
		t.Fatalf("expected {{ .title }} passthrough, got: %s", content)
	}
}
