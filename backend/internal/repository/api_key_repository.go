package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/harshbarnawa/mintok/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// APIKeyRepository persists project-scoped gateway API-key metadata.
type APIKeyRepository struct {
	pool *pgxpool.Pool
}

// NewAPIKeyRepository builds an API-key repository backed by PostgreSQL.
func NewAPIKeyRepository(pool *pgxpool.Pool) *APIKeyRepository {
	return &APIKeyRepository{pool: pool}
}

// Create stores a hashed API key and returns its public metadata.
func (r *APIKeyRepository) Create(ctx context.Context, userID, projectID, name, keyHash, keyPrefix string) (domain.APIKey, error) {
	const query = `
		INSERT INTO api_keys (user_id, project_id, name, key_hash, key_prefix)
		SELECT $1, p.id, $3, $4, $5
		FROM projects p
		WHERE p.id = $2 AND p.owner_id = $1
		RETURNING id, project_id, name, key_prefix, last_used_at, created_at, revoked_at
	`

	key, err := scanAPIKey(r.pool.QueryRow(ctx, query, userID, projectID, name, keyHash, keyPrefix))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.APIKey{}, ErrNotFound
	}
	if err != nil {
		return domain.APIKey{}, fmt.Errorf("insert api key: %w", err)
	}
	return key, nil
}

// ListByProject returns active and revoked API-key metadata for a project the user owns.
func (r *APIKeyRepository) ListByProject(ctx context.Context, userID, projectID string) ([]domain.APIKey, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT k.id, k.project_id, k.name, k.key_prefix, k.last_used_at, k.created_at, k.revoked_at
		FROM api_keys k
		JOIN projects p ON p.id = k.project_id
		WHERE k.project_id = $1 AND p.owner_id = $2
		ORDER BY k.created_at DESC
	`, projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}
	defer rows.Close()

	keys := make([]domain.APIKey, 0)
	for rows.Next() {
		key, err := scanAPIKey(rows)
		if err != nil {
			return nil, fmt.Errorf("scan api key: %w", err)
		}
		keys = append(keys, key)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate api keys: %w", err)
	}
	return keys, nil
}

// Revoke marks an API key inactive if it belongs to a project the user owns.
func (r *APIKeyRepository) Revoke(ctx context.Context, userID, projectID, keyID string) error {
	command, err := r.pool.Exec(ctx, `
		UPDATE api_keys k
		SET revoked_at = now()
		FROM projects p
		WHERE k.id = $1 AND k.project_id = $2 AND k.project_id = p.id
			AND p.owner_id = $3 AND k.revoked_at IS NULL
	`, keyID, projectID, userID)
	if err != nil {
		return fmt.Errorf("revoke api key: %w", err)
	}
	if command.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Authenticate returns the project for an active API-key hash and records its use.
func (r *APIKeyRepository) Authenticate(ctx context.Context, keyHash string) (string, error) {
	var projectID string
	err := r.pool.QueryRow(ctx, `
		UPDATE api_keys SET last_used_at = now()
		WHERE key_hash = $1 AND revoked_at IS NULL AND project_id IS NOT NULL
		RETURNING project_id
	`, keyHash).Scan(&projectID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("authenticate api key: %w", err)
	}
	return projectID, nil
}

type apiKeyRow interface {
	Scan(...any) error
}

func scanAPIKey(row apiKeyRow) (domain.APIKey, error) {
	var key domain.APIKey
	err := row.Scan(&key.ID, &key.ProjectID, &key.Name, &key.KeyPrefix, &key.LastUsedAt, &key.CreatedAt, &key.RevokedAt)
	return key, err
}
