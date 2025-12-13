package domain

import (
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserCreate struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type UserUpdate struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type User struct {
	ID                     uuid.UUID  `db:"id"`
	Username               string     `db:"username"`
	Email                  string     `db:"email"`
	FirstName              string     `db:"first_name"`
	LastName               string     `db:"last_name"`
	PasswordHash           string     `db:"password_hash"`
	LastLogin              *time.Time `db:"last_login"`
	FailedLoginAttempts    uint       `db:"failed_login_attempts"`
	LastFailedLoginAttempt *time.Time `db:"last_failed_login_attempt"`
	IsDisabled             bool       `db:"is_disabled"`
	CreatedAt              time.Time  `db:"created_at"`
	UpdatedAt              time.Time  `db:"updated_at"`
	Roles                  []Role     `db:"-"`
}

func NewUser(
	username string,
	email string,
	firstName string,
	lastName string,
	passwordHash string,
	roles []Role,
) (*User, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	newUser := &User{
		ID:           id,
		Username:     username,
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Roles:        roles,
	}

	return newUser, nil
}

func (u *User) HasRole(name string) bool {
	return slices.ContainsFunc(u.Roles, func(r Role) bool {
		return strings.EqualFold(r.Name, name)
	})
}
