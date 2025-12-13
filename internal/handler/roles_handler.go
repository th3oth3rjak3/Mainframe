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

// HandleListRoles returns a list of all roles.
//
// @Summary      List Roles
// @Description  Get all roles
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Success      200 {object} []domain.Role
// @Router       /api/roles [get]
func HandleListRoles(c echo.Context, roleService services.RoleService) error {
	user, ok := c.Request().Context().Value(mw.UserContextKey).(*domain.User)
	if !ok {
		return handleRoleServiceErrors(c, fmt.Errorf("user was missing from context when expected"))
	}

	roles, err := roleService.GetAllRoles(user)
	if err != nil {
		log.Err(err).
			Msg("failed to get all roles")

		return handleRoleServiceErrors(c, err)
	}

	return c.JSON(http.StatusOK, roles)
}

func handleRoleServiceErrors(c echo.Context, err error) error {
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
