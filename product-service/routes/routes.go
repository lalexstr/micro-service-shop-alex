package routes

import (
	"ooolalex/product-service/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes регистрирует все маршруты приложения
func SetupRoutes(r *gin.Engine) {
	handlers.RegisterProductRoutes(r)
	// Позже можно добавить: handlers.RegisterUserRoutes(r), handlers.RegisterOrderRoutes(r) и т.д.
}
