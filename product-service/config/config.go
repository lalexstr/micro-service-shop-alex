package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
	DBPath    string
}

func LoadConfig() Config {
	// Загружаем .env файл
	_ = godotenv.Load(".env")

	return Config{
		JWTSecret: os.Getenv("JWT_SECRET"),
		DBPath:    os.Getenv("DB_PATH"), // например "./products.db"
	}
}
