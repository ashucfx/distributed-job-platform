package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func SetupRouter(jobController *JobController) *gin.Engine {
	r := gin.Default()

	// CORS for the dashboard
	r.Use(cors.Default())

	// Generic error recovery
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		jobs := v1.Group("/jobs")
		{
			jobs.POST("", jobController.CreateJob)
			jobs.GET("", jobController.ListJobs)
			jobs.GET("/:id", jobController.GetJob)
		}
	}

	return r
}
