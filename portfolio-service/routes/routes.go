package routes

import (
	"portfolio-service/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	handlers.RegisterPortfolioRoutes(r)
}
