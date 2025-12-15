package services

import (
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
	"github.com/th3oth3rjak3/mainframe/internal/shared"
)

type RoleService interface {
	GetAllRoles(user *domain.User) ([]domain.Role, error)
}

type roleService struct {
	roleRepository repository.RoleRepository
}

func NewRoleService(roleRepository repository.RoleRepository) RoleService {
	return &roleService{roleRepository: roleRepository}
}

func (r *roleService) GetAllRoles(user *domain.User) ([]domain.Role, error) {
	if !user.HasRole(domain.Administrator) {
		return nil, shared.ErrForbidden
	}

	return r.roleRepository.GetAll()
}
