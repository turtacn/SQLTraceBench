package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/turtacn/SQLTraceBench/internal/app/conversion"
	"github.com/turtacn/SQLTraceBench/internal/app/execution"
	"github.com/turtacn/SQLTraceBench/internal/app/generation"
	"github.com/turtacn/SQLTraceBench/internal/app/validation"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/reporters"
	"github.com/turtacn/SQLTraceBench/internal/utils"
	"github.com/turtacn/SQLTraceBench/internal/utils/progress"
	"github.com/turtacn/SQLTraceBench/internal/utils/terminal"
)

// Manager coordinates the workflow pipeline.
type Manager struct {
	conversionSvc conversion.Service
	generationSvc generation.Service
	executionSvc  execution.Service
	validationSvc validation.Service
	logger        *utils.Logger
}

// NewManager creates a new WorkflowManager.
func NewManager(
	c conversion.Service,
	g generation.Service,
	e execution.Service,
	v validation.Service,
) *Manager {
	return &Manager{
		conversionSvc: c,
		generationSvc: g,
		executionSvc:  e,
		validationSvc: v,
		logger:        utils.GetGlobalLogger(),
	}
}

// Run executes the full 4-phase pipeline.
func (m *Manager) Run(ctx context.Context, cfg WorkflowConfig) error {
	m.logger.Info("Starting Workflow", utils.Field{Key: "config", Value: cfg})
	fmt.Println(terminal.Info("Starting SQLTraceBench Workflow..."))

	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	// ==========================================
	// Phase 1: Conversion
	// ==========================================
	m.logger.Info("Phase 1: Conversion starting...")
	p1Bar := progress.NewProgressBar(100, "Phase 1: Conversion") // Estimation

	// 1.1 Trace Conversion
	traceReq := conversion.ConvertTraceRequest{
		SourcePath:   cfg.InputTracePath,
		TargetDBType: cfg.TargetPlugin,
	}

	// Simulation of progress for conversion (since streaming isn't fully exposed with progress callback yet)
	p1Bar.Increment(10)
	convRes, err := m.conversionSvc.ConvertFromFile(ctx, traceReq)
	if err != nil {
		return fmt.Errorf("conversion phase failed (traces): %w", err)
	}
	p1Bar.Increment(50)

	// Save converted traces
	convertedTracePath := filepath.Join(cfg.OutputDir, "converted", "traces.jsonl")
	if err := saveJSONL(convertedTracePath, convRes.Traces); err != nil {
		return fmt.Errorf("failed to save converted traces: %w", err)
	}
	p1Bar.Increment(20)

	// 1.2 Schema Conversion (if schema path provided)
	if cfg.InputSchemaPath != "" {
		schemaOutPath := filepath.Join(cfg.OutputDir, "converted", "schema.sql")
		schemaReq := conversion.ConvertRequest{
			SourceSchemaPath: cfg.InputSchemaPath,
			TargetDBType:     cfg.TargetPlugin,
			OutputPath:       schemaOutPath,
		}
		// Ensure output dir exists
		os.MkdirAll(filepath.Dir(schemaOutPath), 0755)

		if err := m.conversionSvc.ConvertSchemaFromFile(ctx, schemaReq); err != nil {
			return fmt.Errorf("conversion phase failed (schema): %w", err)
		}
	}
	p1Bar.Finish()
	fmt.Println(terminal.Success("Phase 1: Conversion complete"))
	m.logger.Info("Phase 1: Conversion complete")

	// ==========================================
	// Phase 2: Generation
	// ==========================================
	m.logger.Info("Phase 2: Generation starting...")
	p2Bar := progress.NewProgressBar(int64(cfg.Generation.Count), "Phase 2: Generation")

	// Update Generation Request with converted traces
	genReq := cfg.Generation
	genReq.SourceTraces = convRes.Traces

	// TODO: Add progress callback to generation service if possible, currently we wait
	workload, err := m.generationSvc.GenerateWorkload(ctx, genReq)
	if err != nil {
		return fmt.Errorf("generation phase failed: %w", err)
	}
	// Complete the bar
	p2Bar.Increment(int64(cfg.Generation.Count))
	p2Bar.Finish()

	workloadPath := filepath.Join(cfg.OutputDir, "workload", "benchmark.jsonl")
	if err := saveJSONL(workloadPath, workload); err != nil {
		return fmt.Errorf("failed to save workload: %w", err)
	}
	fmt.Println(terminal.Success("Phase 2: Generation complete"))
	m.logger.Info("Phase 2: Generation complete")

	// ==========================================
	// Phase 3: Execution
	// ==========================================
	m.logger.Info("Phase 3: Execution starting...")
	totalQueries := int64(len(workload.Queries))
	p3Bar := progress.NewProgressBar(totalQueries, "Phase 3: Execution ")

	execCfg := cfg.Execution
	// Ensure TargetDB is set from top-level config if not in sub-config
	if execCfg.TargetDB == "" {
		execCfg.TargetDB = cfg.TargetPlugin
	}

	// We might need to wrap execution to update progress, but ExecutionService is black box here.
	// For now, we just indicate start and end. Ideally we'd pass a progress channel.
	p3Bar.Increment(1) // Started

	result, err := m.executionSvc.RunBenchmark(ctx, workload, execCfg)
	if err != nil {
		return fmt.Errorf("execution phase failed: %w", err)
	}

	p3Bar.Increment(totalQueries) // Done
	p3Bar.Finish()

	resultPath := filepath.Join(cfg.OutputDir, "results", "metrics.json")
	if err := saveJSON(resultPath, result); err != nil {
		return fmt.Errorf("failed to save metrics: %w", err)
	}
	fmt.Println(terminal.Success("Phase 3: Execution complete"))
	m.logger.Info("Phase 3: Execution complete")

	// ==========================================
	// Phase 4: Validation & Reporting
	// ==========================================
	if cfg.BaselineMetricsPath != "" {
		m.logger.Info("Phase 4: Validation starting...")
		fmt.Println(terminal.Info("Phase 4: Validation starting..."))

		// Load baseline
		var baseline models.BenchmarkResult
		if err := loadJSON(cfg.BaselineMetricsPath, &baseline); err != nil {
			m.logger.Warn("Failed to load baseline metrics, skipping validation", utils.Field{Key: "error", Value: err})
			fmt.Println(terminal.Warning("Skipping validation: could not load baseline metrics"))
		} else {
			report, err := m.validationSvc.ValidateBenchmarks(ctx, &baseline, result)
			if err != nil {
				return fmt.Errorf("validation phase failed: %w", err)
			}

			// Generate HTML Report
			htmlReporter, err := reporters.NewHTMLReporter()
			if err != nil {
				m.logger.Error("Failed to initialize HTML reporter", utils.Field{Key: "error", Value: err})
			} else {
				htmlPath := filepath.Join(cfg.OutputDir, "validation_report.html")
				if err := htmlReporter.GenerateReport(report, cfg.TargetPlugin, htmlPath); err != nil {
					m.logger.Error("Failed to generate HTML report", utils.Field{Key: "error", Value: err})
				} else {
					fmt.Println(terminal.Success(fmt.Sprintf("HTML Report generated: %s", htmlPath)))
				}
			}

			// Save JSON report
			reportPath := filepath.Join(cfg.OutputDir, "report.json")
			if err := saveJSON(reportPath, report); err != nil {
				return fmt.Errorf("failed to save validation report: %w", err)
			}

			statusMsg := fmt.Sprintf("Phase 4: Validation complete. Status: %s", utils.SafeString(report.Pass))
			if report.Pass {
				fmt.Println(terminal.Success(statusMsg))
			} else {
				fmt.Println(terminal.Error(statusMsg))
			}
			m.logger.Info("Phase 4: Validation complete", utils.Field{Key: "status", Value: report.Pass})
		}
	} else {
		m.logger.Info("Phase 4: Validation skipped (no baseline provided)")
		fmt.Println(terminal.Warning("Phase 4: Validation skipped (no baseline provided)"))
	}

	return nil
}

func saveJSONL(path string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)

	switch v := data.(type) {
	case []models.SQLTrace:
		for _, t := range v {
			if err := enc.Encode(t); err != nil {
				return err
			}
		}
	case *models.BenchmarkWorkload:
		for _, q := range v.Queries {
			if err := enc.Encode(q); err != nil {
				return err
			}
		}
	default:
		return enc.Encode(data)
	}
	return nil
}

func saveJSON(path string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func loadJSON(path string, dest interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(dest)
}
