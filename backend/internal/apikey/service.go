package apikey

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/harshbarnawa/mintok/backend/internal/domain"
	"github.com/harshbarnawa/mintok/backend/internal/repository"
)

var ErrValidation = errors.New("invalid input")

// Store persists API-key metadata and validates project ownership.
type Store interface {
	Create(ctx context.Context, userID, projectID, name, keyHash, keyPrefix string) (domain.APIKey, error)
	ListByProject(ctx context.Context, userID, projectID string) ([]domain.APIKey, error)
	Revoke(ctx context.Context, userID, projectID, keyID string) error
}

// Service creates and manages project-scoped gateway credentials.
type Service struct {
	store Store
	rand  func([]byte) (int, error)
}

func NewService(store Store) *Service {
	return &Service{store: store, rand: rand.Read}
}

func (s *Service) Create(ctx context.Context, userID, projectID, name string) (domain.APIKey, string, error) {
	name = strings.TrimSpace(name)
	if userID == "" || projectID == "" || name == "" {
		return domain.APIKey{}, "", fmt.Errorf("%w: project and name are required", ErrValidation)
	}

	plaintext, err := generateKey(s.rand)
	if err != nil {
		return domain.APIKey{}, "", err
	}
	key, err := s.store.Create(ctx, userID, projectID, name, hash(plaintext), prefix(plaintext))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.APIKey{}, "", repository.ErrNotFound
		}
		return domain.APIKey{}, "", err
	}
	return key, plaintext, nil
}

func (s *Service) List(ctx context.Context, userID, projectID string) ([]domain.APIKey, error) {
	if userID == "" || projectID == "" {
		return nil, fmt.Errorf("%w: project is required", ErrValidation)
	}
	return s.store.ListByProject(ctx, userID, projectID)
}

func (s *Service) Revoke(ctx context.Context, userID, projectID, keyID string) error {
	if userID == "" || projectID == "" || keyID == "" {
		return fmt.Errorf("%w: project and key are required", ErrValidation)
	}
	return s.store.Revoke(ctx, userID, projectID, keyID)
}

func generateKey(read func([]byte) (int, error)) (string, error) {
	bytes := make([]byte, 32)
	if _, err := read(bytes); err != nil {
		return "", fmt.Errorf("generate api key: %w", err)
	}
	return "mintok_" + base64.RawURLEncoding.EncodeToString(bytes), nil
}

func hash(value string) string {
	digest := sha256.Sum256([]byte(value))
	return base64.RawURLEncoding.EncodeToString(digest[:])
}

func prefix(value string) string {
	const visibleCharacters = 14
	if len(value) <= visibleCharacters {
		return value
	}
	return value[:visibleCharacters]
}
