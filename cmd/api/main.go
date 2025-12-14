package main

import (
	"context"
	"embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/api"
	"github.com/th3oth3rjak3/mainframe/internal/data"
	_ "github.com/th3oth3rjak3/mainframe/internal/logger"
	"github.com/th3oth3rjak3/mainframe/internal/services"
)

//go:embed web
var webAssets embed.FS

// @title           Mainframe API
// @version         1.0
// @description     Centralized Personal Productivity Application
// @host            localhost:8080
// @BasePath        /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Warn().Msg("Warning: .env file not found, relying on environment variables")
	}

	serverKey := os.Getenv("SERVER_KEY")

	// Initialize database
	db, err := data.InitDB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	container, err := api.NewServiceContainer(db, serverKey)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize service container")
	}

	server := api.NewServer(container, serverKey, webAssets)

	sessionCleanupService := services.NewSessionCleanupService(db)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err = server.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("shutdown error occurred")
		}
	}()

	go services.RunSessionCleanupJob(ctx, sessionCleanupService)

	// Wait for SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Trigger worker shutdown
	cancel()

	// Gracefully shut down Echo
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	_ = server.Shutdown(ctxShutdown)
}
