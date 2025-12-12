package services

import (
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
)

type UserService interface {
}

func NewUserService(
	userRepository repository.UserRepository,
	pwHasher domain.PasswordHasher,
) UserService {
	return &userService{
		userRepository: userRepository,
		passwordHasher: pwHasher,
	}
}

type userService struct {
	userRepository repository.UserRepository
	passwordHasher domain.PasswordHasher
}
