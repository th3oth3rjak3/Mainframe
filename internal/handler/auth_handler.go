package handler

import (
	"strings"

	"github.com/alexedwards/argon2id"

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

// @Summary      Login user
// @Description  Authenticate user credentials
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200 {object} LoginResponse
// @Router       /api/auth/login [post]
func HandleLogin(c echo.Context, userRepo repository.UserRepository) error {
	const INVALID_MESSAGE = "invalid username or password"
	const FAKE_HASH_FOR_TIMING_ATTACK_PREVENTION = "$argon2id$v=19$m=65536,t=1,p=16$5+5ObcY5s1LVbxJ/+Xwajg$EtvmraG0bszkPPJW4k3RFYy6UcXZTQahKIl7TLdJ0TE"
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
		_, err = argon2id.ComparePasswordAndHash(req.Password, FAKE_HASH_FOR_TIMING_ATTACK_PREVENTION)

		if err != nil {
			log.Error().Err(err).Msg("failed to compare password and hash for nil user")
			return JsonError(c, INTERNAL_SERVER_ERROR, nil, 500)
		}

		return JsonError(c, INVALID_MESSAGE, nil, 401)
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, user.PasswordHash)
	if err != nil {
		log.Error().Err(err).Msg("failed to compare password and hash for user")
		return JsonError(c, INTERNAL_SERVER_ERROR, nil, 500)
	}

	if !match {
		return JsonError(c, INVALID_MESSAGE, nil, 401)
	}

	response := LoginResponse{
		Message:  "Login received",
		Username: req.Username,
		Email:    user.Email,
	}

	return c.JSON(200, response)
}
