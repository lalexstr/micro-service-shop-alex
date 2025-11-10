package routes

import (
	"portfolio-service/handlers"
	"portfolio-service/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	handlers.RegisterPortfolioRoutes(r, middleware.AuthMiddleware(), middleware.AdminMiddleware())
}
