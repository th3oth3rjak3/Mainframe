package main

import (
	"github.com/rs/zerolog/log"
	_ "github.com/th3oth3rjak3/mainframe/internal/logger"

	"github.com/th3oth3rjak3/mainframe/data"
	"github.com/th3oth3rjak3/mainframe/internal/api"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
)

// @title           Mainframe API
// @version         1.0
// @description     User authentication and session management
// @host            localhost:8080
// @BasePath        /
func main() {
	// Initialize database
	db, err := data.InitDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	pwHasher := domain.NewPasswordHasher()
	server := api.NewServer(userRepo, pwHasher)
	err = server.Start(":8080")
	log.Fatal().Err(err).Msg("shutting down")
}
