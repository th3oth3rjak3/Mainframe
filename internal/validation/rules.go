// in internal/validation/rules.go

package validation

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// StrongPassword returns a validation rule that checks if a password meets the strength requirements.
func StrongPassword() validation.Rule {
	return validation.By(func(value any) error {
		// First, ensure the value is a string.
		s, ok := value.(string)
		if !ok {
			// This should not happen with a string field, but it's good practice.
			return errors.New("must be a string")
		}

		var (
			hasNumber  int
			hasUpper   int
			hasLower   int
			hasSpecial int
		)

		// Iterate over each character in the string.
		for _, char := range s {
			switch {
			case unicode.IsNumber(char):
				hasNumber++
			case unicode.IsUpper(char):
				hasUpper++
			case unicode.IsLower(char):
				hasLower++
			// Check for common, safe special characters.
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial++
			}
		}

		// Build a list of failures.
		var failures []string
		if hasNumber < 2 {
			failures = append(failures, "at least 2 numbers")
		}
		if hasUpper < 2 {
			failures = append(failures, "at least 2 uppercase letters")
		}
		if hasLower < 2 {
			failures = append(failures, "at least 2 lowercase letters")
		}
		if hasSpecial < 2 {
			failures = append(failures, "at least 2 special characters")
		}

		// If there are any failures, join them into a single, user-friendly error message.
		if len(failures) > 0 {
			return fmt.Errorf("must contain %s", strings.Join(failures, ", "))
		}

		// If we get here, the password is valid.
		return nil
	})
}
