package shared

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
)

// These are the "known" errors our service layer can return.
// The handler will check for these specific errors.
var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrNotFound           = errors.New("resource not found")
	ErrDuplicateValue     = errors.New("a value with this key already exists")
	ErrBadRequest         = errors.New("bad request")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrUsernameTaken      = errors.New("username already exists")
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func JsonError(c *fiber.Ctx, err error, httpStatus int) error {
	return c.JSON(ErrorResponse{
		Message: err.Error(),
	})
}

func InternalServerError(c *fiber.Ctx) error {
	return c.JSON(ErrorResponse{
		Message: "something bad happened",
	})
}

// in your api/server.go

// ResolveError inspects an error and returns the appropriate HTTP status code
// and a user-facing error message. This is our single source of truth.
func ResolveError(err error) (int, string) {
	// 1. Handle echo.HTTPError (from binding, etc.)
	var he *echo.HTTPError
	if errors.As(err, &he) {
		return he.Code, he.Message.(string)
	}

	// 2. Handle validation errors
	if _, ok := err.(validation.Errors); ok {
		return http.StatusBadRequest, "The request contains invalid data."
	}

	// 3. Handle our specific business/sentinel errors
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, ErrUsernameTaken):
		return http.StatusConflict, err.Error()
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, ErrUnauthorized), errors.Is(err, ErrInvalidCredentials):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest, err.Error()
	default:
		// Any error that isn't specifically handled is a 500
		return http.StatusInternalServerError, "An internal server error occurred."
	}
}
