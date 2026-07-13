package config

import (
	"os"
	"strings"
)

type Config struct {
	Environment        string
	Port               string
	ServiceName        string
	DatabaseURL        string
	RedisURL           string
	CORSAllowedOrigins []string
}

type LookupFunc func(string) string

func Load() Config {
	return LoadWithLookup(os.Getenv)
}

func LoadWithLookup(lookup LookupFunc) Config {
	return Config{
		Environment:        getEnv(lookup, "APP_ENV", "development"),
		Port:               getEnv(lookup, "PORT", "8080"),
		ServiceName:        getEnv(lookup, "SERVICE_NAME", "mintok-api"),
		DatabaseURL:        getEnv(lookup, "DATABASE_URL", "postgres://mintok:mintok@localhost:5432/mintok?sslmode=disable"),
		RedisURL:           getEnv(lookup, "REDIS_URL", "redis://localhost:6379/0"),
		CORSAllowedOrigins: getListEnv(lookup, "CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
	}
}

func getEnv(lookup LookupFunc, key, fallback string) string {
	value := strings.TrimSpace(lookup(key))
	if value == "" {
		return fallback
	}

	return value
}

func getListEnv(lookup LookupFunc, key string, fallback []string) []string {
	value := strings.TrimSpace(lookup(key))
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return fallback
	}

	return result
}
