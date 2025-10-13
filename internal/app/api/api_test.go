package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	jobStore := NewJobStore()
	RegisterRoutes(r, jobStore)
	return r
}

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"pong"}`, w.Body.String())
}

func TestCreateAndGetJob(t *testing.T) {
	router := setupRouter()

	// Create a job
	config := models.Config{TracePath: "traces.jsonl"}
	body, _ := json.Marshal(config)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var createdJob models.Job
	err := json.Unmarshal(w.Body.Bytes(), &createdJob)
	require.NoError(t, err)
	assert.NotEmpty(t, createdJob.ID)
	assert.Equal(t, models.JobStatusPending, createdJob.Status)

	// Get the job
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/jobs/"+createdJob.ID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var fetchedJob models.Job
	err = json.Unmarshal(w.Body.Bytes(), &fetchedJob)
	require.NoError(t, err)
	assert.Equal(t, createdJob.ID, fetchedJob.ID)
}