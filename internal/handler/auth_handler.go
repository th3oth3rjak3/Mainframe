package handler

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
)

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (lr *LoginRequest) Validate() []string {
	var errs []string

	if strings.TrimSpace(lr.Username) == "" {
		errs = append(errs, "username is required")
	}

	if strings.TrimSpace(lr.Password) == "" {
		errs = append(errs, "password is required")
	}

	return errs
}

// LoginResponse represents the login response
type LoginResponse struct {
	Message  string `json:"message" example:"Login successful"`
	Username string `json:"username" example:"admin"`
	Email    string `json:"email"`
}

// HandleLogin logs a user into the application.
//
// @Summary      Login user
// @Description  Authenticate user credentials
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200 {object} LoginResponse
// @Router       /api/auth/login [post]
func HandleLogin(c echo.Context, userRepo repository.UserRepository, pwHasher domain.PasswordHasher) error {
	const INVALID_MESSAGE = "invalid username or password"
	const INTERNAL_SERVER_ERROR = "something bad happened"

	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return JsonError(c, "invalid request", nil, 400)
	}

	errs := req.Validate()

	if len(errs) > 0 {
		return JsonError(c, "invalid request", errs, 400)
	}

	var user *domain.User
	user, err := userRepo.GetByUsername(req.Username)

	if err != nil {
		log.Error().Err(err).Msg("failed to get by username")
		return JsonError(c, INTERNAL_SERVER_ERROR, nil, 500)
	}

	if user == nil {
		// Do an arbitrary comparison to prevent timing attacks
		err = pwHasher.FakeVerify(req.Password)

		if err != nil {
			log.Error().Err(err).Msg("failed to compare password and hash for nil user")
			return JsonError(c, INTERNAL_SERVER_ERROR, nil, 500)
		}

		return JsonError(c, INVALID_MESSAGE, nil, 401)
	}

	match, err := pwHasher.Verify(req.Password, user.PasswordHash)
	if err != nil {
		log.Error().Err(err).Msg("failed to compare password and hash for user")
		return JsonError(c, INTERNAL_SERVER_ERROR, nil, 500)
	}

	if !match {
		return JsonError(c, INVALID_MESSAGE, nil, 401)
	}

	// TODO: User has a valid login, add session, make cookie, and then return details for the client.
	response := LoginResponse{
		Message:  "Login received",
		Username: req.Username,
		Email:    user.Email,
	}

	return c.JSON(200, response)
}
