package middleware

import (
	"net/http"

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
				log.Error().Msg("user missing from context where expected in require role middleware")
				return shared.InternalServerError(c)
			}

			if !user.HasRole(roleName) {
				log.Warn().
					Str("user_id", user.ID.String()).
					Str("path", c.Request().URL.Path).
					Msg("user attempting to access forbidden resources")

				return shared.JsonError(c, "forbidden", nil, http.StatusForbidden)
			}

			return next(c)
		}
	}
}
