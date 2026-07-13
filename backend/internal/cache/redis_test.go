package cache

import (
	"testing"

	"github.com/harshbarnawa/mintok/backend/internal/config"
)

func TestNewRedisOptionsParsesURL(t *testing.T) {
	options, err := NewRedisOptions(config.Config{
		RedisURL: "redis://:secret@localhost:6380/3",
	})
	if err != nil {
		t.Fatalf("expected redis options, got error: %v", err)
	}

	if options.Addr != "localhost:6380" {
		t.Fatalf("expected localhost:6380, got %q", options.Addr)
	}

	if options.Password != "secret" {
		t.Fatal("expected redis password to be parsed")
	}

	if options.DB != 3 {
		t.Fatalf("expected redis db 3, got %d", options.DB)
	}
}

func TestNewRedisOptionsRejectsInvalidURL(t *testing.T) {
	_, err := NewRedisOptions(config.Config{
		RedisURL: "://bad-url",
	})
	if err == nil {
		t.Fatal("expected invalid redis url error")
	}
}
