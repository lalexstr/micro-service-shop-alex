package routes

import (
	"user-service/config"
	"user-service/handlers"
	"user-service/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *middleware.Config) {
	// Пользовательские маршруты
	// Преобразуем middleware.Config в config.Config для handlers
	configCfg := &config.Config{
		JWTSecret: cfg.JWTSecret,
	}
	handlers.RegisterUserRoutes(r, configCfg)
	// Можно добавить сюда остальные: проекты, контакты, портфолио
}
