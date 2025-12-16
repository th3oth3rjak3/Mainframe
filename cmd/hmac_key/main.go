package main

import (
	"encoding/base64"

	"github.com/gofiber/fiber/v2/log"
	"github.com/th3oth3rjak3/mainframe/internal/crypto"
)

// generate a new server key used for hmac session token hashing
func main() {
	bytes, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		log.Fatalf("failed to generate server key %w", err)
	}

	encodedKey := base64.RawURLEncoding.EncodeToString(bytes)
	log.Info(encodedKey)
}
