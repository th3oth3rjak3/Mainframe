package repository

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
)

type UserRepository interface {
	// Fetch a user by ID
	GetByID(id uuid.UUID) (*domain.User, error)

	// Fetch a user by Username
	GetByUsername(username string) (*domain.User, error)

	// Create a new user.
	Create(user *domain.User) error

	// Update an existing user (timestamps, failed login attempts, etc)
	Update(user *domain.User) error
}

type sqliteUserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &sqliteUserRepository{db: db}
}

func (r *sqliteUserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, username, email, first_name, last_name, password_hash,
			last_login, failed_login_attempts, last_failed_login_attempt, 
			is_disabled, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	err := r.db.Get(&user, query, id.String())
	if err == sql.ErrNoRows {
		return nil, nil // not found
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *sqliteUserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, username, email, first_name, last_name, password_hash,
			last_login, failed_login_attempts, last_failed_login_attempt,
			is_disabled, created_at, updated_at
		FROM users
		WHERE LOWER(username) = ?
	`

	err := r.db.Get(&user, query, strings.ToLower(username))
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *sqliteUserRepository) Create(user *domain.User) error {
	panic("todo")
}

func (r *sqliteUserRepository) Update(user *domain.User) error {
	panic("todo")
}
