package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	mw "github.com/th3oth3rjak3/mainframe/internal/middleware"
)

func getUserFromContext(c echo.Context) (*domain.User, error) {
	user, ok := c.Request().Context().Value(mw.UserContextKey).(*domain.User)
	if !ok {
		return nil, fmt.Errorf("expected user in context, but found none")
	}

	return user, nil
}
