package helpers

import "os"

// Env reads an env variable with a fallback default.
func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
