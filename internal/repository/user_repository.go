package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/shared"
)

type UserRepository interface {
	// Fetch a user by ID, when not found, returns an error.
	GetByID(id uuid.UUID) (*domain.User, error)

	// Fetch a user by Username, when not found returns an error.
	GetByUsername(username string) (*domain.User, error)

	// Get all users.
	//
	// There are never expected to be more than 25 users since
	// this is a local application. If that changes, we can
	// rewrite this to produce paged results.
	GetAll() ([]domain.User, error)

	// Create a new user.
	Create(user *domain.User) error

	// Update an existing user's basic details, does not update
	// collection objects like roles.
	UpdateBasic(user *domain.User) error

	// Delete an existing user and all of the associated data.
	// This is unrecoverable.
	Delete(user *domain.User) error
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
		return nil, shared.ErrNotFound // not found
	}

	if err != nil {
		return nil, err
	}

	roles, err := r.getRolesForUser(user.ID)
	if err != nil {
		return nil, err
	}

	user.Roles = roles

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
		return nil, shared.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	roles, err := r.getRolesForUser(user.ID)
	if err != nil {
		return nil, err
	}

	user.Roles = roles

	return &user, nil
}

func (r *sqliteUserRepository) GetAll() ([]domain.User, error) {
	var users []domain.User

	query := `
		SELECT id, username, email, first_name, last_name, 
			last_login, failed_login_attempts, last_failed_login_attempt, 
			is_disabled, created_at, updated_at
		FROM users
	`

	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}

	for idx, user := range users {
		roles, err := r.getRolesForUser(user.ID)
		if err != nil {
			return nil, err
		}

		user.Roles = roles
		users[idx] = user
	}

	return users, nil
}

func (r *sqliteUserRepository) Create(user *domain.User) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer tx.Rollback()

	query := `
		INSERT INTO users (id, username, email, first_name, last_name, password_hash)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(query, user.ID, user.Username, user.Email, user.FirstName, user.LastName, user.PasswordHash)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected to create 1 new user row, but rows affected was %d", rows)
	}

	for _, role := range user.Roles {
		query := "INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)"
		result, err := tx.Exec(query, user.ID, role.ID)
		if err != nil {
			return err
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if affected != 1 {
			return fmt.Errorf("expected to create 1 new user role, but affected was %d", affected)
		}
	}

	return tx.Commit()
}

func (r *sqliteUserRepository) UpdateBasic(user *domain.User) error {
	query := `
		UPDATE users SET 
			email = ?, 
			first_name = ?,
			last_name = ?,
			last_login = ?,
			failed_login_attempts = ?,
			last_failed_login_attempt = ?,
			is_disabled = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.LastLogin,
		user.FailedLoginAttempts,
		user.LastFailedLoginAttempt,
		user.IsDisabled,
		user.UpdatedAt,
		user.ID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected to update 1 user row, but rows affected was %d", rows)
	}

	return nil
}

func (r *sqliteUserRepository) Delete(user *domain.User) error {
	query := "DELETE FROM users WHERE id = ?"

	result, err := r.db.Exec(query, user.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if affected != 1 {
		return fmt.Errorf("expected to delete 1 record, rows affected: %d", affected)
	}

	return nil
}

func (r *sqliteUserRepository) getRolesForUser(userID uuid.UUID) ([]domain.Role, error) {
	var roles []domain.Role

	query := `
		SELECT r.id, r.name
		FROM roles r
		INNER JOIN user_roles ur 
			ON ur.role_id = r.id
		WHERE ur.user_id = ?
	`

	err := r.db.Select(&roles, query, userID)
	if err == sql.ErrNoRows {
		return make([]domain.Role, 0), nil
	}

	if err != nil {
		return nil, err
	}

	return roles, nil
}
