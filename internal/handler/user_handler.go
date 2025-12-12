package handler

import "github.com/labstack/echo/v4"

// HandleListUsers returns a list of users.
//
// @Summary      List Users
// @Description  Get all users
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200 {object} string
// @Router       /api/users [get]
func HandleListUsers(c echo.Context) error {
	// TODO: actually return users list
	return c.String(200, "OK!")
}
