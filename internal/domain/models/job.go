package models

import "time"

// JobStatus represents the status of a benchmark job.
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

// Job represents a benchmark job.
type Job struct {
	ID        string        `json:"id"`
	Status    JobStatus     `json:"status"`
	Config    *Config       `json:"config"`
	Report    *Report       `json:"report,omitempty"`
	Error     string        `json:"error,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// Config holds the configuration for a benchmark job.
// This is a simplified version of the main application config.
type Config struct {
	TracePath       string `json:"trace_path"`
	SchemaPath      string `json:"schema_path"`
	WorkloadPath    string `json:"workload_path"`
	BaseMetricsPath string `json:"base_metrics_path"`
	CandMetricsPath string `json:"cand_metrics_path"`
	ReportPath      string `json:"report_path"`
	Executor        string `json:"executor"`
	Target          string `json:"target"`
}