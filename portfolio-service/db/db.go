package db

import (
	"log"
	"portfolio-service/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(path string) {
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	err = database.AutoMigrate(&models.Portfolio{}, &models.Log{})
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}

	DB = database
}
