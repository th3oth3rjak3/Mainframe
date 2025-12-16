package handler

import (
	"github.com/gofiber/fiber/v2"
)

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
func HandleHealthCheck(c *fiber.Ctx) error {
	return c.JSON(HealthCheckResponse{Status: "ok"})
}
