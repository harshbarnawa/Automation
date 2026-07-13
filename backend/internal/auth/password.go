package auth

import "golang.org/x/crypto/bcrypt"

// PasswordHasher hashes and verifies user passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

// BcryptHasher implements PasswordHasher using bcrypt.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher builds a hasher with the given cost, falling back to the
// bcrypt default cost when a non-positive value is supplied.
func NewBcryptHasher(cost int) BcryptHasher {
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	return BcryptHasher{cost: cost}
}

// Hash returns the bcrypt hash of the given password.
func (h BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Compare reports whether the hash matches the password.
func (h BcryptHasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
