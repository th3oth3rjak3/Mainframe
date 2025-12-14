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

// HandleListUsers returns a list of users.
//
// @Summary      List Users
// @Description  Get all users
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200 {object} string
// @Router       /api/users [get]
func HandleListUsers(c echo.Context, userService services.UserService) error {
	user, ok := c.Request().Context().Value(mw.UserContextKey).(*domain.User)
	if !ok {
		return handleUserServiceErrors(c, fmt.Errorf("expected user in context, but found none"))
	}

	users, err := userService.GetAll(user)
	if err != nil {
		return handleUserServiceErrors(c, err)
	}

	return c.JSON(http.StatusOK, users)
}

func handleUserServiceErrors(c echo.Context, err error) error {
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
