package api

import (
	"fmt"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/app"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// Worker is a background worker that processes benchmark jobs.
type Worker struct {
	store    *JobStore
	root     *app.Root
	exporter *MetricsExporter
}

// NewWorker creates a new Worker.
func NewWorker(store *JobStore, exporter *MetricsExporter) *Worker {
	return &Worker{
		store:    store,
		root:     app.NewRoot(),
		exporter: exporter,
	}
}

// Start starts the worker.
func (w *Worker) Start() {
	go func() {
		for job := range JobQueue {
			w.processJob(job)
		}
	}()
}

func (w *Worker) processJob(job *models.Job) {
	job.Status = models.JobStatusRunning
	job.UpdatedAt = time.Now()
	w.store.Update(job)

	// In a real implementation, we would use the job's config to run the pipeline.
	// For now, we'll just simulate the work and update the job status.
	fmt.Printf("Processing job %s...\n", job.ID)
	time.Sleep(5 * time.Second) // Simulate work

	// Here, you would call the actual benchmark pipeline services.
	// For example:
	// report, err := w.root.Validation.Validate(...)
	// if err != nil {
	// 	job.Status = models.JobStatusFailed
	// 	job.Error = err.Error()
	// } else {
	// 	job.Status = models.JobStatusCompleted
	// 	job.Report = report
	// 	w.exporter.RecordMetrics(job.Config.Executor, job.Config.Target, report.Result.CandidateMetrics)
	// }

	job.Status = models.JobStatusCompleted
	job.UpdatedAt = time.Now()
	w.store.Update(job)
	fmt.Printf("Job %s completed.\n", job.ID)
}