package main

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	_ "github.com/th3oth3rjak3/mainframe/internal/logger"
)

// generate an arbitrary UUID
func main() {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Err(err).Msg("error generating new id")
		return
	}

	log.Info().Msg(id.String())
}
