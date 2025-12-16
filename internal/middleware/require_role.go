package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/shared"
)

func RequireRole(roleName string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals(UserContextKey).(*domain.User)
		if !ok || user == nil {
			return fmt.Errorf("user expected in context but was not found")
		}

		if !user.HasRole(roleName) {
			msg := fmt.Sprintf("user with id '%s' trying access forbidden resources", user.ID.String())
			path := c.OriginalURL()
			log.Warnf("%s; path=%s", msg, path)
			return shared.ErrForbidden
		}

		return c.Next()
	}
}
