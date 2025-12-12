package domain

import (
	"strings"
	"time"
)

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (lr *LoginRequest) Validate() []string {
	var errs []string

	if strings.TrimSpace(lr.Username) == "" {
		errs = append(errs, "username is required")
	}

	if strings.TrimSpace(lr.Password) == "" {
		errs = append(errs, "password is required")
	}

	return errs
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
