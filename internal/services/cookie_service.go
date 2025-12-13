package services

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
)

type CookieService interface {
	ClearCookie(c echo.Context)
	SetCookie(c echo.Context, session *domain.Session)
}

type cookieService struct{}

func NewCookieService() CookieService {
	return &cookieService{}
}

func (s *cookieService) ClearCookie(c echo.Context) {
	c.SetCookie(getEmptyCookie())
}

func (s *cookieService) SetCookie(c echo.Context, session *domain.Session) {
	c.SetCookie(createSessionCookie(session))
}

// createSessionCookie makes a new cookie and includes the session details.
func createSessionCookie(session *domain.Session) *http.Cookie {
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
