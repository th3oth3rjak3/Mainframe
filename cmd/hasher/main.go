package main

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
)

// This CLI, is used to generate argon2id hashes for inserting an initial administrator
// into the database in a secure way.
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: hasher <password>")
	}

	hasher := domain.NewPasswordHasher()
	password := os.Args[1]

	hash, err := hasher.HashPassword(password)
	if err != nil {
		log.Fatalf("error occurred during hashing %w", err)
		return
	}

	log.Info(hash)
}
