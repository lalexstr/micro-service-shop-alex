package handlers

import (
	"net/http"
	"strconv"
	"time"

	"user-service/clients"
	"user-service/config"
	"user-service/db"
	"user-service/middleware"
	"user-service/models"

	"github.com/gin-gonic/gin"
)

type changeRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

func RegisterUserRoutes(r *gin.Engine, cfg *config.Config) {
	admin := r.Group("/api/users")
	admin.Use(middleware.AuthRequired(middleware.Config{JWTSecret: cfg.JWTSecret}), middleware.AdminOnly(cfg))

	admin.GET("", listUsers)
	admin.PATCH(":id/role", changeUserRole)
	admin.GET(":id/activity", userActivity)
}

func listUsers(c *gin.Context) {
	authClient := clients.NewAuthClient()
	
	page := 1
	size := 20
	if p := c.Query("page"); p != "" {
		if pv, err := strconv.Atoi(p); err == nil && pv > 0 {
			page = pv
		}
	}
	if s := c.Query("size"); s != "" {
		if sv, err := strconv.Atoi(s); err == nil && sv > 0 {
			size = sv
		}
	}

	filters := make(map[string]string)
	if role := c.Query("role"); role != "" {
		filters["role"] = role
	}
	if email := c.Query("email"); email != "" {
		filters["email"] = email
	}

	users, total, err := authClient.ListUsers(page, size, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users from auth-service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": users,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

func changeUserRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req changeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil || (req.Role != "admin" && req.Role != "user") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	authClient := clients.NewAuthClient()
	
	// Get user from auth-service to verify it exists
	user, err := authClient.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found in auth-service"})
		return
	}

	// Update role in auth-service via HTTP call
	if err := authClient.UpdateUserRole(user.ID, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update role in auth-service"})
		return
	}
	
	actorRaw, _ := c.Get(middleware.ContextUserID)
	actor, _ := actorRaw.(uint)
	db.DB.Create(&models.Log{
		UserID:    actor,
		Action:    "changed role for user id=" + strconv.Itoa(int(user.ID)) + " from " + user.Role + " to " + req.Role,
		Timestamp: time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"id": user.ID, "role": req.Role})
}

func userActivity(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var items []models.Log
	q := db.DB.Where("user_id = ?", id)
	if from := c.Query("from"); from != "" {
		q = q.Where("timestamp >= ?", from)
	}
	if to := c.Query("to"); to != "" {
		q = q.Where("timestamp <= ?", to)
	}
	if err := q.Order("timestamp desc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, items)
}
