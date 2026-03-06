package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomString generates a cryptographically random hex string of n bytes.
func RandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
