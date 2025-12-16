package services

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/th3oth3rjak3/mainframe/internal/crypto"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
	"github.com/th3oth3rjak3/mainframe/internal/shared"
)

const maximumLoginAttemptsAllowed = 5

type AuthenticationService interface {
	Login(request *domain.LoginRequest) (*LoginResult, error)
	Logout(session *domain.Session) error
}

func NewAuthenticationService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	pwHasher domain.PasswordHasher,
	hmacKey string,
) AuthenticationService {
	return &authenticationService{
		userRepository:    userRepo,
		sessionRepository: sessionRepo,
		passwordHasher:    pwHasher,
		hmacKey:           hmacKey,
	}
}

type LoginResult struct {
	User            *domain.User
	Session         *domain.Session
	RawSessionToken []byte
}

type authenticationService struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	passwordHasher    domain.PasswordHasher
	hmacKey           string
}

func (s *authenticationService) Login(request *domain.LoginRequest) (*LoginResult, error) {
	var user *domain.User
	user, err := s.userRepository.GetByUsername(request.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get by username: %w", err)
	}

	if user == nil {
		_ = s.passwordHasher.FakeVerify(request.Password) // Prevent timing attack
		return nil, shared.ErrInvalidCredentials
	}

	match, err := s.passwordHasher.Verify(request.Password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("password verification failed: %w", err)
	}

	if !match {
		err := s.handleFailedLogin(user)
		return nil, err
	}

	if err = s.handleSuccessfulLogin(user); err != nil {
		return nil, err
	}

	session, verifier, err := s.createSessionForUser(user)
	if err != nil {
		return nil, err
	}

	return &LoginResult{User: user, Session: session, RawSessionToken: verifier}, nil
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
		return fmt.Errorf("failed to update user after failed login: %w", err)
	}

	// logging ok for critical business and security events.
	log.Warn(fmt.Sprintf("invalid login attempt #%d for user %s", user.FailedLoginAttempts, user.Username))
	return shared.ErrInvalidCredentials
}

func (s *authenticationService) handleSuccessfulLogin(user *domain.User) error {
	now := time.Now().UTC()
	user.FailedLoginAttempts = 0
	user.LastFailedLoginAttempt = nil
	user.UpdatedAt = now
	user.LastLogin = &now

	err := s.userRepository.UpdateBasic(user)
	if err != nil {
		return fmt.Errorf("failed to update user after successful login: %w", err)
	}

	return nil
}

func (s *authenticationService) createSessionForUser(user *domain.User) (*domain.Session, []byte, error) {
	verifier, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate session verifier: %w", err)
	}

	serverKey := []byte(s.hmacKey)
	token := crypto.ComputeHMACSHA256(verifier, serverKey)

	session, err := domain.NewSession(user.ID, token)
	if err != nil {
		return nil, nil, fmt.Errorf("session creation failed: %w", err)
	}

	err = s.sessionRepository.Create(session)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save session: %w", err)
	}

	return session, verifier, nil
}

func (s *authenticationService) Logout(session *domain.Session) error {
	err := s.sessionRepository.DeleteByID(session.ID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
