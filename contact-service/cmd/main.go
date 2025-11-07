package main

import (
	"log"

	"ooolalex/contact-service/config"
	"ooolalex/contact-service/db"
	"ooolalex/contact-service/models"
	"ooolalex/contact-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db.InitDB(cfg.DBPath, &models.ContactRequest{}, &models.Log{})

	r := gin.Default()
	routes.SetupRoutes(r, &cfg)

	port := cfg.Port
	if port == "" {
		port = "8084"
	}

	log.Printf("Contact Service running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
