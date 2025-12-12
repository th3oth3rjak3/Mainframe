package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/services"
)

// HandleLogin logs a user into the application.
//
// @Summary      Login user
// @Description  Authenticate user credentials
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body domain.LoginRequest true "Login credentials"
// @Success      200 {object} domain.LoginResponse
// @Router       /api/auth/login [post]
func HandleLogin(
	c echo.Context,
	authService services.AuthenticationService,
) error {

	var req domain.LoginRequest

	if err := c.Bind(&req); err != nil {
		return JsonError(c, "invalid request", nil, 400)
	}

	user, session, err := authService.Login(&req)
	if err != nil {
		return handleServiceErrors(c, err)
	}

	cookie := createHttpCookie(session)
	c.SetCookie(cookie)

	response := domain.NewLoginResponse(user)
	return c.JSON(200, response)
}

func createHttpCookie(session *domain.Session) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "session_id"
	cookie.Value = session.ID.String()
	cookie.Expires = session.ExpiresAt
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode

	return cookie
}

func handleServiceErrors(c echo.Context, err error) error {
	var validationError *services.ValidationError

	if errors.As(err, &validationError) {
		return JsonError(c, validationError.Message, validationError.Details, http.StatusBadRequest)
	}

	if errors.Is(err, services.ErrUnauthorized) {
		return JsonError(c, err.Error(), nil, http.StatusUnauthorized)
	}

	if errors.Is(err, services.ErrForbidden) {
		return JsonError(c, err.Error(), nil, http.StatusForbidden)
	}

	log.Error().
		Err(err).
		Str("path", c.Path()).
		Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).
		Msg("Unhandled internal error caught by handler")

	return InternalServerError(c)
}
