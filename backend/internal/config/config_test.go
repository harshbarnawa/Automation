package config

import (
	"reflect"
	"testing"
)

func TestLoadWithLookupUsesDefaults(t *testing.T) {
	cfg := LoadWithLookup(func(string) string {
		return ""
	})

	if cfg.Environment != "development" {
		t.Fatalf("expected development environment, got %q", cfg.Environment)
	}

	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.Port)
	}

	if cfg.ServiceName != "mintok-api" {
		t.Fatalf("expected service name mintok-api, got %q", cfg.ServiceName)
	}

	if cfg.LogLevel != "info" {
		t.Fatalf("expected default log level info, got %q", cfg.LogLevel)
	}

	expectedOrigins := []string{"http://localhost:3000"}
	if !reflect.DeepEqual(cfg.CORSAllowedOrigins, expectedOrigins) {
		t.Fatalf("expected default origins %v, got %v", expectedOrigins, cfg.CORSAllowedOrigins)
	}
}

func TestLoadWithLookupUsesEnvironmentValues(t *testing.T) {
	values := map[string]string{
		"APP_ENV":              "production",
		"PORT":                 "9090",
		"SERVICE_NAME":         "mintok-test",
		"LOG_LEVEL":            "debug",
		"DATABASE_URL":         "postgres://user:pass@db:5432/app",
		"REDIS_URL":            "redis://cache:6379/2",
		"CORS_ALLOWED_ORIGINS": "https://app.example.com, https://admin.example.com",
	}

	cfg := LoadWithLookup(func(key string) string {
		return values[key]
	})

	if cfg.Environment != values["APP_ENV"] {
		t.Fatalf("expected configured environment, got %q", cfg.Environment)
	}

	if cfg.Port != values["PORT"] {
		t.Fatalf("expected configured port, got %q", cfg.Port)
	}

	if cfg.DatabaseURL != values["DATABASE_URL"] {
		t.Fatalf("expected configured database url, got %q", cfg.DatabaseURL)
	}

	if cfg.LogLevel != values["LOG_LEVEL"] {
		t.Fatalf("expected configured log level, got %q", cfg.LogLevel)
	}

	expectedOrigins := []string{"https://app.example.com", "https://admin.example.com"}
	if !reflect.DeepEqual(cfg.CORSAllowedOrigins, expectedOrigins) {
		t.Fatalf("expected origins %v, got %v", expectedOrigins, cfg.CORSAllowedOrigins)
	}
}
