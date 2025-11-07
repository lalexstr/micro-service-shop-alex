package handlers

import (
	"net/http"
	"strconv"
	"time"

	"portfolio-service/db"
	"portfolio-service/models"

	"github.com/gin-gonic/gin"
)

type CreatePortfolioRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type UpdatePortfolioRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url"`
}

func RegisterPortfolioRoutes(r *gin.Engine) {
	api := r.Group("/api/portfolio")

	// Create
	api.POST("", func(c *gin.Context) {
		var req CreatePortfolioRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		p := models.Portfolio{
			Title:       req.Title,
			Description: req.Description,
			ImageURL:    req.ImageURL,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := db.DB.Create(&p).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create"})
			return
		}

		db.DB.Create(&models.Log{
			Action:    "Created portfolio ID=" + strconv.Itoa(int(p.ID)),
			CreatedAt: time.Now(),
		})

		c.JSON(http.StatusCreated, p)
	})

	// List
	api.GET("", func(c *gin.Context) {
		var items []models.Portfolio
		if err := db.DB.Order("created_at desc").Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list"})
			return
		}
		c.JSON(http.StatusOK, items)
	})

	// Update
	api.PATCH("/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var req UpdatePortfolioRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		var p models.Portfolio
		if err := db.DB.First(&p, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		if req.Title != nil {
			p.Title = *req.Title
		}
		if req.Description != nil {
			p.Description = *req.Description
		}
		if req.ImageURL != nil {
			p.ImageURL = *req.ImageURL
		}
		p.UpdatedAt = time.Now()

		if err := db.DB.Save(&p).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
			return
		}

		db.DB.Create(&models.Log{
			Action:    "Updated portfolio ID=" + strconv.Itoa(int(p.ID)),
			CreatedAt: time.Now(),
		})

		c.JSON(http.StatusOK, p)
	})

	// Delete
	api.DELETE("/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		if err := db.DB.Delete(&models.Portfolio{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
			return
		}

		db.DB.Create(&models.Log{
			Action:    "Deleted portfolio ID=" + strconv.Itoa(id),
			CreatedAt: time.Now(),
		})

		c.Status(http.StatusNoContent)
	})
}
