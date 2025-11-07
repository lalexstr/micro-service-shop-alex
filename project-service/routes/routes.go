package routes

import (
	"ooolalex/project-service/config"
	"ooolalex/project-service/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg config.Config) {
	handlers.RegisterProjectRoutes(r)
}
