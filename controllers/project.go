package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"marilancy/config"
	"marilancy/models"

	"github.com/gin-gonic/gin"
)

func GetProjectDetail(c *gin.Context) {
	projectID := c.Param("id")
	var project models.Project

	if err := config.DB.Preload("Job").Preload("Client").Preload("Freelancer").Preload("Tasks").Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project tidak ditemukan"})
		return
	}

	totalTasks := len(project.Tasks)
	completedTasks := 0
	for _, task := range project.Tasks {
		if task.Status == "done" {
			completedTasks++
		}
	}

	progress := 0
	if totalTasks > 0 {
		progress = (completedTasks * 100) / totalTasks
	}

	c.JSON(http.StatusOK, gin.H{
		"project":  project,
		"progress": progress,
	})
}

func CreateTask(c *gin.Context) {
	var input struct {
		ProjectID uint   `json:"project_id"`
		Title     string `json:"title"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := models.Task{
		ProjectID: input.ProjectID,
		Title:     input.Title,
		Status:    "todo",
	}

	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task ditambahkan", "task": task})
}

func UpdateTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	var input struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status diperbarui"})
}

func CompleteProject(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project
	if err := config.DB.Preload("Tasks").Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project tidak ditemukan"})
		return
	}

	totalTasks := len(project.Tasks)
	completedTasks := 0
	for _, task := range project.Tasks {
		if task.Status == "done" {
			completedTasks++
		}
	}

	if totalTasks == 0 || completedTasks < totalTasks {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua task harus selesai 100% terlebih dahulu!"})
		return
	}

	link := c.PostForm("submission_link")
	file, errFile := c.FormFile("submission_file")

	if link == "" && errFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Harap sertakan link atau upload file hasil kerja!"})
		return
	}

	if errFile == nil {
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		savePath := filepath.Join("uploads", fileName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file hasil kerja"})
			return
		}
		project.SubmissionFile = "/" + filepath.ToSlash(savePath)
	}

	project.SubmissionLink = link
	project.Status = "completed"

	if err := config.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyelesaikan proyek"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proyek berhasil diselesaikan"})
}

func RequestRevision(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project
	if err := config.DB.Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project tidak ditemukan"})
		return
	}

	project.Status = "active"

	if err := config.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meminta pengiriman ulang"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proyek dibuka kembali untuk revisi"})
}

func GetMyProjects(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var projects []models.Project
	if err := config.DB.Preload("Job").Preload("Client").Preload("Tasks").Where("freelancer_id = ?", userID).Order("created_at desc").Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data proyek"})
		return
	}

	var result []gin.H
	for _, p := range projects {
		totalTasks := len(p.Tasks)
		completedTasks := 0
		for _, t := range p.Tasks {
			if t.Status == "done" {
				completedTasks++
			}
		}

		progress := 0
		if totalTasks > 0 {
			progress = (completedTasks * 100) / totalTasks
		}

		result = append(result, gin.H{
			"id":       p.ID,
			"status":   p.Status,
			"job":      p.Job,
			"client":   p.Client,
			"progress": progress,
			"tasks":    p.Tasks,
		})
	}

	c.JSON(http.StatusOK, result)
}

func GetClientProjects(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var projects []models.Project
	if err := config.DB.Preload("Job").Preload("Freelancer").Preload("Tasks").Where("client_id = ?", userID).Order("created_at desc").Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data proyek"})
		return
	}

	var result []gin.H
	for _, p := range projects {
		totalTasks := len(p.Tasks)
		completedTasks := 0
		for _, t := range p.Tasks {
			if t.Status == "done" {
				completedTasks++
			}
		}

		progress := 0
		if totalTasks > 0 {
			progress = (completedTasks * 100) / totalTasks
		}

		result = append(result, gin.H{
			"id":         p.ID,
			"status":     p.Status,
			"job":        p.Job,
			"freelancer": p.Freelancer,
			"progress":   progress,
		})
	}

	c.JSON(http.StatusOK, result)
}

func DeleteTask(c *gin.Context) {
	taskID := c.Param("task_id")

	if err := config.DB.Delete(&models.Task{}, taskID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task berhasil dihapus"})
}

func UpdateTaskTitle(c *gin.Context) {
	taskID := c.Param("task_id")
	var input struct {
		Title string `json:"title"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if err := config.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("title", input.Title).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate judul task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Judul task diperbarui"})
}
func UpdateTaskPriority(c *gin.Context) {
	taskID := c.Param("task_id")
	var input struct {
		Priority string `json:"priority"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if err := config.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("priority", input.Priority).Error; err != nil {
		c.JSON(500, gin.H{"error": "Gagal update prioritas"})
		return
	}

	c.JSON(200, gin.H{"message": "Prioritas diperbarui"})
}
