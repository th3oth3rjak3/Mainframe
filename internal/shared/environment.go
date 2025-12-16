package shared

import (
	"os"
	"strings"
)

func IsProduction() bool {
	env := os.Getenv("APP_ENV")
	return strings.EqualFold(env, "production")
}
