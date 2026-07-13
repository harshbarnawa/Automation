package apikey

import (
	"context"
	"errors"
	"testing"

	"github.com/harshbarnawa/mintok/backend/internal/domain"
	"github.com/harshbarnawa/mintok/backend/internal/repository"
)

type fakeStore struct {
	created domain.APIKey
	hash    string
	err     error
}

func (f *fakeStore) Create(_ context.Context, _, projectID, name, keyHash, keyPrefix string) (domain.APIKey, error) {
	if f.err != nil {
		return domain.APIKey{}, f.err
	}
	f.hash = keyHash
	f.created = domain.APIKey{ID: "key-1", ProjectID: projectID, Name: name, KeyPrefix: keyPrefix}
	return f.created, nil
}
func (f *fakeStore) ListByProject(context.Context, string, string) ([]domain.APIKey, error) {
	return []domain.APIKey{f.created}, f.err
}
func (f *fakeStore) Revoke(context.Context, string, string, string) error { return f.err }

func TestCreateIssuesOpaqueKeyAndStoresOnlyHash(t *testing.T) {
	store := &fakeStore{}
	service := NewService(store)
	service.rand = func(bytes []byte) (int, error) {
		for i := range bytes {
			bytes[i] = byte(i)
		}
		return len(bytes), nil
	}

	metadata, plaintext, err := service.Create(context.Background(), "user-1", "project-1", " Production ")
	if err != nil {
		t.Fatalf("create api key: %v", err)
	}
	if plaintext == "" || plaintext[:7] != "mintok_" {
		t.Fatalf("expected Mintok key, got %q", plaintext)
	}
	if metadata.Name != "Production" {
		t.Fatalf("expected trimmed name, got %q", metadata.Name)
	}
	if store.hash == plaintext || store.hash != hash(plaintext) {
		t.Fatal("expected only a hash to be persisted")
	}
	if metadata.KeyPrefix != prefix(plaintext) {
		t.Fatal("expected key prefix metadata")
	}
}

func TestCreateValidatesInputAndReturnsMissingProject(t *testing.T) {
	service := NewService(&fakeStore{})
	if _, _, err := service.Create(context.Background(), "user", "", "key"); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}

	service = NewService(&fakeStore{err: repository.ErrNotFound})
	if _, _, err := service.Create(context.Background(), "user", "missing", "key"); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}
