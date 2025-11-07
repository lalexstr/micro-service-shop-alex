package handlers

import (
	"net/http"
	"strconv"
	"time"

	"ooolalex/contact-service/config"
	"ooolalex/contact-service/db"
	"ooolalex/contact-service/models"

	"github.com/gin-gonic/gin"
)

type createContactRequest struct {
	Contact string `json:"contact" binding:"required"`
}

type updateStatusRequest struct {
	Status  models.ContactStatus `json:"status" binding:"required"`
	AdminID *uint                `json:"admin_id"`
}

func RegisterContactRoutes(r *gin.Engine, cfg *config.Config) {
	r.POST("api/contact-requests", func(c *gin.Context) {
		var req createContactRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Contact == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		cr := models.ContactRequest{
			Contact: req.Contact,
			Status:  models.ContactNew,
		}

		if err := db.DB.Create(&cr).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
			return
		}

		// ✅ Записываем лог при создании
		db.DB.Create(&models.Log{
			Action:    "create_contact_request",
			Message:   "New contact request: " + req.Contact,
			Timestamp: time.Now(),
		})

		c.JSON(http.StatusOK, gin.H{"message": "Ваш запрос успешно отправлен! Мы свяжемся с вами в ближайшее время."})
	})

	admin := r.Group("/api/contact-requests/admin")

	admin.GET("", func(c *gin.Context) {
		var items []models.ContactRequest
		q := db.DB.Model(&models.ContactRequest{})
		if s := c.Query("status"); s != "" {
			q = q.Where("status = ?", s)
		}
		if err := q.Order("created_at desc").Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, items)
	})

	admin.PATCH(":id/status", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var req updateStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		var cr models.ContactRequest
		if err := db.DB.First(&cr, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		cr.Status = req.Status
		cr.AdminID = req.AdminID
		if err := db.DB.Save(&cr).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}

		// ✅ Логируем изменение статуса
		db.DB.Create(&models.Log{
			UserID:    req.AdminID,
			Action:    "update_contact_status",
			Message:   "Updated contact_request id=" + strconv.Itoa(id) + " to status " + string(req.Status),
			Timestamp: time.Now(),
		})

		c.JSON(http.StatusOK, cr)
	})

	admin.DELETE(":id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := db.DB.Delete(&models.ContactRequest{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}

		// ✅ Логируем удаление
		db.DB.Create(&models.Log{
			Action:    "delete_contact_request",
			Message:   "Deleted contact_request id=" + strconv.Itoa(id),
			Timestamp: time.Now(),
		})

		c.Status(http.StatusNoContent)
	})
}
