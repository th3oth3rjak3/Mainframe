package services

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
)

const maximumLoginAttemptsAllowed = 5

type AuthenticationService interface {
	Login(request *domain.LoginRequest) (*domain.User, *domain.Session, error)
	Logout(session *domain.Session) error
}

func NewAuthenticationService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	pwHasher domain.PasswordHasher,
) AuthenticationService {
	return &authenticationService{
		userRepository:    userRepo,
		sessionRepository: sessionRepo,
		passwordHasher:    pwHasher,
	}
}

type authenticationService struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	passwordHasher    domain.PasswordHasher
}

func (s *authenticationService) Login(request *domain.LoginRequest) (*domain.User, *domain.Session, error) {
	errs := request.Validate()

	if len(errs) > 0 {
		return nil, nil, NewValidationError("invalid request", errs)
	}

	var user *domain.User
	user, err := s.userRepository.GetByUsername(request.Username)

	if err != nil {
		log.Err(err).
			Str("username", request.Username).
			Msg("a critical datastore error occurred while fetching user by username")

		return nil, nil, fmt.Errorf("datastore error: %w", err)
	}

	if user == nil {
		_ = s.passwordHasher.FakeVerify(request.Password) // Prevent timing attack
		log.Warn().Msg("authentication attempt for non-existent user")
		return nil, nil, ErrUnauthorized
	}

	match, err := s.passwordHasher.Verify(request.Password, user.PasswordHash)
	if err != nil {
		log.Error().Err(err).Str("username", user.Username).Msg("password verification error")
		return nil, nil, fmt.Errorf("hasher error: %w", err)
	}

	if !match {
		err := s.handleFailedLogin(user)
		return nil, nil, err
	}

	if err = s.handleSuccessfulLogin(user); err != nil {
		return nil, nil, err
	}

	session, err := s.createSessionForUser(user)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (s *authenticationService) handleFailedLogin(user *domain.User) error {
	// update login error details
	user.FailedLoginAttempts += 1
	now := time.Now().UTC()
	user.LastFailedLoginAttempt = &now
	user.UpdatedAt = now
	user.IsDisabled = user.FailedLoginAttempts >= maximumLoginAttemptsAllowed

	// update user in database
	err := s.userRepository.UpdateBasic(user)
	if err != nil {
		log.Err(err).Msg("failed to update user last failed login attempt")
		return fmt.Errorf("datastore error: %w", err)
	}

	log.Warn().Msg(fmt.Sprintf("invalid login attempt #%d for user %s", user.FailedLoginAttempts, user.Username))
	return ErrUnauthorized
}

func (s *authenticationService) handleSuccessfulLogin(user *domain.User) error {
	now := time.Now().UTC()
	user.FailedLoginAttempts = 0
	user.LastFailedLoginAttempt = nil
	user.UpdatedAt = now
	user.LastLogin = &now

	err := s.userRepository.UpdateBasic(user)
	if err != nil {
		log.Err(err).Msg("failed to update user after successful login")
		return fmt.Errorf("datastore error: %w", err)
	}

	return nil
}

func (s *authenticationService) createSessionForUser(user *domain.User) (*domain.Session, error) {
	session, err := domain.NewSession(user.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create new session")
		return nil, fmt.Errorf("session creation error: %w", err)
	}

	err = s.sessionRepository.Create(session)
	if err != nil {
		log.Error().Err(err).Msg("failed to save new session to the database")
		return nil, fmt.Errorf("datastore error while saving session: %w", err)
	}

	return session, nil
}

func (s *authenticationService) Logout(session *domain.Session) error {
	err := s.sessionRepository.DeleteByID(session.ID)
	if err != nil {
		log.Err(err).Str("session_id", session.ID.String()).Msg("error deleting session by id")
		return fmt.Errorf("session data error: %w", err)
	}
	return nil
}
