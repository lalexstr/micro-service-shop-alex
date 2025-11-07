package handlers

import (
	"net/http"
	"ooolalex/project-service/db"
	"ooolalex/project-service/logs"
	"ooolalex/project-service/middleware"
	"ooolalex/project-service/models"
	"time"

	"github.com/gin-gonic/gin"
)

type createProjectRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type updateProjectRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type updateProjectProgressRequest struct {
	Status   string `json:"status" binding:"required"`
	Progress int    `json:"progress" binding:"required"`
}

func RegisterProjectRoutes(r *gin.Engine) {
	admin := r.Group("/api/projects")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.POST("", CreateProject)
		admin.GET("", ListProjects)
		admin.PATCH(":id", UpdateProject)
		admin.PATCH(":id/progress", UpdateProjectProgress)
		admin.DELETE(":id", DeleteProject)
	}

	me := r.Group("/api/me/projects")
	me.Use(middleware.AuthMiddleware())
	{
		me.GET("", ListMyProjects)
		me.GET("/summary", ProjectSummary)
	}
}

func CreateProject(c *gin.Context) {
	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project := models.Project{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "pending",
		Progress:    0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create project"})
		return
	}

	logs.SendLog(req.UserID, "Admin created project: "+req.Title)
	c.JSON(http.StatusCreated, project)
}

func ListProjects(c *gin.Context) {
	var projects []models.Project
	if err := db.DB.Order("created_at desc").Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch projects"})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var req updateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project models.Project
	if err := db.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	if req.Title != nil {
		project.Title = *req.Title
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	project.UpdatedAt = time.Now()

	if err := db.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update project"})
		return
	}
	logs.SendLog(project.UserID, "Admin updated project: "+project.Title)
	c.JSON(http.StatusOK, project)
}

func UpdateProjectProgress(c *gin.Context) {
	id := c.Param("id")

	var req updateProjectProgressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Progress < 0 || req.Progress > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "progress must be between 0 and 100"})
		return
	}

	var project models.Project
	if err := db.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	project.Progress = req.Progress
	project.Status = req.Status
	project.UpdatedAt = time.Now()

	if err := db.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update progress"})
		return
	}

	logs.SendLog(project.UserID, "Project progress updated: "+project.Title)
	c.JSON(http.StatusOK, project)
}

func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	if err := db.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	if err := db.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete project"})
		return
	}

	logs.SendLog(project.UserID, "Admin deleted project: "+project.Title)
	c.JSON(http.StatusOK, gin.H{"message": "project deleted"})
}

// ✅ Получение всех проектов конкретного пользователя
func ListMyProjects(c *gin.Context) {
	userID := c.MustGet(middleware.ContextUserID).(uint)

	var projects []models.Project
	if err := db.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func ProjectSummary(c *gin.Context) {
	userID := c.MustGet(middleware.ContextUserID).(uint)

	var total, completed, inProgress int64

	db.DB.Model(&models.Project{}).Where("user_id = ?", userID).Count(&total)
	db.DB.Model(&models.Project{}).Where("user_id = ? AND progress = 100", userID).Count(&completed)
	db.DB.Model(&models.Project{}).Where("user_id = ? AND progress < 100", userID).Count(&inProgress)

	summary := gin.H{
		"total":       total,
		"completed":   completed,
		"in_progress": inProgress,
	}

	c.JSON(http.StatusOK, summary)
}
