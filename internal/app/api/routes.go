package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the API routes with the Gin router.
func RegisterRoutes(r *gin.Engine, jobStore *JobStore) {
	// Simple ping endpoint for health checks.
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Job management endpoints.
	jobHandler := NewJobHandler(jobStore)

	jobs := r.Group("/jobs")
	{
		jobs.POST("", jobHandler.CreateJob)
		jobs.GET("/:id", jobHandler.GetJob)
	}

	// Metrics endpoint for Prometheus scraping.
	metricsExporter := NewMetricsExporter()
	r.GET("/metrics", metricsExporter.Handler())
}