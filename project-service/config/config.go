package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath string
}

func LoadConfig() Config {
	// Загружаем .env файл
	_ = godotenv.Load(".env")

	return Config{
		DBPath: os.Getenv("DB_PATH"),
	}
}
