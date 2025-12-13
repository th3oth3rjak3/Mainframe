package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `db:"id"`
	Token     string    `db:"token"`
	UserID    uuid.UUID `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
}

func NewSession(userID uuid.UUID, token string) (*Session, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:        id,
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 2),
	}

	return session, nil
}
