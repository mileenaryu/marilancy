package controllers

import (
	"marilancy/config"
	"marilancy/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminDashboardData(c *gin.Context) {
	var totalFreelancers, totalClients, totalJobs int64
	config.DB.Model(&models.Freelancer{}).Count(&totalFreelancers)
	config.DB.Model(&models.Client{}).Count(&totalClients)
	config.DB.Model(&models.Job{}).Where("status != ?", "dihapus").Count(&totalJobs)

	c.JSON(http.StatusOK, gin.H{
		"freelancers": totalFreelancers,
		"clients":     totalClients,
		"jobs":        totalJobs,
	})
}

func GetFreelancers(c *gin.Context) {
	var users []models.Freelancer
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil freelancer"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetClients(c *gin.Context) {
	var clients []models.Client
	if err := config.DB.Find(&clients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil client"})
		return
	}
	c.JSON(http.StatusOK, clients)
}

func DeleteFreelancer(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Freelancer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Berhasil dihapus"})
}

func DeleteClient(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Client{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Berhasil dihapus"})
}

func AdminGetJobs(c *gin.Context) {
	var jobs []models.Job
	if err := config.DB.Preload("Client").Where("status != ?", "dihapus").Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil job"})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func DeleteJobs(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Model(&models.Job{}).Where("id = ?", id).Update("status", "dihapus").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus job"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Berhasil dihapus"})
}
