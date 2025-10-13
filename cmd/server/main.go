package main

import (
	"github.com/gin-gonic/gin"
	"github.com/turtacn/SQLTraceBench/internal/app/api"
)

func setupServer() *gin.Engine {
	r := gin.Default()

	// Create the job store, metrics exporter, and worker.
	jobStore := api.NewJobStore()
	metricsExporter := api.NewMetricsExporter()
	worker := api.NewWorker(jobStore, metricsExporter)

	// Register the routes and pass the job store to the handlers.
	api.RegisterRoutes(r, jobStore)

	// Start the worker.
	worker.Start()

	return r
}

func main() {
	r := setupServer()
	r.Run(":8080")
}