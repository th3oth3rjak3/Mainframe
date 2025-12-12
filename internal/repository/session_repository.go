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
		SELECT id, user_id, expires_at
		FROM sessions
		WHERE id = ?`

	err := r.db.Get(&session, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sqliteSessionRepository) Create(session *domain.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES (?, ?, ?)
	`

	rows, err := r.db.Exec(query, session.ID, session.UserID, session.ExpiresAt)
	if err != nil {
		return err
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return fmt.Errorf("expected to create 1 session, rows affected: %d", rows)
	}

	return nil
}

func (r *sqliteSessionRepository) DeleteByID(id uuid.UUID) error {
	query := `
		DELETE FROM sessions
		WHERE id = ?
	`

	rows, err := r.db.Exec(query, id)
	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		return err
	}

	log.Info().Msg(fmt.Sprintf("deleted %d session rows for user %s", rows, id))

	return nil
}
