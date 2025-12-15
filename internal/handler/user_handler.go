package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
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
		return err
	}

	users, err := userService.GetAll(user)
	if err != nil {
		return err
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
		return err
	}

	idString := c.Param("id")
	userID, err := uuid.Parse(idString)

	if err != nil {
		return fmt.Errorf("the id parameter was malformed or invalid: %w", shared.ErrBadRequest)
	}

	foundUser, err := userService.GetByID(user, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, foundUser)
}

// HandleCreateUser creates a new user.
//
// @Summary      Create User
// @Description  Create a new application user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Param        request body domain.UserCreate true "New User"
// @Router       /api/users [post]
func HandleCreateUser(c echo.Context, userService services.UserService) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var request domain.UserCreate

	err = c.Bind(&request)
	if err != nil {
		return fmt.Errorf("the request body is malformed or invalid: %w", shared.ErrBadRequest)
	}

	err = request.Validate()
	if err != nil {
		return err
	}

	id, err := userService.Create(user, request)
	if err != nil {
		return err
	}

	locationUrl := fmt.Sprintf("api/users/%s", id)
	c.Response().Header().Set("Location", locationUrl)
	return c.JSON(http.StatusCreated, map[string]string{"id": id.String()})
}
