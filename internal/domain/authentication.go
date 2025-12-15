package domain

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	v "github.com/th3oth3rjak3/mainframe/internal/validation"
)

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Username, validation.Required, validation.Length(3, 50)),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 100), v.StrongPassword()),
	)
}

// LoginResponse represents the login response
type LoginResponse struct {
	Username  string     `json:"username" example:"admin"`
	Email     string     `json:"email"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	LastLogin *time.Time `json:"lastLogin"`
	Roles     []string   `json:"roles"`
}

func NewLoginResponse(user *User) *LoginResponse {
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	return &LoginResponse{
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		LastLogin: user.LastLogin,
		Roles:     roles,
	}
}
