package domain

import (
	"github.com/alexedwards/argon2id"
)

// fakeHash is used in conjunction with the FakeVerify method to prevent timing attack.
const fakeHash = "$argon2id$v=19$m=65536,t=1,p=16$5+5ObcY5s1LVbxJ/+Xwajg$EtvmraG0bszkPPJW4k3RFYy6UcXZTQahKIl7TLdJ0TE"

type PasswordHasher interface {
	// HashPassword calculates a new password hash for the input password.
	//
	// When the hash is successful, it returns the new hash, otherwise an error.
	HashPassword(password string) (string, error)

	// Verify compares the provided password to an existing password hash.
	//
	// When hashing is successful, a true value indicates a match, a false value indicates
	// a mismatch.
	Verify(password string, hash string) (bool, error)

	// FakeVerify does a constant-time password verification for the purposes of preventing
	// timing attacks.
	//
	// Since the verification process is only to prevent timing attacks,
	// it only requires the input password and will be compared against
	// a constant hash value.
	FakeVerify(password string) error
}

// NewPasswordHasher creates a new password hasher.
func NewPasswordHasher() PasswordHasher {
	return &argon2PasswordHasher{}
}

type argon2PasswordHasher struct{}

func (h *argon2PasswordHasher) HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func (h *argon2PasswordHasher) Verify(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func (h *argon2PasswordHasher) FakeVerify(password string) error {

	_, err := h.Verify(password, fakeHash)
	return err
}
