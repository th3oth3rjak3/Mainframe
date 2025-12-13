package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type SessionCleanupService interface {
	DeleteExpired() error
}

type sessionCleanupService struct {
	db *sqlx.DB
}

func NewSessionCleanupService(db *sqlx.DB) SessionCleanupService {
	return &sessionCleanupService{db: db}
}

func (s *sessionCleanupService) DeleteExpired() error {
	query := "DELETE FROM sessions WHERE expires_at < ?"
	result, err := s.db.Exec(query, time.Now().UTC())
	if err != nil {
		log.Err(err).Msg("failed to execute cleanup command")
		return fmt.Errorf("failed to execute delete session command: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		log.Err(err).Msg("failed to get number of rows affected")
		return fmt.Errorf("unhandled error: %w", err)
	}

	log.Info().Int64("rows_affected", affected).Msg("expired sessions deleted")
	return nil
}

func RunSessionCleanupJob(ctx context.Context, service SessionCleanupService) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	log.Info().Msg("starting background session cleanup")

	if err := service.DeleteExpired(); err != nil {
		log.Err(err).Msg("session cleanup failed")
	}

	for {
		select {
		case <-ticker.C:
			if err := service.DeleteExpired(); err != nil {
				log.Err(err).Msg("session cleanup failed")
			}

		case <-ctx.Done():
			log.Info().Msg("background session cleanup service stopping")
			return
		}
	}
}
