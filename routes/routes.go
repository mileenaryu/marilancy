package routes

import (
	"marilancy/controllers"
	"marilancy/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.GET("/jobs", controllers.GetJobs)
	r.GET("/jobs/:id", controllers.GetJobDetail)
	r.GET("/clients/:id", controllers.GetClientByID)

	r.GET("/freelancer/:id/rating-summary", controllers.GetFreelancerRatingSummary)

	r.GET("/client/rating", func(c *gin.Context) {
		c.HTML(200, "rating.html", nil)
	})
	r.GET("/freelancer/projects", func(c *gin.Context) {
		c.HTML(200, "my-projects.html", nil)
	})

	r.GET("/client/projects", func(c *gin.Context) {
		c.HTML(200, "my-projects-client.html", nil)
	})
	r.GET("/lihatprofileclient.html", func(c *gin.Context) {
		c.HTML(200, "lihatprofileclient.html", nil) 
	})
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(200, "dashboard_admin.html", nil)
	})

	r.GET("/job/detail", func(c *gin.Context) {
		c.HTML(200, "job_detail.html", nil)
	})

	r.GET("/client/profile/view", func(c *gin.Context) {
		c.HTML(200, "profile.html", nil)
	})

	r.GET("/freelancer/profile/view", func(c *gin.Context) {
		c.HTML(200, "profilefree.html", nil)
	})
	r.GET("/client/applicants", func(c *gin.Context) {
		c.HTML(200, "pendaftar.html", nil)
	})
	r.GET("/freelancer/notification", func(c *gin.Context) {
		c.HTML(200, "notification.html", nil)
	})


	r.GET("/freelancer/workspace", func(c *gin.Context) {
		c.HTML(200, "workspace-freelancer.html", nil)
	})
	r.GET("/client/workspace", func(c *gin.Context) {
		c.HTML(200, "workspace-client.html", nil)
	})
	
	apiProjects := r.Group("/api/projects")
	apiProjects.Use(middleware.AuthMiddleware())
	{
		apiProjects.GET("/:id", controllers.GetProjectDetail)
		apiProjects.POST("/task", controllers.CreateTask)
		apiProjects.PUT("/task/:task_id", controllers.UpdateTaskStatus)
		apiProjects.PUT("/:id/complete", controllers.CompleteProject)
		apiProjects.PUT("/:id/revision", controllers.RequestRevision)
		apiProjects.DELETE("/task/:task_id", controllers.DeleteTask) 
		apiProjects.PATCH("/task/:task_id/title", controllers.UpdateTaskTitle)
		apiProjects.PUT("/task/:task_id/priority", controllers.UpdateTaskPriority)
	}

	freelancer := r.Group("/freelancer")
	freelancer.Use(
		middleware.AuthMiddleware(),
		middleware.RoleMiddleware("freelancer"),
	)
	{
		freelancer.GET("/profile", controllers.GetFreelancerProfile)
		freelancer.PUT("/profile", controllers.UpdateFreelancerProfile)
		freelancer.POST("/apply", controllers.ApplyJob)
		freelancer.GET("/applications", controllers.GetMyApplications)
		freelancer.POST("/withdraw", controllers.WithdrawApplication)
		freelancer.GET("/my-projects", controllers.GetMyProjects)
	}

	client := r.Group("/client")
	client.Use(
		middleware.AuthMiddleware(),
		middleware.RoleMiddleware("client"),
	)
	{
		client.GET("/profile", controllers.GetClientProfile)
		client.PUT("/profile", controllers.UpdateClientProfile)
		client.POST("/jobs", controllers.CreateJob)
		client.PUT("/jobs/:id", controllers.UpdateJob)
		client.GET("/jobs/:id/applicants", controllers.GetJobApplicants)
		client.PUT("/application/status", controllers.UpdateApplicationStatus)
		client.DELETE("/jobs/:id", controllers.DeleteJob)
		client.GET("/my-projects", controllers.GetClientProjects)
		client.POST("/rating", controllers.CreateRating)
		client.GET("/check-rating", controllers.CheckRating)
	}

	admin := r.Group("/admin")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.RoleMiddleware("admin"),
	)
	{
		admin.GET("/data", controllers.AdminDashboardData)

		admin.GET("/freelancers", controllers.GetFreelancers)
		admin.GET("/clients", controllers.GetClients)

		admin.DELETE("/freelancers/:id", controllers.DeleteFreelancer)
		admin.DELETE("/clients/:id", controllers.DeleteClient)

		admin.GET("/jobs", controllers.AdminGetJobs)
		admin.DELETE("/jobs/:id", controllers.DeleteJobs)
	}
}
