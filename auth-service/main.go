package main

import (
	"log"

	"auth-service/config"
	"auth-service/db"
	"auth-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()

	r := gin.Default()

	routes.Setup(r)

	addr := ":" + config.Port
	log.Printf("Starting auth service on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
