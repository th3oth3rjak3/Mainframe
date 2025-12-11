package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// Pretty human-readable console output
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// Configure the global logger
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	// If you want debug logs visible, uncomment this:
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
