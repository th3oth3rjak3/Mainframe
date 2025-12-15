package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
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
		return fmt.Errorf("the request body is malformed or invalid: %w", shared.ErrBadRequest)
	}

	if err := req.Validate(); err != nil {
		return err
	}

	result, err := authService.Login(&req)
	if err != nil {
		return err
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
		return fmt.Errorf("could not get session from context")
	}

	if err := authService.Logout(session); err != nil {
		return err
	}

	cookieService.ClearCookie(c)
	c.NoContent(http.StatusNoContent)

	return nil
}
