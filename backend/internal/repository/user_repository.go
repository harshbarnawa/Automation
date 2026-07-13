package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/harshbarnawa/mintok/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a requested record does not exist.
var ErrNotFound = errors.New("record not found")

// ErrEmailTaken is returned when a user email already exists.
var ErrEmailTaken = errors.New("email already registered")

const uniqueViolationCode = "23505"

// UserRepository persists users in PostgreSQL.
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository builds a UserRepository backed by the given pool.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create inserts a new user and returns the stored record.
func (r *UserRepository) Create(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	const query = `
		INSERT INTO users (email, name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, name, password_hash, created_at, updated_at
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, params.Email, params.Name, params.PasswordHash).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return domain.User{}, ErrEmailTaken
		}
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}

// GetByEmail returns the user with the given email, or ErrNotFound.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	const query = `
		SELECT id, email, name, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ErrNotFound
		}
		return domain.User{}, fmt.Errorf("select user by email: %w", err)
	}

	return user, nil
}
