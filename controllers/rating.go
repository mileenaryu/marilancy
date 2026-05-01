package controllers

import (
	"errors"
	"marilancy/config"
	"marilancy/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateRating(c *gin.Context) {
	var input struct {
		Nilai        float64 `json:"nilai"`
		Komentar     string  `json:"komentar"`
		FreelancerID uint    `json:"freelancer_id"`
		ProjectID    uint    `json:"project_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}


	if input.Nilai < 1 || input.Nilai > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nilai harus 1 - 5"})
		return
	}

	clientIDRaw, _ := c.Get("user_id")
	clientID := clientIDRaw.(uint)


	var existing models.Rating
	err := config.DB.
		Where("client_id = ? AND freelancer_id = ? AND project_id = ?", clientID, input.FreelancerID, input.ProjectID).
		First(&existing).Error


	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah memberi penilaian untuk proyek ini"})
		return
	}


	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal cek rating"})
		return
	}

	rating := models.Rating{
		Nilai:        input.Nilai,
		Komentar:     input.Komentar,
		ClientID:     clientID,
		FreelancerID: input.FreelancerID,
		ProjectID:    input.ProjectID,
		CreatedAt:    time.Now(),
	}

	if err := config.DB.Create(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating berhasil disimpan"})
}
func CheckRating(c *gin.Context) {
	freelancerIDStr := c.Query("freelancer_id")
	projectIDStr := c.Query("project_id")

	freelancerID, _ := strconv.Atoi(freelancerIDStr)
	projectID, _ := strconv.Atoi(projectIDStr)

	clientIDRaw, _ := c.Get("user_id")
	clientID := clientIDRaw.(uint)

	var rating models.Rating

	err := config.DB.
		Where("client_id = ? AND freelancer_id = ? AND project_id = ?", clientID, freelancerID, projectID).
		First(&rating).Error


	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{"sudah_rating": false})
		return
	}


	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal cek rating"})
		return
	}


	c.JSON(http.StatusOK, gin.H{"sudah_rating": true})
}
func GetFreelancerRatingSummary(c *gin.Context) {
	freelancerID := c.Param("id")

	var total int64
	var avg float64

	config.DB.Model(&models.Rating{}).
		Where("freelancer_id = ?", freelancerID).
		Count(&total)

	config.DB.Model(&models.Rating{}).
		Select("AVG(nilai)").
		Where("freelancer_id = ?", freelancerID).
		Scan(&avg)


	if total == 0 {
		avg = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"total_penilaian": total,
		"avg_rating":      avg,
	})
}
