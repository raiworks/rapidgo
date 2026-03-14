package auth

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtConfig reads JWT_SECRET and JWT_EXPIRY from environment.
func jwtConfig() (secret string, expiry int, err error) {
	secret = os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", 0, errors.New("JWT_SECRET is not set")
	}
	if len(secret) < 32 {
		return "", 0, errors.New("JWT_SECRET must be at least 32 bytes")
	}

	expiry = 3600 // default 1 hour
	if v := os.Getenv("JWT_EXPIRY"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			expiry = parsed
		}
	}
	return secret, expiry, nil
}

// GenerateToken creates a signed JWT for the given user ID.
// Reads JWT_SECRET and JWT_EXPIRY (seconds) from environment.
func GenerateToken(userID uint) (string, error) {
	secret, expiry, err := jwtConfig()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(expiry) * time.Second).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateTokenFromString creates a signed JWT for the given user ID string.
// Use this for UUID or other string-based primary keys.
// Reads JWT_SECRET and JWT_EXPIRY (seconds) from environment.
func GenerateTokenFromString(userID string) (string, error) {
	secret, expiry, err := jwtConfig()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(expiry) * time.Second).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken parses and validates a JWT string.
// Returns the claims if the token is valid.
// Note: user_id claim is float64 (from GenerateToken) or string (from GenerateTokenFromString).
func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET must be at least 32 bytes")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenUnverifiable
	}

	return claims, nil
}
