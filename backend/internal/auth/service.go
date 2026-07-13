package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/harshbarnawa/mintok/backend/internal/domain"
	"github.com/harshbarnawa/mintok/backend/internal/repository"
)

// ErrInvalidCredentials is returned when login fails.
var ErrInvalidCredentials = errors.New("invalid email or password")

// ErrEmailTaken is returned when registering an already-used email.
var ErrEmailTaken = errors.New("email already registered")

// ErrValidation is returned when registration or login input is invalid.
var ErrValidation = errors.New("invalid input")

const minPasswordLength = 8

// UserRepository is the persistence dependency required by the auth service.
type UserRepository interface {
	Create(ctx context.Context, params domain.CreateUserParams) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

// Service implements registration and login.
type Service struct {
	users  UserRepository
	hasher PasswordHasher
}

// NewService builds an auth Service.
func NewService(users UserRepository, hasher PasswordHasher) *Service {
	return &Service{users: users, hasher: hasher}
}

// RegisterParams holds the input for account registration.
type RegisterParams struct {
	Email    string
	Name     string
	Password string
}

// LoginParams holds the input for authentication.
type LoginParams struct {
	Email    string
	Password string
}

// Register validates input, hashes the password and stores the new user.
func (s *Service) Register(ctx context.Context, params RegisterParams) (domain.User, error) {
	email := normalizeEmail(params.Email)
	name := strings.TrimSpace(params.Name)

	if email == "" || !strings.Contains(email, "@") {
		return domain.User{}, fmt.Errorf("%w: email is required", ErrValidation)
	}
	if name == "" {
		return domain.User{}, fmt.Errorf("%w: name is required", ErrValidation)
	}
	if len(params.Password) < minPasswordLength {
		return domain.User{}, fmt.Errorf("%w: password must be at least %d characters", ErrValidation, minPasswordLength)
	}

	hash, err := s.hasher.Hash(params.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.users.Create(ctx, domain.CreateUserParams{
		Email:        email,
		Name:         name,
		PasswordHash: hash,
	})
	if err != nil {
		if errors.Is(err, repository.ErrEmailTaken) {
			return domain.User{}, ErrEmailTaken
		}
		return domain.User{}, err
	}

	return user, nil
}

// Login verifies the credentials and returns the matching user.
func (s *Service) Login(ctx context.Context, params LoginParams) (domain.User, error) {
	email := normalizeEmail(params.Email)
	if email == "" || params.Password == "" {
		return domain.User{}, ErrInvalidCredentials
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.User{}, ErrInvalidCredentials
		}
		return domain.User{}, err
	}

	if err := s.hasher.Compare(user.PasswordHash, params.Password); err != nil {
		return domain.User{}, ErrInvalidCredentials
	}

	return user, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
