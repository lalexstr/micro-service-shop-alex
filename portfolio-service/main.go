package main

import (
	"log"
	"os"
	"portfolio-service/config"
	"portfolio-service/db"
	"portfolio-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	
	if cfg.DBPath == "" {
		cfg.DBPath = "./portfolio.db"
	}
	
	db.InitDB(cfg.DBPath)

	r := gin.Default()
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Portfolio Service running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
