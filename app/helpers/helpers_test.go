package helpers

import (
	"os"
	"sort"
	"testing"
	"time"
)

// --- Password ---

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "secret" {
		t.Fatal("hash should not equal plain-text password")
	}
}

func TestCheckPassword_Valid(t *testing.T) {
	hash, _ := HashPassword("secret")
	if !CheckPassword(hash, "secret") {
		t.Fatal("expected CheckPassword to return true for correct password")
	}
}

func TestCheckPassword_Invalid(t *testing.T) {
	hash, _ := HashPassword("secret")
	if CheckPassword(hash, "wrong") {
		t.Fatal("expected CheckPassword to return false for wrong password")
	}
}

// --- Random ---

func TestRandomString_Length(t *testing.T) {
	s := RandomString(16)
	if len(s) != 32 {
		t.Fatalf("expected hex length 32, got %d", len(s))
	}
}

func TestRandomString_Unique(t *testing.T) {
	a := RandomString(16)
	b := RandomString(16)
	if a == b {
		t.Fatal("two random strings should not be equal")
	}
}

// --- String ---

func TestSlugify(t *testing.T) {
	got := Slugify("Hello World!")
	if got != "hello-world" {
		t.Fatalf("expected 'hello-world', got %q", got)
	}
}

func TestTruncate_Short(t *testing.T) {
	got := Truncate("hi", 10)
	if got != "hi" {
		t.Fatalf("expected 'hi', got %q", got)
	}
}

func TestTruncate_Long(t *testing.T) {
	got := Truncate("Hello World!", 8)
	if got != "Hello..." {
		t.Fatalf("expected 'Hello...', got %q", got)
	}
}

func TestContains(t *testing.T) {
	if !Contains("Hello World", "hello") {
		t.Fatal("expected case-insensitive match")
	}
}

func TestContains_NoMatch(t *testing.T) {
	if Contains("Hello World", "xyz") {
		t.Fatal("expected no match")
	}
}

func TestTitle(t *testing.T) {
	got := Title("hello world")
	if got != "Hello World" {
		t.Fatalf("expected 'Hello World', got %q", got)
	}
}

func TestExcerpt(t *testing.T) {
	got := Excerpt("The quick brown fox jumps over", 3)
	if got != "The quick brown..." {
		t.Fatalf("expected 'The quick brown...', got %q", got)
	}
}

func TestStripHTML(t *testing.T) {
	got := StripHTML("<p>Hello <b>World</b></p>")
	if got != "Hello World" {
		t.Fatalf("expected 'Hello World', got %q", got)
	}
}

func TestMask(t *testing.T) {
	got := Mask("secret123", 2, 2)
	if got != "se*****23" {
		t.Fatalf("expected 'se*****23', got %q", got)
	}
}

// --- Number ---

func TestFormatBytes_Zero(t *testing.T) {
	got := FormatBytes(0)
	if got != "0 B" {
		t.Fatalf("expected '0 B', got %q", got)
	}
}

func TestFormatBytes_KB(t *testing.T) {
	got := FormatBytes(1536)
	if got != "1.50 KB" {
		t.Fatalf("expected '1.50 KB', got %q", got)
	}
}

func TestClamp_InRange(t *testing.T) {
	if Clamp(5, 1, 10) != 5 {
		t.Fatal("expected 5")
	}
}

func TestClamp_BelowMin(t *testing.T) {
	if Clamp(-1, 0, 10) != 0 {
		t.Fatal("expected 0")
	}
}

func TestClamp_AboveMax(t *testing.T) {
	if Clamp(20, 0, 10) != 10 {
		t.Fatal("expected 10")
	}
}

// --- Time ---

func TestTimeAgo_JustNow(t *testing.T) {
	got := TimeAgo(time.Now())
	if got != "just now" {
		t.Fatalf("expected 'just now', got %q", got)
	}
}

func TestTimeAgo_MinutesAgo(t *testing.T) {
	got := TimeAgo(time.Now().Add(-5 * time.Minute))
	if got != "5 minutes ago" {
		t.Fatalf("expected '5 minutes ago', got %q", got)
	}
}

func TestFormatDate(t *testing.T) {
	tm := time.Date(2026, 3, 6, 14, 30, 0, 0, time.UTC)
	got := FormatDate(tm)
	if got != "Mar 6, 2026 2:30 PM" {
		t.Fatalf("expected 'Mar 6, 2026 2:30 PM', got %q", got)
	}
}

// --- Data ---

func TestStructToMap(t *testing.T) {
	type S struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	m, err := StructToMap(S{Name: "Alice", Age: 30})
	if err != nil {
		t.Fatalf("StructToMap error: %v", err)
	}
	if m["name"] != "Alice" {
		t.Fatalf("expected 'Alice', got %v", m["name"])
	}
}

func TestMapKeys(t *testing.T) {
	m := map[string]interface{}{"a": 1, "b": 2, "c": 3}
	keys := MapKeys(m)
	sort.Strings(keys)
	if len(keys) != 3 || keys[0] != "a" || keys[1] != "b" || keys[2] != "c" {
		t.Fatalf("expected [a b c], got %v", keys)
	}
}

// --- Env ---

func TestEnv_Set(t *testing.T) {
	t.Setenv("RGO_TEST_VAR", "hello")
	if Env("RGO_TEST_VAR", "fallback") != "hello" {
		t.Fatal("expected 'hello'")
	}
}

func TestEnv_Fallback(t *testing.T) {
	os.Unsetenv("RGO_TEST_MISSING")
	if Env("RGO_TEST_MISSING", "default") != "default" {
		t.Fatal("expected 'default'")
	}
}
