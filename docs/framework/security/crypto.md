---
title: "Crypto & Security Utilities"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Crypto & Security Utilities

## Abstract

This document covers the built-in cryptographic utilities in
`core/crypto/` — random generation, hashing, HMAC signing, and
AES-256-GCM encryption/decryption — plus the bcrypt password helpers
in `app/helpers/`.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Random Generation](#2-random-generation)
3. [Hashing](#3-hashing)
4. [HMAC Signing](#4-hmac-signing)
5. [AES-256-GCM Encryption](#5-aes-256-gcm-encryption)
6. [Password Hashing (bcrypt)](#6-password-hashing-bcrypt)
7. [Security Considerations](#7-security-considerations)
8. [References](#8-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **HMAC** — Hash-based Message Authentication Code; used for
  verifying message integrity and authenticity.
- **GCM** — Galois/Counter Mode; an authenticated encryption mode
  for AES.

## 2. Random Generation

All random generation uses `crypto/rand` (cryptographically secure):

```go
// RandomBytes returns n cryptographically random bytes.
func RandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return nil, err
    }
    return b, nil
}

// RandomHex returns a random hex string of n bytes (2n hex chars).
func RandomHex(n int) string {
    b, _ := RandomBytes(n)
    return hex.EncodeToString(b)
}

// RandomBase64 returns a URL-safe base64-encoded random string.
func RandomBase64(n int) string {
    b, _ := RandomBytes(n)
    return base64.URLEncoding.EncodeToString(b)
}
```

Usage:

```go
token := crypto.RandomHex(32)       // 64-char hex string
apiKey := crypto.RandomBase64(24)   // URL-safe base64 string
```

## 3. Hashing

SHA-256 hash for non-password data:

```go
func SHA256Hash(data string) string {
    h := sha256.Sum256([]byte(data))
    return hex.EncodeToString(h[:])
}
```

Usage:

```go
hash := crypto.SHA256Hash("sensitive-data")
```

## 4. HMAC Signing

Create and verify HMAC-SHA256 signatures (e.g., webhook
verification):

```go
func HMACSign(message, key string) string {
    mac := hmac.New(sha256.New, []byte(key))
    mac.Write([]byte(message))
    return hex.EncodeToString(mac.Sum(nil))
}

func HMACVerify(message, signature, key string) bool {
    expected := HMACSign(message, key)
    return hmac.Equal([]byte(expected), []byte(signature))
}
```

Usage:

```go
sig := crypto.HMACSign(payload, secret)
valid := crypto.HMACVerify(payload, sig, secret)
```

`HMACVerify` uses constant-time comparison via `hmac.Equal` to
prevent timing attacks.

## 5. AES-256-GCM Encryption

Symmetric encryption for sensitive data (cookie values, tokens):

```go
func Encrypt(plaintext string, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encrypted string, key []byte) (string, error) {
    data, err := base64.URLEncoding.DecodeString(encrypted)
    if err != nil {
        return "", err
    }
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return "", errors.New("ciphertext too short")
    }
    plaintext, err := gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
    if err != nil {
        return "", err
    }
    return string(plaintext), nil
}
```

Usage:

```go
key := []byte(os.Getenv("APP_KEY")) // must be 32 bytes
enc, _ := crypto.Encrypt("user@example.com", key)
plain, _ := crypto.Decrypt(enc, key)
```

## 6. Password Hashing (bcrypt)

Located in `app/helpers/`:

```go
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(hash, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
```

Usage:

```go
hash, _ := helpers.HashPassword("secret123")
valid := helpers.CheckPassword(hash, "secret123") // true
```

## 7. Security Considerations

- `APP_KEY` **MUST** be exactly 32 bytes for AES-256 and **MUST NOT**
  be committed to version control.
- All random generation **MUST** use `crypto/rand`, never
  `math/rand`.
- HMAC verification **MUST** use `hmac.Equal` (constant-time) to
  prevent timing side-channel attacks.
- Nonces **MUST** be unique per encryption operation — the
  implementation generates them from `crypto/rand`.
- bcrypt is intentionally slow — do not use SHA-256 for password
  hashing.

## 8. References

- [Authentication](authentication.md)
- [Sessions](sessions.md)
- [Configuration](../core/configuration.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
