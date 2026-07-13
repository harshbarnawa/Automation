package config

import "os"

type Config struct {
	Environment string
	Port        string
	ServiceName string
}

func Load() Config {
	return Config{
		Environment: getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		ServiceName: "mintok-api",
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
