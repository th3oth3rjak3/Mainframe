package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
)

type RoleRepository interface {
	// GetAll returns all roles or an error.
	GetAll() ([]domain.Role, error)
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
		log.Err(err).Msg("failed to get all roles")
		return nil, fmt.Errorf("role repository error: %w", err)
	}

	return roles, nil
}
