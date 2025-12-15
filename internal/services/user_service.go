package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
	"github.com/th3oth3rjak3/mainframe/internal/shared"
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

	// Create makes a new user from the provided request. The ID of the new
	// user will be returned upon success.
	Create(user *domain.User, request domain.UserCreate) (uuid.UUID, error)
}

func NewUserService(
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	pwHasher domain.PasswordHasher,
) UserService {
	return &userService{
		userRepository: userRepository,
		roleRepository: roleRepository,
		passwordHasher: pwHasher,
	}
}

type userService struct {
	userRepository repository.UserRepository
	roleRepository repository.RoleRepository
	passwordHasher domain.PasswordHasher
}

func (s *userService) GetAll(user *domain.User) ([]domain.UserRead, error) {
	if user == nil || !user.HasRole(domain.Administrator) {
		return nil, shared.ErrForbidden
	}

	users, err := s.userRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	userList := make([]domain.UserRead, len(users))

	for idx, user := range users {
		userList[idx] = domain.NewUserRead(&user)
	}

	return userList, nil
}

func (s *userService) GetByID(user *domain.User, userID uuid.UUID) (*domain.UserRead, error) {
	if user == nil || !user.HasRole(domain.Administrator) {
		return nil, shared.ErrForbidden
	}

	foundUser, err := s.userRepository.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	userRead := domain.NewUserRead(foundUser)
	return &userRead, nil
}

func (s *userService) Create(user *domain.User, request domain.UserCreate) (uuid.UUID, error) {
	if user == nil || !user.HasRole(domain.Administrator) {
		return uuid.UUID{}, shared.ErrForbidden
	}

	// Ensure the unique username constraint in the database is not violated
	existing, err := s.userRepository.GetByUsername(request.Username)
	if err != nil {
		return uuid.UUID{}, err
	}

	if existing != nil {
		return uuid.UUID{}, shared.ErrUsernameTaken
	}

	pwHash, err := s.passwordHasher.HashPassword(request.Password)
	if err != nil {
		return uuid.UUID{}, err
	}

	roles := make([]domain.Role, 1)

	role, err := s.roleRepository.GetByName(domain.BasicUser)
	if err != nil {
		return uuid.UUID{}, err
	}

	roles[0] = *role

	newUser, err := domain.NewUser(request.Username, request.Email, request.FirstName, request.LastName, pwHash, roles)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = s.userRepository.Create(newUser)
	if err != nil {
		return uuid.UUID{}, err
	}

	return newUser.ID, nil
}
