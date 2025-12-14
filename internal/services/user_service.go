package services

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
)

type UserService interface {
	GetAll(user *domain.User) ([]domain.UserRead, error)
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

func (s *userService) GetAll(user *domain.User) ([]domain.UserRead, error) {
	if !user.HasRole(domain.Administrator) {
		return nil, ErrForbidden
	}

	users, err := s.userRepository.GetAll()
	if err != nil {
		log.Err(err).Msg("user repository error while getting all users")
		return nil, fmt.Errorf("user repository error: %w", err)
	}

	userList := make([]domain.UserRead, len(users))

	for idx, user := range users {
		userList[idx] = domain.NewUserRead(&user)
	}

	return userList, nil
}
