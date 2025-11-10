package config

import (
	"log"
	"os"
)

type Config struct {
	DBPath    string
	JWTSecret string
}

func LoadConfig() Config {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./portfolio.db"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	return Config{
		DBPath:    dbPath,
		JWTSecret: jwtSecret,
	}
}
