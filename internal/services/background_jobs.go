package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
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

// centerPad centers a string within a given width
func centerPad(s string, width int) string {
	if len(s) >= width {
		return s
	}
	padding := width - len(s)
	left := padding / 2
	right := padding - left
	return fmt.Sprintf("%s%s%s", spaces(left), s, spaces(right))
}

func spaces(n int) string {
	return fmt.Sprintf("%*s", n, "")
}

// LogInfo prints an info line with centered dashes
func LogInfo(msg string) {
	now := time.Now().Format("15:04:05")
	fmt.Printf("%s | %s | %s | %s | %s | %s | %s\n",
		now,
		centerPad("-", 3),  // status
		centerPad("-", 13), // latency
		centerPad("-", 9),  // IP
		centerPad("-", 3),  // method
		centerPad("-", 5),  // path
		msg,
	)
}

// LogError prints an error line with centered dashes
func LogError(msg string, err error) {
	now := time.Now().Format("15:04:05")
	fmt.Printf("%s | %s | %s | %s | %s | %s | %s: %v\n",
		now,
		centerPad("-", 3),  // status
		centerPad("-", 12), // latency
		centerPad("-", 11), // IP
		centerPad("-", 3),  // method
		centerPad("-", 5),  // path
		msg,
		err,
	)
}

func (s *sessionCleanupService) DeleteExpired() error {
	query := "DELETE FROM sessions WHERE expires_at < ?"
	result, err := s.db.Exec(query, time.Now().UTC())
	if err != nil {
		LogError("failed to execute cleanup command", err)
		return fmt.Errorf("failed to execute delete session command: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		LogError("failed to get number of rows affected", err)
		return fmt.Errorf("unhandled error: %w", err)
	}

	LogInfo(fmt.Sprintf("%d expired sessions deleted", affected))
	return nil
}

func RunSessionCleanupJob(ctx context.Context, service SessionCleanupService) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	LogInfo("starting background session cleanup")

	if err := service.DeleteExpired(); err != nil {
		LogError("session cleanup failed", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := service.DeleteExpired(); err != nil {
				LogError("session cleanup failed", err)
			}

		case <-ctx.Done():
			LogInfo("background session cleanup service stopping")
			return
		}
	}
}
