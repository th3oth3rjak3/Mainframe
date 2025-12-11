package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
)

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Message  string `json:"message" example:"Login successful"`
	Username string `json:"username" example:"admin"`
	Email    string `json:"email"`
}

// @Summary      Login user
// @Description  Authenticate user credentials
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200 {object} LoginResponse
// @Router       /api/auth/login [post]
func HandleLogin(c echo.Context, userRepo repository.UserRepository) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{
			"error": "invalid request",
		})
	}

	var user *domain.User
	user, err := userRepo.GetByUsername(req.Username)

	if err != nil {
		return c.String(500, err.Error())
	}

	response := LoginResponse{
		Message:  "Login received",
		Username: req.Username,
		Email:    user.Email,
	}

	return c.JSON(200, response)
}
