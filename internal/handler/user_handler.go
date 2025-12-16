package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
func HandleListUsers(c *fiber.Ctx, userService services.UserService) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	users, err := userService.GetAll(user)
	if err != nil {
		return err
	}

	return c.JSON(users)
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
func HandleGetUserByID(c *fiber.Ctx, userService services.UserService) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	idString := c.Params("id")
	userID, err := uuid.Parse(idString)

	if err != nil {
		return fmt.Errorf("%w: the id parameter was malformed or invalid", shared.ErrBadRequest)
	}

	foundUser, err := userService.GetByID(user, userID)
	if err != nil {
		return err
	}

	return c.JSON(foundUser)
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
func HandleCreateUser(c *fiber.Ctx, userService services.UserService) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	var request domain.UserCreate

	err = c.BodyParser(&request)
	if err != nil {
		return fmt.Errorf("%w: the request body is malformed or invalid", shared.ErrBadRequest)
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
	c.Set("Location", locationUrl)
	return c.JSON(map[string]string{"id": id.String()})
}

// HandleUpdateUser updates an existing user.
//
// @Summary      Update User
// @Description  Update an application user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      204
// @Param        request body domain.UserUpdate true "Update User"
// @Param        id path string true "User ID"
// @Router       /api/users/:id [put]
func HandleUpdateUser(c *fiber.Ctx, userService services.UserService) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	idString := c.Params("id")
	userID, err := uuid.Parse(idString)

	if err != nil {
		return fmt.Errorf("%w: the id parameter was malformed or invalid", shared.ErrBadRequest)
	}

	var request domain.UserUpdate
	err = c.BodyParser(&request)
	if err != nil {
		return fmt.Errorf("%w: the request body is malformed or invalid", shared.ErrBadRequest)
	}

	err = request.Validate()
	if err != nil {
		return err
	}

	err = userService.Update(user, userID, request)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
