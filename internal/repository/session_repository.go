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
	GetByID(id uuid.UUID) (*domain.Session, error)
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
