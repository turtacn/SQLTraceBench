package main

import (
	"github.com/gin-gonic/gin"
	"github.com/turtacn/SQLTraceBench/internal/app/api"
)

func main() {
	r := gin.Default()

	// Create the job store, metrics exporter, and worker.
	jobStore := api.NewJobStore()
	metricsExporter := api.NewMetricsExporter()
	worker := api.NewWorker(jobStore, metricsExporter)

	// Register the routes and pass the job store to the handlers.
	api.RegisterRoutes(r, jobStore)

	// Start the worker.
	worker.Start()

	r.Run(":8080")
}