package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/th3oth3rjak3/mainframe/internal/services"
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
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	roles, err := roleService.GetAllRoles(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, roles)
}
