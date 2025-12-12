package main

import (
	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/api"
	"github.com/th3oth3rjak3/mainframe/internal/data"
	_ "github.com/th3oth3rjak3/mainframe/internal/logger"
)

// @title           Mainframe API
// @version         1.0
// @description     Centralized Personal Productivity Application
// @host            localhost:8080
// @BasePath        /
func main() {
	// Initialize database
	db, err := data.InitDB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	container, err := api.NewServiceContainer(db)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize service container")
	}

	server := api.NewServer(container)
	err = server.Start(":8080")
	log.Fatal().Err(err).Msg("shutting down")
}
