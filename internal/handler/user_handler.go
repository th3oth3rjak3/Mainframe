package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
// @Success      200 {object} []domain.UserRead
// @Router       /api/users [get]
func HandleListUsers(c echo.Context, userService services.UserService) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return handleUserServiceErrors(c, err)
	}

	users, err := userService.GetAll(user)
	if err != nil {
		return handleUserServiceErrors(c, err)
	}

	return c.JSON(http.StatusOK, users)
}

// HandleGetUserByID returns a user by ID.
//
// @Summary      Get User
// @Description  Get one user by ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200 {object} domain.UserRead
// @Param        id path string true "User ID"
// @Router       /api/users/:id [get]
func HandleGetUserByID(c echo.Context, userService services.UserService) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return handleUserServiceErrors(c, err)
	}

	idString := c.Param("id")
	userID, err := uuid.Parse(idString)

	if err != nil {
		return handleUserServiceErrors(c, err)
	}

	foundUser, err := userService.GetByID(user, userID)
	if err != nil {
		return handleUserServiceErrors(c, err)
	}

	return c.JSON(http.StatusOK, foundUser)
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
