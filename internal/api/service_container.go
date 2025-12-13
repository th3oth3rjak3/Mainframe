package api

import (
	"github.com/jmoiron/sqlx"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
	"github.com/th3oth3rjak3/mainframe/internal/services"
)

type ServiceContainer struct {
	// Infrastructure
	DB             *sqlx.DB
	PasswordHasher domain.PasswordHasher

	// Repositories
	UserRepository    repository.UserRepository
	SessionRepository repository.SessionRepository
	RoleRepository    repository.RoleRepository

	// Services
	UserService           services.UserService
	AuthenticationService services.AuthenticationService
	CookieService         services.CookieService
	RoleService           services.RoleService
}

// NewServiceContainer builds and returns a new dependency container.
// This is the single place where all application components are instantiated.
func NewServiceContainer(db *sqlx.DB, hmacKey string) (*ServiceContainer, error) {
	// Infrastructure
	pwHasher := domain.NewPasswordHasher()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Services
	userService := services.NewUserService(userRepo, pwHasher)
	authService := services.NewAuthenticationService(userRepo, sessionRepo, pwHasher, hmacKey)
	cookieService := services.NewCookieService()
	roleService := services.NewRoleService(roleRepo)

	// Return the fully-built container
	return &ServiceContainer{
		DB:                    db,
		PasswordHasher:        pwHasher,
		UserRepository:        userRepo,
		RoleRepository:        roleRepo,
		SessionRepository:     sessionRepo,
		UserService:           userService,
		RoleService:           roleService,
		AuthenticationService: authService,
		CookieService:         cookieService,
	}, nil
}
