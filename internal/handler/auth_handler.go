package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	mw "github.com/th3oth3rjak3/mainframe/internal/middleware"
	"github.com/th3oth3rjak3/mainframe/internal/services"
	"github.com/th3oth3rjak3/mainframe/internal/shared"
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
	cookieService services.CookieService,
) error {

	var req domain.LoginRequest

	if err := c.Bind(&req); err != nil {
		return shared.JsonError(c, "invalid request", nil, 400)
	}

	result, err := authService.Login(&req)
	if err != nil {
		return handleServiceErrors(c, err)
	}

	cookieService.SetCookie(c, result.Session, result.RawSessionToken)

	response := domain.NewLoginResponse(result.User)
	return c.JSON(200, response)
}

// HandleLogout logs a user out of the application.
//
// @Summary      Logout user
// @Description  Log out of application
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Success      204
// @Router       /api/auth/logout [post]
func HandleLogout(
	c echo.Context,
	authService services.AuthenticationService,
	cookieService services.CookieService,
) error {
	session, ok := c.Request().Context().Value(mw.SessionContextKey).(*domain.Session)
	if !ok {
		return handleServiceErrors(c, fmt.Errorf("could not get session from context"))
	}

	if err := authService.Logout(session); err != nil {
		return handleServiceErrors(c, err)
	}

	cookieService.ClearCookie(c)
	c.NoContent(http.StatusNoContent)

	return nil
}

func handleServiceErrors(c echo.Context, err error) error {
	var validationError *services.ValidationError

	if errors.As(err, &validationError) {
		return shared.JsonError(c, validationError.Message, validationError.Details, http.StatusBadRequest)
	}

	if errors.Is(err, services.ErrUnauthorized) {
		return shared.JsonError(c, err.Error(), nil, http.StatusUnauthorized)
	}

	if errors.Is(err, services.ErrForbidden) {
		return shared.JsonError(c, err.Error(), nil, http.StatusForbidden)
	}

	log.Error().
		Err(err).
		Str("path", c.Path()).
		Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).
		Msg("Unhandled internal error caught by handler")

	return shared.InternalServerError(c)
}
