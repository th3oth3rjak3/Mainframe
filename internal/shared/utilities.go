package shared

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func JsonError(c echo.Context, message string, details []string, httpStatus int) error {
	return c.JSON(httpStatus, ErrorResponse{
		Message: message,
		Details: details,
	})
}

func InternalServerError(c echo.Context) error {
	return c.JSON(
		http.StatusInternalServerError,
		ErrorResponse{
			Message: "something bad happened",
			Details: nil,
		})
}
