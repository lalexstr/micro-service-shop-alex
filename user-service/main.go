package main

import (
	"log"

	"user-service/config"
	"user-service/db"
	"user-service/middleware"
	"user-service/models"
	"user-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Инициализация БД (только Log, User хранится в auth-service)
	db.InitDB(cfg.DBPath, &models.Log{})

	r := gin.Default()

	// Регистрируем маршруты
	mwCfg := &middleware.Config{JWTSecret: cfg.JWTSecret}
	routes.SetupRoutes(r, mwCfg)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Service running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
