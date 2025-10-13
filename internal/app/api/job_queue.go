package api

import "github.com/turtacn/SQLTraceBench/internal/domain/models"

// JobQueue is a simple in-memory job queue.
var JobQueue = make(chan *models.Job, 100)