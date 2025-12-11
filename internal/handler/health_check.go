package handler

import "github.com/labstack/echo/v4"

// @Summary      Health Check
// @Description  Perform Health Check
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health [get]
func HandleHealthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "ok",
	})
}
