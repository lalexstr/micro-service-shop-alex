package routes

import (
	"ooolalex/contact-service/config"
	"ooolalex/contact-service/handlers"
	"ooolalex/contact-service/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	authMiddleware := middleware.AuthMiddleware(cfg)
	adminMiddleware := middleware.AdminMiddleware()
	
	handlers.RegisterContactRoutes(r, cfg, authMiddleware, adminMiddleware)
	handlers.RegisterLogRoutes(r)
}
