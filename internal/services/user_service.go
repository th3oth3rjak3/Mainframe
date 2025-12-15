package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
)

type UserService interface {
	// GetAll returns all users or an error. There aren't assumed
	// to be more than 25 users for this application, so paging
	// isn't necessary. The supplied user is the one performing the
	// request and must have the correct role to access this resource.
	GetAll(user *domain.User) ([]domain.UserRead, error)

	// GetByID gets a user by their ID. If the user is not found
	// no error will be returned and the user will be nil.
	// The supplied user must have the correct role to access this
	// protected resource.
	GetByID(user *domain.User, userID uuid.UUID) (*domain.UserRead, error)
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

func (s *userService) GetByID(user *domain.User, userID uuid.UUID) (*domain.UserRead, error) {
	if user == nil || !user.HasRole(domain.Administrator) {
		return nil, ErrForbidden
	}

	foundUser, err := s.userRepository.GetByID(userID)
	if err != nil {
		log.Err(err).
			Str("user_ID", userID.String()).
			Msg("repository error getting by id")

		return nil, fmt.Errorf("user repository error: %w", err)
	}

	userRead := domain.NewUserRead(foundUser)
	return &userRead, nil
}
