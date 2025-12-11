package handler

import "github.com/labstack/echo/v4"

type HealthCheckResponse struct {
	Status string `json:"status"`
}

// @Summary      Health Check
// @Description  Perform Health Check
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200 {object} HealthCheckResponse
// @Router       /health [get]
func HandleHealthCheck(c echo.Context) error {
	return c.JSON(200, HealthCheckResponse{Status: "ok"})
}
