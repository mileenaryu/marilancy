package config

import (
	"fmt"
	"marilancy/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {

	var admin models.Admin

	err := DB.Where("email = ?", "admin@marilancy.com").First(&admin).Error

	if err == nil {
		fmt.Println("✅ Admin sudah ada, skip seeding")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)

	admin = models.Admin{
		NamaAdmin: "Super Admin",
		Email:     "admin@marilancy.com",
		Password:  string(hash),
		Role:      "admin",
	}

	if err := DB.Create(&admin).Error; err != nil {
		fmt.Println("❌ Gagal seed admin:", err)
		return
	}

	fmt.Println("🔥 Admin berhasil dibuat!")
	fmt.Println("📧 Email: admin@marilancy.com")
	fmt.Println("🔑 Password: admin123")
}
