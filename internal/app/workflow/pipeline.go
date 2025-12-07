package workflow

import (
	"github.com/turtacn/SQLTraceBench/internal/app/execution"
	"github.com/turtacn/SQLTraceBench/internal/app/generation"
)

// WorkflowConfig defines the configuration for the full pipeline.
type WorkflowConfig struct {
	// Paths
	InputTracePath  string `yaml:"input_trace_path"`
	InputSchemaPath string `yaml:"input_schema_path"`
	OutputDir       string `yaml:"output_dir"`

	// Settings
	TargetPlugin string `yaml:"target_plugin"`
	ReportStyle  string `yaml:"report_style"` // html, json

	// Phase Configs
	Generation generation.GenerateRequest `yaml:"generation"`
	Execution  execution.ExecutionConfig  `yaml:"execution"`

	// Validation
	BaselineMetricsPath string `yaml:"baseline_metrics_path"`
}

// PipelineStatus tracks the progress of the workflow.
type PipelineStatus struct {
	CurrentPhase string
	Error        error
}
