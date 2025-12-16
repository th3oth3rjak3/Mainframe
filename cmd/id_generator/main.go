package main

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

// generate an arbitrary UUID
func main() {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("error generating new id %w", err)
		return
	}

	log.Info(id.String())
}
