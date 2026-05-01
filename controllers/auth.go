package controllers

import (
	"fmt"
	"marilancy/config"
	"marilancy/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Nama     string `json:"nama"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal hash password"})
		return
	}

	switch input.Role {

	case "freelancer":
		err = config.DB.Create(&models.Freelancer{
			Nama:     input.Nama,
			Email:    input.Email,
			Password: string(hash),
			Role:     "freelancer",
		}).Error

	case "client":
		err = config.DB.Create(&models.Client{
			NamaClient: input.Nama,
			Email:      input.Email,
			Password:   string(hash),
			Role:       "client",
		}).Error

	case "admin":
		err = config.DB.Create(&models.Admin{
			NamaAdmin: input.Nama,
			Email:     input.Email,
			Password:  string(hash),
			Role:      "admin",
		}).Error

	default:
		c.JSON(400, gin.H{"error": "Role tidak valid"})
		return
	}

	if err != nil {
		fmt.Println("❌ REGISTER ERROR:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Register berhasil"})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var userID uint
	var email, pass, role string

	var client models.Client
	if err := config.DB.Where("email = ?", input.Email).First(&client).Error; err == nil {
		userID = client.ID
		email = client.Email
		pass = client.Password
		role = "client"
	}

	if role == "" {
		var freelancer models.Freelancer
		if err := config.DB.Where("email = ?", input.Email).First(&freelancer).Error; err == nil {
			userID = freelancer.ID
			email = freelancer.Email
			pass = freelancer.Password
			role = "freelancer"
		}
	}

	if role == "" {
		var admin models.Admin
		if err := config.DB.Where("email = ?", input.Email).First(&admin).Error; err == nil {
			userID = admin.ID
			email = admin.Email
			pass = admin.Password
			role = "admin"
		}
	}

	if role == "" {
		c.JSON(400, gin.H{"error": "User tidak ditemukan"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(input.Password)); err != nil {
		c.JSON(400, gin.H{"error": "Password salah"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal generate token"})
		return
	}

	c.JSON(200, gin.H{
		"token": tokenString,
		"role":  role,
	})
}
