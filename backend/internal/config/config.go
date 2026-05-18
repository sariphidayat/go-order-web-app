package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv            string
	AppPort           string
	FrontendOrigin    string
	JWTSecret         string
	JWTExpiresMinutes int
	DatabaseDSN       string
	AdminEmail        string
	AdminPassword     string
	AdminTOTPSecret   string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppEnv:            getEnv("APP_ENV", "development"),
		AppPort:           getEnv("APP_PORT", "8080"),
		FrontendOrigin:    getEnv("APP_FRONTEND_ORIGIN", "http://localhost:5173"),
		JWTSecret:         getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiresMinutes: getEnvInt("JWT_EXPIRES_MINUTES", 120),
		DatabaseDSN:       getEnv("DATABASE_DSN", "app.db"),
		AdminEmail:        getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword:     getEnv("ADMIN_PASSWORD", "admin12345"),
		AdminTOTPSecret:   getEnv("ADMIN_TOTP_SECRET", "JBSWY3DPEHPK3PXP"),
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
