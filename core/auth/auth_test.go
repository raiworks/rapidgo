package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TC-01: GenerateToken returns valid JWT string
func TestGenerateToken_ReturnsValidJWT(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-for-testing-32by")
	t.Setenv("JWT_EXPIRY", "3600")

	token, err := GenerateToken(1)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

// TC-02: ValidateToken parses valid token and returns claims
func TestValidateToken_ParsesValidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-for-testing-32by")
	t.Setenv("JWT_EXPIRY", "3600")

	token, err := GenerateToken(1)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		t.Fatal("expected user_id claim to be float64")
	}
	if userID != 1 {
		t.Fatalf("expected user_id=1, got %v", userID)
	}

	if _, ok := claims["exp"]; !ok {
		t.Fatal("expected exp claim")
	}
	if _, ok := claims["iat"]; !ok {
		t.Fatal("expected iat claim")
	}
}

// TC-03: ValidateToken rejects expired token
func TestValidateToken_RejectsExpiredToken(t *testing.T) {
	secret := "test-secret-key-for-testing-32by"
	t.Setenv("JWT_SECRET", secret)

	// Create token that expired 1 hour ago
	claims := jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(-1 * time.Hour).Unix(),
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to create test token: %v", err)
	}

	_, err = ValidateToken(tokenStr)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

// TC-04: ValidateToken rejects malformed token
func TestValidateToken_RejectsMalformedToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-for-testing-32by")

	_, err := ValidateToken("not.a.jwt")
	if err == nil {
		t.Fatal("expected error for malformed token")
	}
}

// TC-05: GenerateToken fails when JWT_SECRET is empty
func TestGenerateToken_FailsWithoutSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	_, err := GenerateToken(1)
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is empty")
	}
	if err.Error() != "JWT_SECRET is not set" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TC-06: ValidateToken fails when JWT_SECRET is empty
func TestValidateToken_FailsWithoutSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	_, err := ValidateToken("some.token.here")
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is empty")
	}
	if err.Error() != "JWT_SECRET is not set" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TC-07: ValidateToken rejects token signed with wrong secret
func TestValidateToken_RejectsWrongSecret(t *testing.T) {
	// Generate with one secret
	t.Setenv("JWT_SECRET", "secret-one-padded-to-32-bytes!!!")
	t.Setenv("JWT_EXPIRY", "3600")

	token, err := GenerateToken(1)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Validate with different secret
	t.Setenv("JWT_SECRET", "secret-two-padded-to-32-bytes!!!")
	_, err = ValidateToken(token)
	if err == nil {
		t.Fatal("expected error for token signed with wrong secret")
	}
}

// TC-08: GenerateToken respects JWT_EXPIRY env var
func TestGenerateToken_RespectsExpiry(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-for-testing-32by")
	t.Setenv("JWT_EXPIRY", "7200") // 2 hours

	token, err := GenerateToken(1)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("expected exp claim to be float64")
	}

	iat, ok := claims["iat"].(float64)
	if !ok {
		t.Fatal("expected iat claim to be float64")
	}

	// exp - iat should be ~7200 seconds
	diff := exp - iat
	if diff < 7190 || diff > 7210 {
		t.Fatalf("expected expiry ~7200s, got %.0f", diff)
	}
}

// TC-09: ValidateToken returns user_id as float64
func TestValidateToken_UserIDIsFloat64(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-for-testing-32by")
	t.Setenv("JWT_EXPIRY", "3600")

	token, err := GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		t.Fatalf("expected user_id to be float64, got %T", claims["user_id"])
	}
	if userID != 42 {
		t.Fatalf("expected user_id=42, got %v", userID)
	}
}

// TC-10: GenerateToken rejects secret shorter than 32 bytes
func TestGenerateToken_RejectsShortSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "too-short")

	_, err := GenerateToken(1)
	if err == nil {
		t.Fatal("expected error for short JWT_SECRET")
	}
	if err.Error() != "JWT_SECRET must be at least 32 bytes" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TC-11: ValidateToken rejects secret shorter than 32 bytes
func TestValidateToken_RejectsShortSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "too-short")

	_, err := ValidateToken("some.token.here")
	if err == nil {
		t.Fatal("expected error for short JWT_SECRET")
	}
	if err.Error() != "JWT_SECRET must be at least 32 bytes" {
		t.Fatalf("unexpected error: %v", err)
	}
}
