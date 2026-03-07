package models

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TC-01: BaseModel has expected fields
func TestBaseModel_Fields(t *testing.T) {
	rt := reflect.TypeOf(BaseModel{})

	checks := []struct {
		name     string
		wantType string
	}{
		{"ID", "uint"},
		{"CreatedAt", "time.Time"},
		{"UpdatedAt", "time.Time"},
		{"DeletedAt", "gorm.DeletedAt"},
	}

	for _, c := range checks {
		f, ok := rt.FieldByName(c.name)
		if !ok {
			t.Fatalf("expected field %q on BaseModel", c.name)
		}
		if f.Type.String() != c.wantType {
			t.Fatalf("expected %q to be %s, got %s", c.name, c.wantType, f.Type.String())
		}
	}
}

// TC-02: User embeds BaseModel
func TestUser_EmbedsBaseModel(t *testing.T) {
	user := User{BaseModel: BaseModel{ID: 42}}
	if user.ID != 42 {
		t.Fatalf("expected user.ID == 42, got %d", user.ID)
	}
}

// TC-03: User password excluded from JSON
func TestUser_PasswordExcludedFromJSON(t *testing.T) {
	user := User{
		BaseModel: BaseModel{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Name:      "Alice",
		Email:     "alice@test.com",
		Password:  "secret",
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if _, ok := m["password"]; ok {
		t.Fatal("expected password to be excluded from JSON")
	}
	if _, ok := m["name"]; !ok {
		t.Fatal("expected name to be present in JSON")
	}
}

// helper: open SQLite :memory: for integration tests
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&User{}, &Post{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return db
}

// TC-04: GORM AutoMigrate succeeds
func TestModels_AutoMigrate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := db.AutoMigrate(&User{}, &Post{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
}

// TC-05: GORM creates and queries User
func TestUser_CreateAndQuery(t *testing.T) {
	db := setupTestDB(t)

	user := User{Name: "Alice", Email: "alice@test.com", Password: "hash"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	var found User
	if err := db.Where("email = ?", "alice@test.com").First(&found).Error; err != nil {
		t.Fatalf("Query user failed: %v", err)
	}

	if found.Name != "Alice" {
		t.Fatalf("expected name Alice, got %s", found.Name)
	}
	if found.ID == 0 {
		t.Fatal("expected ID > 0")
	}
	if found.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
}

// TC-06: GORM creates Post with foreign key
func TestPost_CreateWithForeignKey(t *testing.T) {
	db := setupTestDB(t)

	user := User{Name: "Bob", Email: "bob@test.com", Password: "hash"}
	db.Create(&user)

	post := Post{Title: "Hello", Slug: "hello", Body: "World", UserID: user.ID}
	if err := db.Create(&post).Error; err != nil {
		t.Fatalf("Create post failed: %v", err)
	}

	var found Post
	if err := db.First(&found, post.ID).Error; err != nil {
		t.Fatalf("Query post failed: %v", err)
	}
	if found.UserID != user.ID {
		t.Fatalf("expected UserID %d, got %d", user.ID, found.UserID)
	}
}

// TC-07: GORM Preload loads User→Posts
func TestUser_PreloadPosts(t *testing.T) {
	db := setupTestDB(t)

	user := User{Name: "Carol", Email: "carol@test.com", Password: "hash"}
	db.Create(&user)

	db.Create(&Post{Title: "Post 1", Slug: "post-1", Body: "Body 1", UserID: user.ID})
	db.Create(&Post{Title: "Post 2", Slug: "post-2", Body: "Body 2", UserID: user.ID})

	var found User
	if err := db.Preload("Posts").First(&found, user.ID).Error; err != nil {
		t.Fatalf("Preload query failed: %v", err)
	}
	if len(found.Posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(found.Posts))
	}
}

// TC-08: User defaults (role, active)
func TestUser_Defaults(t *testing.T) {
	db := setupTestDB(t)

	user := User{Name: "Dave", Email: "dave@test.com", Password: "hash"}
	db.Create(&user)

	var found User
	db.First(&found, user.ID)

	if found.Role != "user" {
		t.Fatalf("expected default role 'user', got %q", found.Role)
	}
	if !found.Active {
		t.Fatal("expected default active to be true")
	}
}

// --- BeforeCreate Hook ---

// TC-15: BeforeCreate hashes plaintext password
func TestUser_BeforeCreate_HashesPassword(t *testing.T) {
	db := setupTestDB(t)

	user := User{Name: "Eve", Email: "eve@test.com", Password: "mypassword"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.Password == "mypassword" {
		t.Fatal("expected password to be hashed, but it was unchanged")
	}
	if !strings.HasPrefix(user.Password, "$2a$") {
		t.Fatalf("expected bcrypt hash starting with $2a$, got %q", user.Password[:10])
	}
}

// TC-16: BeforeCreate skips already-hashed password
func TestUser_BeforeCreate_SkipsHashed(t *testing.T) {
	db := setupTestDB(t)

	hashed := "$2a$10$abcdefghijklmnopqrstuuABCDEFGHIJKLMNOPQRSTUVWXYZ012"
	user := User{Name: "Frank", Email: "frank@test.com", Password: hashed}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.Password != hashed {
		t.Fatal("expected already-hashed password to be unchanged")
	}
}
