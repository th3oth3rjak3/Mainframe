package domain

import "github.com/google/uuid"

const (
	Administrator string = "Administrator" // Someone who adminsters the application
	BasicUser     string = "Basic User"    // All users of the application are considered basic users
	RecipeUser    string = "Recipe User"   // Users who can access the recipe features
)

type Role struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}
