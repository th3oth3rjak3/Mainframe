package services

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
)

type CookieService interface {
	ClearCookie(c echo.Context)
	SetCookie(c echo.Context, session *domain.Session, rawToken []byte)
	ParseSessionCookie(cookie *http.Cookie) (sessionID string, rawToken []byte, err error)
}

type cookieService struct{}

func NewCookieService() CookieService {
	return &cookieService{}
}

func (s *cookieService) ClearCookie(c echo.Context) {
	c.SetCookie(getEmptyCookie())
}

func (s *cookieService) SetCookie(c echo.Context, session *domain.Session, rawToken []byte) {
	token := base64.RawURLEncoding.EncodeToString(rawToken)
	c.SetCookie(createSessionCookie(session, token))
}

func (s *cookieService) ParseSessionCookie(cookie *http.Cookie) (sessionID string, rawToken []byte, err error) {
	parts := strings.SplitN(cookie.Value, ":", 2)
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
func createSessionCookie(session *domain.Session, rawToken string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "session_id"
	cookie.Value = fmt.Sprintf("%s:%s", session.ID.String(), rawToken)
	cookie.Expires = session.ExpiresAt
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true // For production
	cookie.SameSite = http.SameSiteStrictMode
	return cookie
}

// getEmptyCookie creates an empty cookie that is used to replace the existing one
// in the browser.
func getEmptyCookie() *http.Cookie {
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
