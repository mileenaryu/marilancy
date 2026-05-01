package controllers

import (
	"net/http"

	"marilancy/config"
	"marilancy/models"

	"github.com/gin-gonic/gin"
)

func ApplyJob(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		JobID uint `json:"job_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	var freelancer models.Freelancer
	if err := config.DB.First(&freelancer, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Freelancer tidak ditemukan"})
		return
	}
	if freelancer.Resume == "" || freelancer.Certificates == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lengkapi resume & certificates dulu"})
		return
	}

	var job models.Job
	if err := config.DB.First(&job, input.JobID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job tidak ditemukan"})
		return
	}
	if job.Status == "dihapus" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job sudah dihapus"})
		return
	}
	if job.Status == "ditutup" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job sudah ditutup"})
		return
	}

	var existing models.Application
	err := config.DB.Where("freelancer_id = ? AND job_id = ?", userID, input.JobID).First(&existing).Error

	if err == nil {

		if existing.Status == "withdrawn" {
			existing.Status = "pending"

			if err := config.DB.Save(&existing).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal re-apply"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Berhasil melamar ulang job"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "Sudah pernah melamar job ini"})
		return
	}

	app := models.Application{
		FreelancerID: userID,
		JobID:        input.JobID,
		Status:       "pending",
	}

	if err := config.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melamar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil melamar job"})
}

func WithdrawApplication(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		JobID uint `json:"job_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var app models.Application
	err := config.DB.Where("freelancer_id = ? AND job_id = ?", userID, input.JobID).First(&app).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	if app.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak bisa withdraw, status sudah berubah"})
		return
	}

	app.Status = "withdrawn"

	if err := config.DB.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal withdraw"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil mengundurkan diri"})
}

func GetJobApplicants(c *gin.Context) {
	jobID := c.Param("id")

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var job models.Job
	if err := config.DB.First(&job, jobID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job tidak ditemukan"})
		return
	}

	if job.ClientID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tidak bisa lihat pelamar job orang lain"})
		return
	}

	var apps []models.Application
	err := config.DB.
		Preload("Freelancer").
		Where("job_id = ?", jobID).
		Order("tanggal_daftar desc").
		Find(&apps).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil pelamar"})
		return
	}

	c.JSON(http.StatusOK, apps)
}

func UpdateApplicationStatus(c *gin.Context) {

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		ApplicationID uint   `json:"application_id"`
		Status        string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	validStatus := map[string]bool{
		"pending":  true,
		"accepted": true,
		"rejected": true,
		"dihapus":  true,
	}

	if !validStatus[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak valid"})
		return
	}

	var app models.Application
	if err := config.DB.First(&app, input.ApplicationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	var job models.Job
	if err := config.DB.First(&job, app.JobID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job tidak ditemukan"})
		return
	}
	if job.ClientID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tidak diizinkan mengubah job klien lain"})
		return
	}

	app.Status = input.Status

	if err := config.DB.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status"})
		return
	}

	if input.Status == "accepted" {
		var existingProject models.Project

		err := config.DB.Where("job_id = ? AND freelancer_id = ?", app.JobID, app.FreelancerID).First(&existingProject).Error

		if err != nil {

			newProject := models.Project{
				JobID:        app.JobID,
				ClientID:     job.ClientID,
				FreelancerID: app.FreelancerID,
				Status:       "active",
			}
			config.DB.Create(&newProject)

		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status updated",
		"data":    app,
	})
}
