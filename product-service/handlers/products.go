package handlers

import (
	"math"
	"net/http"
	"ooolalex/product-service/db"
	"ooolalex/product-service/logs"
	"ooolalex/product-service/middleware"
	"ooolalex/product-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createProductRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	ImageURL    string  `json:"image_url"`
}

type updateProductRequest struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	ImageURL    *string  `json:"image_url"`
}

func RegisterProductRoutes(r *gin.Engine) {
	admin := r.Group("/api/products")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware()) // JWT проверка через Auth Service

	admin.POST("", CreateProduct)
	admin.GET("", ListProducts)
	admin.PATCH(":id", UpdateProduct)
	admin.DELETE(":id", DeleteProduct)

	r.GET("/api/products/public", PublicListProducts)
}

// CreateProduct создаёт новый продукт
func CreateProduct(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	p := models.Product{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    req.ImageURL,
	}
	if err := db.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create"})
		return
	}

	userID := c.GetUint("userID")
	go logs.SendLog(userID, "created product-service id="+strconv.Itoa(int(p.ID)))

	c.JSON(http.StatusCreated, p)
}

// ListProducts возвращает список продуктов для админов с пагинацией
func ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	var total int64
	db.DB.Model(&models.Product{}).Count(&total)

	var items []models.Product
	db.DB.Order("created_at desc").Offset((page - 1) * size).Limit(size).Find(&items)

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"page":  page,
		"size":  size,
		"total": total,
		"pages": int(math.Ceil(float64(total) / float64(size))),
	})
}

// UpdateProduct обновляет существующий продукт
func UpdateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req updateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var p models.Product
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
	if req.Price != nil {
		p.Price = *req.Price
	}
	if req.ImageURL != nil {
		p.ImageURL = *req.ImageURL
	}

	if err := db.DB.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	userID := c.GetUint("userID")
	go logs.SendLog(userID, "updated product-service id="+strconv.Itoa(int(p.ID)))

	c.JSON(http.StatusOK, p)
}

// DeleteProduct удаляет продукт
func DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Сначала получаем продукт, чтобы иметь p.ID
	var p models.Product
	if err := db.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	if err := db.DB.Delete(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	userID := c.GetUint("userID")
	go logs.SendLog(userID, "deleted product-service id="+strconv.Itoa(int(p.ID)))

	c.Status(http.StatusNoContent)
}

// PublicListProducts возвращает список продуктов без авторизации
func PublicListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	var total int64
	db.DB.Model(&models.Product{}).Count(&total)

	var items []models.Product
	db.DB.Order("created_at desc").Offset((page - 1) * size).Limit(size).Find(&items)

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"page":  page,
		"size":  size,
		"total": total,
		"pages": int(math.Ceil(float64(total) / float64(size))),
	})
}
