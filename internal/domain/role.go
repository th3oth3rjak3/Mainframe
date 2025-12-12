package domain

import "github.com/google/uuid"

const (
	Administrator string = "Administrator" // Someone who adminsters the application
	BasicUser     string = "Basic User"    // All users of the application
	RecipeUser    string = "Recipe User"   // Users who can access recipes
)

type Role struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
