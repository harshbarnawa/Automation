package domain

import "time"

// APIKey is the persisted metadata for a project gateway credential. The
// plaintext credential is intentionally never stored.
type APIKey struct {
	ID         string
	ProjectID  string
	Name       string
	KeyPrefix  string
	LastUsedAt *time.Time
	CreatedAt  time.Time
	RevokedAt  *time.Time
}
