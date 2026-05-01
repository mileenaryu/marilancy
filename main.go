package main

import (
	"net/http"
	"os"

	"marilancy/config"
	"marilancy/models"
	"marilancy/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	config.ConnectDB()
	config.InitJWT()

	config.DB.AutoMigrate(
		&models.Freelancer{},
		&models.Client{},
		&models.Admin{},
		&models.Job{},
		&models.Application{},
		&models.Project{},
		&models.Task{},
		&models.Rating{},
	)
	config.SeedAdmin()

	r := gin.Default()

	r.Static("/static", "./static")
	r.Static("/uploads", "./uploads")
	r.LoadHTMLGlob("templates/*")

	routes.SetupRoutes(r)

	r.GET("/guest", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard_guest.html", nil)
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/guest")
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	r.GET("/freelancer", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard_freelancer.html", nil)
	})

	r.GET("/client", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard_client.html", nil)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
