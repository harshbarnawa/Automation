package auth

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/harshbarnawa/mintok/backend/internal/repository"
)

type fakeRefreshTokenStore struct {
	mu     sync.Mutex
	tokens map[string]string
}

func newFakeRefreshTokenStore() *fakeRefreshTokenStore {
	return &fakeRefreshTokenStore{tokens: map[string]string{}}
}

func (f *fakeRefreshTokenStore) Create(_ context.Context, userID, tokenHash string, _ time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.tokens[tokenHash] = userID
	return nil
}

func (f *fakeRefreshTokenStore) Consume(_ context.Context, tokenHash string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	userID, ok := f.tokens[tokenHash]
	if !ok {
		return "", repository.ErrNotFound
	}
	delete(f.tokens, tokenHash)
	return userID, nil
}

func TestRefreshTokenManagerRotatesSingleUseTokens(t *testing.T) {
	manager := NewRefreshTokenManager(newFakeRefreshTokenStore(), time.Hour)
	first, err := manager.Issue(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	userID, second, err := manager.Rotate(context.Background(), first)
	if err != nil {
		t.Fatalf("rotate token: %v", err)
	}
	if userID != "user-1" {
		t.Fatalf("expected user-1, got %q", userID)
	}
	if second == first {
		t.Fatal("expected a new refresh token")
	}

	if _, _, err := manager.Rotate(context.Background(), first); err != ErrInvalidRefreshToken {
		t.Fatalf("expected invalid reused token, got %v", err)
	}
}
