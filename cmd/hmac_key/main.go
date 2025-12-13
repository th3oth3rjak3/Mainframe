package main

import (
	"encoding/base64"

	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/crypto"
	_ "github.com/th3oth3rjak3/mainframe/internal/logger"
)

// generate a new server key used for hmac session token hashing
func main() {
	bytes, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to generate server key")
	}

	encodedKey := base64.RawURLEncoding.EncodeToString(bytes)
	log.Info().Msg(encodedKey)
}
