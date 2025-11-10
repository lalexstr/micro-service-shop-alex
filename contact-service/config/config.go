package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBPath    string
	LogLevel  string
	JWTSecret string
}

func LoadConfig() Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./contact.db"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// JWT_SECRET is required for admin endpoints
		// But we allow it to be empty for backward compatibility
		// Admin endpoints will fail if JWT_SECRET is not set
	}

	return Config{
		Port:      port,
		DBPath:    dbPath,
		LogLevel:  logLevel,
		JWTSecret: jwtSecret,
	}
}
