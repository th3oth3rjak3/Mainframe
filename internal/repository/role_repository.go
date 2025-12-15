package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/shared"
)

type RoleRepository interface {
	// GetAll returns all roles or an error.
	GetAll() ([]domain.Role, error)

	// GetByName searches for a role by its name. If not found,
	// then an error will be returned.
	GetByName(name string) (*domain.Role, error)
}

func NewRoleRepository(db *sqlx.DB) RoleRepository {
	return &roleRepository{DB: db}
}

type roleRepository struct {
	DB *sqlx.DB
}

func (r *roleRepository) GetAll() ([]domain.Role, error) {
	var roles []domain.Role

	query := "SELECT id, name FROM roles"

	err := r.DB.Select(&roles, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles from database: %w", err)
	}

	return roles, nil
}

func (r *roleRepository) GetByName(name string) (*domain.Role, error) {
	var role domain.Role
	query := "SELECT id, name FROM roles WHERE name = ?"
	err := r.DB.Get(&role, query, name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, shared.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}

	return &role, nil
}
