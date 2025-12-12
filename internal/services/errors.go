package services

import "errors"

// These are the "known" errors our service layer can return.
// The handler will check for these specific errors.
var (
	ErrInvalidCredentials = errors.New("invalid credentials provided")
	ErrNotFound           = errors.New("resource not found")
	ErrDuplicateValue     = errors.New("a value with this key already exists")
	ErrValidation         = errors.New("input validation failed")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
)

type ValidationError struct {
	Message string
	Details []string
}

func (v *ValidationError) Error() string {
	return v.Message
}

func NewValidationError(message string, details []string) error {
	return &ValidationError{
		Message: message,
		Details: details,
	}
}
