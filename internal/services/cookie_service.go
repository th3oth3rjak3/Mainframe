package services

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/shared"
)

type CookieService interface {
	ClearCookie(c *fiber.Ctx)
	SetCookie(c *fiber.Ctx, session *domain.Session, rawToken []byte)
	ParseSessionCookie(cookie string) (sessionID string, rawToken []byte, err error)
}

type cookieService struct{}

func NewCookieService() CookieService {
	return &cookieService{}
}

func (s *cookieService) ClearCookie(c *fiber.Ctx) {
	c.Cookie(getEmptyCookie())
}

func (s *cookieService) SetCookie(c *fiber.Ctx, session *domain.Session, rawToken []byte) {
	token := base64.RawURLEncoding.EncodeToString(rawToken)
	c.Cookie(createSessionCookie(session, token))
}

func (s *cookieService) ParseSessionCookie(cookie string) (sessionID string, rawToken []byte, err error) {
	parts := strings.SplitN(cookie, ":", 2)
	if len(parts) != 2 {
		return "", nil, fmt.Errorf("invalid session cookie format")
	}
	tokenBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", nil, fmt.Errorf("invalid base64 token: %w", err)
	}
	return parts[0], tokenBytes, nil
}

// createSessionCookie makes a new cookie and includes the session details.
func createSessionCookie(session *domain.Session, rawToken string) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = fmt.Sprintf("%s:%s", session.ID.String(), rawToken)
	cookie.Expires = session.ExpiresAt
	cookie.Path = "/"
	cookie.HTTPOnly = true
	cookie.Secure = shared.IsProduction() // For production
	cookie.SameSite = fiber.CookieSameSiteLaxMode
	return cookie
}

// getEmptyCookie creates an empty cookie that is used to replace the existing one
// in the browser.
func getEmptyCookie() *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0) // Set to a time in the past
	cookie.Path = "/"
	cookie.HTTPOnly = true
	cookie.Secure = shared.IsProduction()
	cookie.SameSite = fiber.CookieSameSiteLaxMode
	return cookie
}
