package api

import (
	"sync"

	"github.com/google/uuid"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

// JobStore is an in-memory store for benchmark jobs.
type JobStore struct {
	mu   sync.RWMutex
	jobs map[string]*models.Job
}

// NewJobStore creates a new JobStore.
func NewJobStore() *JobStore {
	return &JobStore{
		jobs: make(map[string]*models.Job),
	}
}

// Create creates a new job and adds it to the store.
func (s *JobStore) Create(job *models.Job) (*models.Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job.ID = uuid.New().String()
	s.jobs[job.ID] = job
	return job, nil
}

// Get retrieves a job from the store by its ID.
func (s *JobStore) Get(id string) (*models.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, ok := s.jobs[id]
	if !ok {
		return nil, types.NewError(types.ErrNotFound, "job not found")
	}
	return job, nil
}

// Update updates an existing job in the store.
func (s *JobStore) Update(job *models.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs[job.ID] = job
	return nil
}