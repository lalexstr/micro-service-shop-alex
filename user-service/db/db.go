package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(path string, models ...interface{}) {
	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Миграции
	for _, m := range models {
		if err := DB.AutoMigrate(m); err != nil {
			log.Fatal("failed to migrate model:", err)
		}
	}
}
