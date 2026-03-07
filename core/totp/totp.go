package totp

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

// Key wraps an OTP key with its secret and provisioning URI.
type Key struct {
	Secret string // Base32-encoded secret
	URL    string // otpauth:// provisioning URI for QR codes
}

// GenerateKey creates a new TOTP secret for the given issuer and account.
// The issuer appears in the authenticator app (e.g. "RapidGo").
// The account identifies the user (typically their email).
func GenerateKey(issuer, account string) (*Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, err
	}
	return &Key{
		Secret: key.Secret(),
		URL:    key.URL(),
	}, nil
}

// ValidateCode checks whether a TOTP code is valid for the given secret.
// Uses a time window of ±1 period (30 seconds) for clock drift tolerance.
func ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

// GenerateBackupCodes creates n random backup codes in "XXXX-XXXX" format.
func GenerateBackupCodes(n int) ([]string, error) {
	if n <= 0 {
		return []string{}, nil
	}
	codes := make([]string, 0, n)
	for i := 0; i < n; i++ {
		b := make([]byte, 4)
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}
		codes = append(codes, fmt.Sprintf("%04X-%04X", int(b[0])<<8|int(b[1]), int(b[2])<<8|int(b[3])))
	}
	return codes, nil
}

// HashBackupCode returns a bcrypt hash of a backup code.
// The code is normalized to uppercase before hashing.
func HashBackupCode(code string) (string, error) {
	normalized := strings.ToUpper(code)
	hash, err := bcrypt.GenerateFromPassword([]byte(normalized), bcrypt.DefaultCost)
	return string(hash), err
}

// VerifyBackupCode checks a plaintext code against a bcrypt hash.
// The code is normalized to uppercase before comparison.
func VerifyBackupCode(code, hash string) bool {
	normalized := strings.ToUpper(code)
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(normalized)) == nil
}
