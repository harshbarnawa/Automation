package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

// ErrInvalidRefreshToken is returned when a refresh token is expired, unknown,
// or has already been consumed.
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

// RefreshTokenStore persists only hashes of opaque refresh tokens.
type RefreshTokenStore interface {
	Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	Consume(ctx context.Context, tokenHash string) (string, error)
}

// RefreshTokenManager creates opaque refresh tokens and stores only their hashes.
type RefreshTokenManager struct {
	store RefreshTokenStore
	ttl   time.Duration
	now   Clock
}

// NewRefreshTokenManager builds a manager. A non-positive TTL defaults to 30 days.
func NewRefreshTokenManager(store RefreshTokenStore, ttl time.Duration) *RefreshTokenManager {
	if ttl <= 0 {
		ttl = 30 * 24 * time.Hour
	}
	return &RefreshTokenManager{store: store, ttl: ttl, now: time.Now}
}

// Issue creates and stores a refresh token for a user.
func (m *RefreshTokenManager) Issue(ctx context.Context, userID string) (string, error) {
	token, err := newRefreshToken()
	if err != nil {
		return "", err
	}
	if err := m.store.Create(ctx, userID, hashRefreshToken(token), m.now().Add(m.ttl)); err != nil {
		return "", fmt.Errorf("store refresh token: %w", err)
	}
	return token, nil
}

// Rotate consumes a refresh token and creates a replacement for the same user.
func (m *RefreshTokenManager) Rotate(ctx context.Context, token string) (string, string, error) {
	if token == "" {
		return "", "", ErrInvalidRefreshToken
	}

	userID, err := m.store.Consume(ctx, hashRefreshToken(token))
	if err != nil {
		return "", "", ErrInvalidRefreshToken
	}

	replacement, err := m.Issue(ctx, userID)
	if err != nil {
		return "", "", err
	}
	return userID, replacement, nil
}

func newRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func hashRefreshToken(token string) string {
	digest := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(digest[:])
}
