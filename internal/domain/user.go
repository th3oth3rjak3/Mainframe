package domain

import (
	"time"

	"github.com/google/uuid"
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
}
