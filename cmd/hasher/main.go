package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	_ "github.com/th3oth3rjak3/mainframe/internal/logger"
)

// This CLI, is used to generate argon2id hashes for inserting an initial administrator
// into the database in a secure way.
func main() {
	if len(os.Args) < 2 {
		log.Fatal().Msg("Usage: hasher <password>")
	}

	hasher := domain.NewPasswordHasher()
	password := os.Args[1]

	hash, err := hasher.HashPassword(password)
	if err != nil {
		log.Err(err).Msg("error occurred during hashing")
		return
	}

	log.Info().Msg(hash)
}
