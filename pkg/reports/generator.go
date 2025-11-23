package reports

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// Generator handles report generation.
type Generator struct{}

// NewGenerator creates a new report generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateJSONReport generates a JSON report.
func (g *Generator) GenerateJSONReport(report *models.Report, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

// GenerateHTMLReport generates an HTML report.
func (g *Generator) GenerateHTMLReport(report *models.Report, outputPath string) error {
	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"sub": func(a, b int64) int64 { return a - b },
		"subf": func(a, b float64) float64 { return a - b },
	}).Parse(HTMLTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, report)
}

// CompareAndGenerateReports compares two runs and generates both JSON and HTML reports.
func (g *Generator) CompareAndGenerateReports(
	baseMetrics, candMetrics *models.PerformanceMetrics,
	metadata *models.ReportMetadata,
	outputPrefix string,
) (*models.Report, error) {

	// Calculate statistics
	pass := true
	if baseMetrics.QPS() > 0 {
		pass = candMetrics.QPS() >= baseMetrics.QPS()*(1-metadata.Threshold)
	}

	var reason string
	qpsDiff := candMetrics.QPS() - baseMetrics.QPS()
	var qpsDiffPercent float64
	if baseMetrics.QPS() > 0 {
		qpsDiffPercent = (qpsDiff / baseMetrics.QPS()) * 100
	}

	if pass {
		reason = fmt.Sprintf(
			"Validation passed. Candidate QPS of %.2f is within the %.2f%% threshold of the base QPS of %.2f (difference of %.2f, %.2f%%).",
			candMetrics.QPS(),
			metadata.Threshold*100,
			baseMetrics.QPS(),
			qpsDiff,
			qpsDiffPercent,
		)
	} else {
		reason = fmt.Sprintf(
			"Validation failed. Candidate QPS of %.2f is below the %.2f%% threshold of the base QPS of %.2f (difference of %.2f, %.2f%%).",
			candMetrics.QPS(),
			metadata.Threshold*100,
			baseMetrics.QPS(),
			qpsDiff,
			qpsDiffPercent,
		)
	}

	report := &models.Report{
		Version:   "report.v1",
		Timestamp: time.Now(),
		Metadata:  metadata,
		Result: &models.ValidationResult{
			BaseMetrics:      baseMetrics,
			CandidateMetrics: candMetrics,
			Pass:             pass,
			Reason:           reason,
		},
	}

	if err := g.GenerateJSONReport(report, outputPrefix+".json"); err != nil {
		return nil, fmt.Errorf("failed to generate JSON report: %w", err)
	}

	if err := g.GenerateHTMLReport(report, outputPrefix+".html"); err != nil {
		return nil, fmt.Errorf("failed to generate HTML report: %w", err)
	}

	return report, nil
}
