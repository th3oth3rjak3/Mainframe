package domain

import (
	"slices"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	v "github.com/th3oth3rjak3/mainframe/internal/validation"
)

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

type UserRead struct {
	ID                     uuid.UUID  `json:"id"`
	Username               string     `json:"username"`
	Email                  string     `json:"email"`
	FirstName              string     `json:"firstName"`
	LastName               string     `json:"lastName"`
	LastLogin              *time.Time `json:"lastLogin"`
	FailedLoginAttempts    uint       `json:"failedLoginAttempts"`
	LastFailedLoginAttempt *time.Time `json:"lastFailedLoginAttempt"`
	IsDisabled             bool       `json:"isDisabled"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
	Roles                  []Role     `json:"roles"`
}

func NewUserRead(user *User) UserRead {
	return UserRead{
		ID:                     user.ID,
		Username:               user.Username,
		Email:                  user.Email,
		FirstName:              user.FirstName,
		LastName:               user.LastName,
		LastLogin:              user.LastLogin,
		FailedLoginAttempts:    user.FailedLoginAttempts,
		LastFailedLoginAttempt: user.LastFailedLoginAttempt,
		IsDisabled:             user.IsDisabled,
		CreatedAt:              user.CreatedAt,
		UpdatedAt:              user.UpdatedAt,
		Roles:                  user.Roles,
	}
}

type UserCreate struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (u *UserCreate) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.FirstName, validation.Required, validation.Length(1, 100), is.Alpha),
		validation.Field(&u.LastName, validation.Required, validation.Length(1, 100), is.Alpha),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Username, validation.Required, validation.Length(3, 50)),
		validation.Field(&u.Password, v.StrongPassword()),
	)
}

type UserUpdate struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}
