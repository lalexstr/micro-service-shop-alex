package config

import "os"

type Config struct {
	DBPath string
}

func LoadConfig() Config {
	return Config{
		DBPath: os.Getenv("DB_PATH"),
	}
}
