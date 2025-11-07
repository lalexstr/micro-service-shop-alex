package routes

import (
	"ooolalex/contact-service/config"
	"ooolalex/contact-service/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	handlers.RegisterContactRoutes(r, cfg)
	handlers.RegisterLogRoutes(r)
}
