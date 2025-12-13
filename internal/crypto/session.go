package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

// GenerateRandomBytes returns n cryptographically secure random bytes.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ComputeHMACSHA256 returns a base64.RawURLEncoding HMAC of the message using the key.
func ComputeHMACSHA256(message, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(message)
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// EncodeSessionToken combines sessionID and verifier into a single string for storage in cookies.
func EncodeSessionToken(sessionID string, verifier []byte) string {
	return sessionID + ":" + base64.RawURLEncoding.EncodeToString(verifier)
}

// DecodeSessionToken splits a token into sessionID and verifier bytes.
func DecodeSessionToken(token string) (sessionID string, verifier []byte, err error) {
	parts := strings.SplitN(token, ":", 2)
	if len(parts) != 2 {
		return "", nil, fmt.Errorf("invalid token format")
	}
	sessionID = parts[0]
	verifier, err = base64.RawURLEncoding.DecodeString(parts[1])
	return
}

// VerifyVerifier checks if the provided token matches the stored HMAC, in constant time.
func VerifyVerifier(token []byte, key []byte, expectedHMAC string) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(token)
	expectedMAC, err := base64.RawURLEncoding.DecodeString(expectedHMAC)
	if err != nil {
		return false
	}
	return hmac.Equal(mac.Sum(nil), expectedMAC)
}
