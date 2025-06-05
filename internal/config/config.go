package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Address string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	Auth struct {
		SecretKey   string
		TokenExpiry time.Duration
	}
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	var cfg Config

	// Server config
	cfg.Server.Address = getEnv("SERVER_ADDRESS", ":8080")

	// Database config
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.Database.Name = getEnv("DB_NAME", "inventory")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// Auth config
	cfg.Auth.SecretKey = getEnv("AUTH_SECRET_KEY", "secret-key")
	tokenExpiry, err := time.ParseDuration(getEnv("AUTH_TOKEN_EXPIRY", "24h"))
	if err != nil {
		return nil, err
	}
	cfg.Auth.TokenExpiry = tokenExpiry

	return &cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
