package database

import (
	"testing"

	"github.com/harshbarnawa/mintok/backend/internal/config"
)

func TestNewPoolConfigAppliesDatabaseSettings(t *testing.T) {
	poolConfig, err := NewPoolConfig(config.Config{
		DatabaseURL:      "postgres://mintok:mintok@localhost:5432/mintok?sslmode=disable",
		DatabaseMaxConns: 12,
	})
	if err != nil {
		t.Fatalf("expected valid pool config, got error: %v", err)
	}

	if poolConfig.ConnConfig.Database != "mintok" {
		t.Fatalf("expected mintok database, got %q", poolConfig.ConnConfig.Database)
	}

	if poolConfig.MaxConns != 12 {
		t.Fatalf("expected max conns 12, got %d", poolConfig.MaxConns)
	}
}

func TestNewPoolConfigRejectsInvalidURL(t *testing.T) {
	_, err := NewPoolConfig(config.Config{
		DatabaseURL:      "://bad-url",
		DatabaseMaxConns: 10,
	})
	if err == nil {
		t.Fatal("expected invalid database url error")
	}
}
