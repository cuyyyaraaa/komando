package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv        string
	AppPort       string
	DBDSN         string
	JWTSecret     string
	JWTExpiresMin int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	expiresMin := 120
	if v := os.Getenv("JWT_EXPIRES_MIN"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			expiresMin = n
		}
	}

	return &Config{
		AppEnv:        getenv("APP_ENV", "local"),
		AppPort:       getenv("APP_PORT", "8080"),
		DBDSN:         os.Getenv("DB_DSN"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiresMin: expiresMin,
	}, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
