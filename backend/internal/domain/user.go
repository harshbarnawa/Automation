package domain

import "time"

// User is the source-of-truth representation of an authenticated account.
type User struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CreateUserParams carries the fields required to persist a new user.
type CreateUserParams struct {
	Email        string
	Name         string
	PasswordHash string
}

// UpdateUserParams carries the mutable profile fields for an existing user.
type UpdateUserParams struct {
	Name string
}
