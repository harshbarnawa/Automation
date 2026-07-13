package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/harshbarnawa/mintok/backend/internal/domain"
	"github.com/harshbarnawa/mintok/backend/internal/repository"
)

// fakeUserRepo is an in-memory UserRepository for tests.
type fakeUserRepo struct {
	byEmail map[string]domain.User
	seq     int
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{byEmail: map[string]domain.User{}}
}

func (f *fakeUserRepo) Create(_ context.Context, params domain.CreateUserParams) (domain.User, error) {
	if _, ok := f.byEmail[params.Email]; ok {
		return domain.User{}, repository.ErrEmailTaken
	}
	f.seq++
	user := domain.User{
		ID:           string(rune('a' + f.seq)),
		Email:        params.Email,
		Name:         params.Name,
		PasswordHash: params.PasswordHash,
	}
	f.byEmail[params.Email] = user
	return user, nil
}

func (f *fakeUserRepo) GetByEmail(_ context.Context, email string) (domain.User, error) {
	user, ok := f.byEmail[email]
	if !ok {
		return domain.User{}, repository.ErrNotFound
	}
	return user, nil
}

func newTestService() *Service {
	return NewService(newFakeUserRepo(), NewBcryptHasher(4))
}

func TestRegisterSuccess(t *testing.T) {
	svc := newTestService()

	user, err := svc.Register(context.Background(), RegisterParams{
		Email:    "User@Example.com ",
		Name:     " Ada ",
		Password: "supersecret",
	})
	if err != nil {
		t.Fatalf("expected registration to succeed, got %v", err)
	}

	if user.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", user.Email)
	}
	if user.Name != "Ada" {
		t.Fatalf("expected trimmed name, got %q", user.Name)
	}
	if user.PasswordHash == "supersecret" {
		t.Fatal("expected password to be hashed")
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	svc := newTestService()
	params := RegisterParams{Email: "dup@example.com", Name: "Dup", Password: "supersecret"}

	if _, err := svc.Register(context.Background(), params); err != nil {
		t.Fatalf("unexpected error on first register: %v", err)
	}

	_, err := svc.Register(context.Background(), params)
	if !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestRegisterValidation(t *testing.T) {
	svc := newTestService()

	cases := []RegisterParams{
		{Email: "", Name: "A", Password: "supersecret"},
		{Email: "bad", Name: "A", Password: "supersecret"},
		{Email: "a@b.com", Name: "", Password: "supersecret"},
		{Email: "a@b.com", Name: "A", Password: "short"},
	}

	for _, params := range cases {
		if _, err := svc.Register(context.Background(), params); !errors.Is(err, ErrValidation) {
			t.Fatalf("expected ErrValidation for %+v, got %v", params, err)
		}
	}
}

func TestLoginSuccess(t *testing.T) {
	svc := newTestService()
	if _, err := svc.Register(context.Background(), RegisterParams{
		Email:    "login@example.com",
		Name:     "Login",
		Password: "supersecret",
	}); err != nil {
		t.Fatalf("unexpected register error: %v", err)
	}

	user, err := svc.Login(context.Background(), LoginParams{
		Email:    "LOGIN@example.com",
		Password: "supersecret",
	})
	if err != nil {
		t.Fatalf("expected login to succeed, got %v", err)
	}
	if user.Email != "login@example.com" {
		t.Fatalf("unexpected user %q", user.Email)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	svc := newTestService()
	if _, err := svc.Register(context.Background(), RegisterParams{
		Email:    "wrong@example.com",
		Name:     "Wrong",
		Password: "supersecret",
	}); err != nil {
		t.Fatalf("unexpected register error: %v", err)
	}

	_, err := svc.Login(context.Background(), LoginParams{
		Email:    "wrong@example.com",
		Password: "incorrect",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLoginUnknownUser(t *testing.T) {
	svc := newTestService()

	_, err := svc.Login(context.Background(), LoginParams{
		Email:    "missing@example.com",
		Password: "supersecret",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
