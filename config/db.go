package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	host := os.Getenv("MYSQLHOST")
	port := os.Getenv("MYSQLPORT")
	user := os.Getenv("MYSQLUSER")
	pass := os.Getenv("MYSQLPASSWORD")
	name := os.Getenv("MYSQLDATABASE")

	// fallback untuk local
	if host == "" {
		host = "127.0.0.1"
		port = "3306"
		user = "root"
		pass = ""
		name = "marilancy"
	}

	// validasi
	if host == "" || user == "" || name == "" {
		log.Fatal("Database ENV not set properly")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB ERROR:", err)
	}

	// connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("DB instance error:", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	DB = db
}
