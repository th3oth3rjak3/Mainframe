package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
)

// A private key type to prevent collisions in context
type contextKey string

// UserContextKey is the key used to store the user object in the request context.
const UserContextKey = contextKey("user")

// sessionDuration defines how long a session is valid for after the last activity.
const sessionDuration = 2 * time.Hour

// AuthMiddleware holds the dependencies for our authentication middleware.
type AuthMiddleware struct {
	sessionRepo repository.SessionRepository
	userRepo    repository.UserRepository
}

// NewAuthMiddleware creates a new instance of our AuthMiddleware.
func NewAuthMiddleware(sessionRepo repository.SessionRepository, userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}

// SessionAuth is the actual middleware function.
func (m *AuthMiddleware) SessionAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1. Check for the session cookie.
		cookie, err := c.Cookie("session_id")
		if err != nil {
			// If the cookie is not found, it's a standard unauthorized error.
			return echo.NewHTTPError(http.StatusUnauthorized, "missing session cookie")
		}
		sessionIDString := cookie.Value
		sessionID, err := uuid.Parse(sessionIDString)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "session_id invalid format")
		}

		// 2. Validate the session from the database.
		session, err := m.sessionRepo.GetByID(sessionID)
		if err != nil {
			// This is a server error (e.g., DB down), not an auth failure.
			return echo.NewHTTPError(http.StatusInternalServerError, "could not verify session")
		}

		// 3. Check if session is found and not expired.
		if session == nil || session.ExpiresAt.Before(time.Now().UTC()) {
			// If session is expired or not found, clear the bad cookie and deny access.
			c.SetCookie(clearSessionCookie())
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired session")
		}

		// 4. Implement Sliding Window: Update the expiration time.
		newExpiration := time.Now().UTC().Add(sessionDuration)
		session.ExpiresAt = newExpiration
		if err := m.sessionRepo.Update(session); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not update session")
		}

		// Create and set a new cookie with the updated expiration.
		newCookie := createHttpCookie(session)
		c.SetCookie(newCookie)

		// 5. Fetch the user associated with the valid session.
		user, err := m.userRepo.GetByID(session.UserID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not fetch user")
		}
		if user == nil {
			// This is an edge case: the session is valid but the user was deleted.
			return echo.NewHTTPError(http.StatusUnauthorized, "user for session not found")
		}

		// 6. Attach the user object to the context for downstream handlers.
		ctx := context.WithValue(c.Request().Context(), UserContextKey, user)
		c.SetRequest(c.Request().WithContext(ctx))

		// 7. Proceed to the next handler in the chain.
		return next(c)
	}
}

// Helper function to create a session cookie.
func createHttpCookie(session *domain.Session) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "session_id"
	cookie.Value = session.ID.String()
	cookie.Expires = session.ExpiresAt
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true // For production
	cookie.SameSite = http.SameSiteStrictMode
	return cookie
}

// Helper function to create a cookie that clears the session from the browser.
func clearSessionCookie() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "session_id"
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0) // Set to a time in the past
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	return cookie
}
