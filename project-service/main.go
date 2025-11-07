package main

import (
	"log"
	"ooolalex/project-service/config"
	"ooolalex/project-service/db"
	"ooolalex/project-service/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	if cfg.DBPath == "" {
		cfg.DBPath = "./project.db"
	}

	db.InitDB(cfg.DBPath)
	r := gin.Default()
	routes.SetupRoutes(r, cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Project Service started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
