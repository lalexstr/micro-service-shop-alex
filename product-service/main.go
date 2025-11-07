package main

import (
	"log"
	"ooolalex/product-service/config"
	"ooolalex/product-service/db"
	"ooolalex/product-service/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	if cfg.DBPath == "" {
		cfg.DBPath = "./product-service.db"
	}

	db.InitDB(cfg.DBPath)

	r := gin.Default()

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Product Service started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
