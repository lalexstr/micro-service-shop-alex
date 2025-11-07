package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBPath    string
	JWTSecret string
	JWTTTLMin int
}

func LoadConfig() Config {
	// Загружаем .env файл
	execDir, err := os.Executable()
	if err == nil {
		envPath := filepath.Join(filepath.Dir(execDir), ".env")
		_ = godotenv.Load(envPath)
	}
	_ = godotenv.Load(".env")

	ttl := 60
	if val := os.Getenv("JWT_TTL_MIN"); val != "" {
		ttl, _ = strconv.Atoi(val)
	}

	cfg := Config{
		Port:      os.Getenv("PORT"),
		DBPath:    os.Getenv("DB_PATH"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTTTLMin: ttl,
	}

	if cfg.DBPath == "" {
		cfg.DBPath = "./user.db"
	}
	if cfg.Port == "" {
		cfg.Port = "8085"
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	return cfg
}
