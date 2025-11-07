package handlers

import (
	"net/http"
	"ooolalex/contact-service/db"
	"ooolalex/contact-service/models"

	"github.com/gin-gonic/gin"
)

func RegisterLogRoutes(r *gin.Engine) {
	r.GET("/api/logs", func(c *gin.Context) {
		var logs []models.Log
		if err := db.DB.Order("timestamp desc").Find(&logs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch logs"})
			return
		}
		c.JSON(http.StatusOK, logs)
	})
}
