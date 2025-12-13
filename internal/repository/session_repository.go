package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
)

type SessionRepository interface {
	// GetByID gets a session by its id. If a session is not found
	// then the returned session will be nil and no error will be returned.
	GetByID(id uuid.UUID) (*domain.Session, error)

	// Create saves a new session.
	Create(session *domain.Session) error

	// Update saves the session details
	Update(session *domain.Session) error

	// DeleteByID deletes the session with the given id.
	DeleteByID(id uuid.UUID) error
}

type sqliteSessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) SessionRepository {
	return &sqliteSessionRepository{db: db}
}

func (r *sqliteSessionRepository) GetByID(id uuid.UUID) (*domain.Session, error) {
	var session domain.Session

	query := `
		SELECT id, token, user_id, expires_at
		FROM sessions
		WHERE id = ?`

	err := r.db.Get(&session, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Err(err).
			Str("user_id", id.String()).
			Msg("failed to get by id")

		return nil, fmt.Errorf("session repository get by id error: %w", err)
	}

	return &session, nil
}

func (r *sqliteSessionRepository) Create(session *domain.Session) error {
	query := `
		INSERT INTO sessions (id, token, user_id, expires_at)
		VALUES (?, ?, ?, ?)
	`

	rows, err := r.db.Exec(query, session.ID, session.Token, session.UserID, session.ExpiresAt)
	if err != nil {
		log.Err(err).
			Str("session_id", session.ID.String()).
			Msg("error creating user session")

		return fmt.Errorf("create session repository error: %w", err)
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		log.Err(err).
			Msg("error getting rows affected create session")

		return fmt.Errorf("checking rows affected err: %w", err)
	}

	if affected != 1 {
		log.Error().
			Int64("rows_affected", affected).
			Msg("expected 1 row to be created")

		return fmt.Errorf("expected to create 1 session, rows affected: %d", rows)
	}

	return nil
}

func (r *sqliteSessionRepository) Update(session *domain.Session) error {
	query := `
		UPDATE sessions SET expires_at = ?
		WHERE id = ?
	`

	rows, err := r.db.Exec(query, session.ExpiresAt, session.ID)
	if err != nil {
		log.Err(err).
			Str("session_id", session.ID.String()).
			Msg("an error occurred updating the session")

		return fmt.Errorf("session repository update error: %w", err)
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		log.Err(err).
			Msg("error getting rows affected updating session")

		return fmt.Errorf("repo update rows affected error: %w", err)
	}

	if affected != 1 {
		log.Error().
			Int64("rows_affected", affected).
			Msg("expected 1 row to be updated")

		return fmt.Errorf("expected to update 1 session, rows affected: %d", affected)
	}

	return nil
}

func (r *sqliteSessionRepository) DeleteByID(id uuid.UUID) error {
	query := `
		DELETE FROM sessions
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Err(err).
			Str("session_id", id.String()).
			Msg("error deleting session by id")

		return fmt.Errorf("delete session by id error: %w", err)
	}

	return nil
}
