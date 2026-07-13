package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment        string
	Port               string
	ServiceName        string
	LogLevel           string
	DatabaseURL        string
	DatabaseMaxConns   int32
	RedisURL           string
	CORSAllowedOrigins []string
	JWTSecret          string
	JWTAccessTTL       time.Duration
	JWTRefreshTTL      time.Duration
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
		LogLevel:           getEnv(lookup, "LOG_LEVEL", "info"),
		DatabaseURL:        getEnv(lookup, "DATABASE_URL", "postgres://mintok:mintok@localhost:5432/mintok?sslmode=disable"),
		DatabaseMaxConns:   int32(getIntEnv(lookup, "DATABASE_MAX_CONNS", 10)),
		RedisURL:           getEnv(lookup, "REDIS_URL", "redis://localhost:6379/0"),
		CORSAllowedOrigins: getListEnv(lookup, "CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		JWTSecret:          getEnv(lookup, "JWT_SECRET", "mintok-development-secret"),
		JWTAccessTTL:       time.Duration(getIntEnv(lookup, "JWT_ACCESS_TTL_MINUTES", 15)) * time.Minute,
		JWTRefreshTTL:      time.Duration(getIntEnv(lookup, "JWT_REFRESH_TTL_HOURS", 720)) * time.Hour,
	}
}

func getEnv(lookup LookupFunc, key, fallback string) string {
	value := strings.TrimSpace(lookup(key))
	if value == "" {
		return fallback
	}

	return value
}

func getIntEnv(lookup LookupFunc, key string, fallback int) int {
	value := strings.TrimSpace(lookup(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
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
