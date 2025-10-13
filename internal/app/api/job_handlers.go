package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// JobHandler handles job-related API requests.
type JobHandler struct {
	store *JobStore
}

// NewJobHandler creates a new JobHandler.
func NewJobHandler(store *JobStore) *JobHandler {
	return &JobHandler{store: store}
}

// CreateJob creates a new benchmark job and adds it to the queue.
func (h *JobHandler) CreateJob(c *gin.Context) {
	var config models.Config
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job := &models.Job{
		Status:    models.JobStatusPending,
		Config:    &config,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdJob, err := h.store.Create(job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add the job to the queue for processing.
	JobQueue <- createdJob

	c.JSON(http.StatusAccepted, createdJob)
}

// GetJob retrieves a job by its ID.
func (h *JobHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	job, err := h.store.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}