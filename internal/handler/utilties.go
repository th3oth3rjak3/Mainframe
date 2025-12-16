package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	mw "github.com/th3oth3rjak3/mainframe/internal/middleware"
)

func getUserFromContext(c *fiber.Ctx) (*domain.User, error) {
	user, ok := c.Locals(mw.UserContextKey).(*domain.User)
	if !ok {
		return nil, fmt.Errorf("expected user in context, but found none")
	}

	return user, nil
}
