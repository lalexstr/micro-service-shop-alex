package main

import (
	"log"
	"os"
	"portfolio-service/db"
	"portfolio-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB("./portfolio.db")

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
