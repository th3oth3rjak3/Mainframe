package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
	"github.com/th3oth3rjak3/mainframe/internal/services"
	"github.com/th3oth3rjak3/mainframe/internal/shared"
)

// A private key type to prevent collisions in context
type contextKey string

// UserContextKey is the key used to store the user object in the request context.
const UserContextKey = contextKey("user")

// SessionContextKey is the key used to store the session object in the request context.
const SessionContextKey = contextKey("session")

// sessionDuration defines how long a session is valid for after the last activity.
const sessionDuration = 2 * time.Hour

// AuthMiddleware holds the dependencies for our authentication middleware.
type AuthMiddleware struct {
	sessionRepo   repository.SessionRepository
	userRepo      repository.UserRepository
	cookieService services.CookieService
}

// NewAuthMiddleware creates a new instance of our AuthMiddleware.
func NewAuthMiddleware(
	sessionRepo repository.SessionRepository,
	userRepo repository.UserRepository,
	cookieService services.CookieService,
) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepo:   sessionRepo,
		userRepo:      userRepo,
		cookieService: cookieService,
	}
}

// SessionAuth is the actual middleware function.
func (m *AuthMiddleware) SessionAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the session cookie
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return shared.JsonError(c, "unauthorized", nil, http.StatusUnauthorized)
		}
		sessionIDString := cookie.Value
		sessionID, err := uuid.Parse(sessionIDString)
		if err != nil {
			log.Err(err).
				Str("session_id", sessionIDString).
				Msg("could not parse session_id as a UUID")

			return shared.JsonError(c, "unauthorized", nil, http.StatusUnauthorized)
		}

		// Find the session
		session, err := m.sessionRepo.GetByID(sessionID)
		if err != nil {
			log.Err(err).
				Str("session_id", sessionID.String()).
				Msg("failed to get session by id")

			return shared.InternalServerError(c)
		}

		// deal with expired and not found sessions
		if session == nil || session.ExpiresAt.Before(time.Now().UTC()) {
			m.cookieService.ClearCookie(c)
			return shared.JsonError(c, "unauthorized", nil, http.StatusUnauthorized)
		}

		// Update sliding expiration window
		newExpiration := time.Now().UTC().Add(sessionDuration)
		session.ExpiresAt = newExpiration
		if err := m.sessionRepo.Update(session); err != nil {
			log.Err(err).
				Str("session_id", sessionID.String()).
				Msg("failed to update session in the database")

			return shared.InternalServerError(c)
		}

		// Create and set a new cookie with the updated expiration.
		m.cookieService.SetCookie(c, session)

		// Fetch the user associated with the valid session.
		user, err := m.userRepo.GetByID(session.UserID)
		if err != nil {
			log.Err(err).
				Str("session_id", sessionID.String()).
				Str("user_id", session.UserID.String()).
				Msg("failed to get user by id")

			return shared.InternalServerError(c)
		}

		if user == nil {
			return shared.JsonError(c, "unauthorized", nil, http.StatusUnauthorized)
		}

		// Attach the user object to the context for downstream handlers.
		ctx := context.WithValue(c.Request().Context(), UserContextKey, user)
		ctx = context.WithValue(ctx, SessionContextKey, session)
		c.SetRequest(c.Request().WithContext(ctx))

		// Proceed to the next handler in the chain.
		return next(c)
	}
}
