package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/shared"
)

func RequireRole(roleName string) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Request().Context().Value(UserContextKey).(*domain.User)
			if !ok || user == nil {
				return fmt.Errorf("user expected in context but was not found")
			}

			if !user.HasRole(roleName) {
				log.Warn().
					Str("user_id", user.ID.String()).
					Str("path", c.Request().URL.Path).
					Msg("user attempting to access forbidden resources")
				return shared.ErrForbidden
			}

			return next(c)
		}
	}
}
