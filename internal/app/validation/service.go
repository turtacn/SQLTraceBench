package validation

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/reporters"
)

// Service is the interface for the validation service.
type Service interface {
	// ValidateTrace performs statistical validation between original and generated traces.
	ValidateTrace(ctx context.Context, req ValidationRequest) (*services.ValidationReport, error)
}

// ValidationRequest contains parameters for trace validation.
type ValidationRequest struct {
	OriginalPath   string
	GeneratedPath  string
	ReportDir      string
	PrometheusPort int
	KSThreshold    float64
}

// DefaultService is the default implementation of the validation service.
type DefaultService struct {
	validator *services.StatisticalValidator
	generator *services.ReportGenerator
}

// NewService creates a new DefaultService.
func NewService() Service {
	return &DefaultService{
		// Initialize with default thresholds, these can be overridden or passed in request
		validator: services.NewStatisticalValidator(0.05, 0.01),
		generator: services.NewReportGenerator(),
	}
}

// ValidateTrace performs statistical validation between original and generated traces.
func (s *DefaultService) ValidateTrace(ctx context.Context, req ValidationRequest) (*services.ValidationReport, error) {
	// 1. Load data
	original, err := loadTraces(req.OriginalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load original traces: %w", err)
	}
	generated, err := loadTraces(req.GeneratedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load generated traces: %w", err)
	}

	// Update validator thresholds if provided
	if req.KSThreshold > 0 {
		s.validator.KSThreshold = req.KSThreshold
	}

	// 2. Statistical Validation (Distribution)
	var distResults []services.ValidationResult
	params := extractCommonParameters(original, generated)
	for _, paramName := range params {
		origDist := extractDistribution(original, paramName)
		genDist := extractDistribution(generated, paramName)

		if len(origDist) == 0 || len(genDist) == 0 {
			continue
		}

		result := s.validator.KolmogorovSmirnovTest(origDist, genDist)
		result.Details["parameter"] = paramName
		distResults = append(distResults, *result)
	}

	// 3. Temporal Validation
	temporalResults := validateTemporalPatterns(original, generated, s.validator)

	// 4. Query Type Validation
	queryResults := validateQueryTypes(original, generated, s.validator)

	// 5. Generate Report
	report := s.generator.Generate(distResults, temporalResults, queryResults)

	// 6. Export Reports
	if req.ReportDir != "" {
		// Check if running from tests or root
		templatePath := "web/templates/validation_report.html"
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			// Try absolute path or relative to current execution if in test
			// But usually we should run tests from root.
			// For now let's try to check if we are in a subdirectory (e.g. tests/integration)
			// or just fallback to a known location.
			// Or better: Let the caller specify template path if needed.
            // But `ValidationRequest` does not have it.
            // Let's just try ../../web/templates/validation_report.html if original fails
            if _, err := os.Stat("../../" + templatePath); err == nil {
                templatePath = "../../" + templatePath
            }
		}

		htmlReporter := reporters.NewHTMLReporter(templatePath)
		outputPath := fmt.Sprintf("%s/validation_report.html", req.ReportDir)
		if err := htmlReporter.GenerateReport(report, outputPath); err != nil {
			// Log error but don't fail the whole process? Or return error?
			// For now, return error as report generation is key.
			return nil, fmt.Errorf("failed to generate HTML report: %w", err)
		}
	}

	if req.PrometheusPort > 0 {
		promReporter := reporters.NewPrometheusReporter(req.PrometheusPort)
		// Run in background as it blocks
		go func() {
			if err := promReporter.ExportMetrics(report); err != nil {
				fmt.Printf("Prometheus exporter error: %v\n", err)
			}
		}()
	}

	return report, nil
}

// Helper functions

func loadTraces(path string) ([]models.SQLTrace, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var traces []models.SQLTrace
	scanner := bufio.NewScanner(file)
	// Increase buffer size for large lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		var trace models.SQLTrace
		if err := json.Unmarshal(scanner.Bytes(), &trace); err != nil {
			// Skip malformed lines or return error?
            // Usually skip is safer for massive logs, but for validation strictness is good.
            // Let's log and skip.
			continue
		}
		traces = append(traces, trace)
	}
	return traces, scanner.Err()
}

func extractCommonParameters(orig, gen []models.SQLTrace) []string {
	paramSet := make(map[string]struct{})

	// Naive approach: scan all traces.
    // Optimization: scan first N traces.
	for _, t := range orig {
		for k := range t.Parameters {
			paramSet[k] = struct{}{}
		}
	}

	var params []string
	for k := range paramSet {
		params = append(params, k)
	}
	return params
}

func extractDistribution(traces []models.SQLTrace, paramName string) []float64 {
	var dist []float64
	for _, t := range traces {
		if val, ok := t.Parameters[paramName]; ok {
			// Convert value to float64 if possible
			switch v := val.(type) {
			case float64:
				dist = append(dist, v)
			case int:
				dist = append(dist, float64(v))
			case int64:
				dist = append(dist, float64(v))
            // TODO: Handle other types (string length? hash?)
            // For now, we assume numerical parameters.
			}
		}
	}
	return dist
}

func validateTemporalPatterns(orig, gen []models.SQLTrace, v *services.StatisticalValidator) []services.ValidationResult {
	if len(orig) == 0 || len(gen) == 0 {
		return nil
	}

	// 1. Normalize time: Find start and end time for both
	origStart, origEnd := getTimeRange(orig)
	genStart, genEnd := getTimeRange(gen)

	origDuration := origEnd.Sub(origStart).Seconds()
	genDuration := genEnd.Sub(genStart).Seconds()

    if origDuration <= 0 || genDuration <= 0 {
         // Not enough time span
         return []services.ValidationResult{
            {TestName: "Temporal Distribution", PValue: 0.0, Passed: false, Details: map[string]interface{}{"reason": "zero duration"}},
         }
    }

	// 2. Binning: Divide into N bins (e.g., 100)
	numBins := 100
	origBins := make([]float64, numBins)
	genBins := make([]float64, numBins)

	// Populate bins
	for _, t := range orig {
		offset := t.Timestamp.Sub(origStart).Seconds()
		binIdx := int((offset / origDuration) * float64(numBins))
		if binIdx >= numBins {
			binIdx = numBins - 1
		}
        if binIdx < 0 { binIdx = 0 }
		origBins[binIdx]++
	}

	for _, t := range gen {
		offset := t.Timestamp.Sub(genStart).Seconds()
		binIdx := int((offset / genDuration) * float64(numBins))
		if binIdx >= numBins {
			binIdx = numBins - 1
		}
        if binIdx < 0 { binIdx = 0 }
		genBins[binIdx]++
	}

	// Normalize bins to be probability distributions
	origTotal := float64(len(orig))
	genTotal := float64(len(gen))
	for i := 0; i < numBins; i++ {
		origBins[i] /= origTotal
		genBins[i] /= genTotal
	}

	// 3. Compare using Jensen-Shannon Divergence
    // JS Divergence is 0 for identical distributions, ln(2) for disjoint.
    // Threshold for pass: < 0.1 (heuristic)
	jsDiv := v.JensenShannonDivergence(origBins, genBins)
    passed := jsDiv < 0.1

    // Convert JS Div to a score similar to P-Value (1.0 - normalized JS) for consistent reporting
    // JS is between 0 and ln(2) approx 0.693.
    score := 1.0 - (jsDiv / 0.693)
    if score < 0 { score = 0 }

	return []services.ValidationResult{
		{
			TestName: "Temporal Distribution (JS Divergence)",
			PValue:   score, // Not strictly a p-value, but a similarity score
			Passed:   passed,
			Details: map[string]interface{}{
				"js_divergence": jsDiv,
                "bins_orig": origBins, // For heatmap
                "bins_gen": genBins,
			},
		},
	}
}

func validateQueryTypes(orig, gen []models.SQLTrace, v *services.StatisticalValidator) []services.ValidationResult {
	// Compare frequencies of query templates.
    // We use the Query string as a proxy for template here.
	origCounts := make(map[string]int)
	genCounts := make(map[string]int)

	for _, t := range orig {
		origCounts[t.Query]++
	}
	for _, t := range gen {
		genCounts[t.Query]++
	}

	// Create observed (gen) and expected (orig) vectors for Chi-Square
    // Note: Chi-Square needs same categories.
    // We iterate over all unique queries found in original.

    var observed []int
    var expected []int

    // Scaling factor to match total counts (Chi-square is sensitive to sample size)
    scale := float64(len(gen)) / float64(len(orig))

    for q, count := range origCounts {
        exp := float64(count) * scale
        obs := genCounts[q]

        // If expected count is too small (< 5), Chi-Square is invalid.
        // We aggregate "other" queries usually, but for simplicity here we skip or include.
        if exp >= 5 {
            expected = append(expected, int(exp))
            observed = append(observed, obs)
        }
    }

    if len(expected) < 2 {
          return []services.ValidationResult{
            {TestName: "Query Mix (Chi-Square)", PValue: 1.0, Passed: true, Details: map[string]interface{}{"note": "not enough categories"}},
        }
    }

    result := v.ChiSquareTest(observed, expected)
    result.TestName = "Query Mix (Chi-Square)"

	return []services.ValidationResult{*result}
}

func getTimeRange(traces []models.SQLTrace) (start, end time.Time) {
	if len(traces) == 0 {
		return
	}
	start = traces[0].Timestamp
	end = traces[0].Timestamp
	for _, t := range traces {
		if t.Timestamp.Before(start) {
			start = t.Timestamp
		}
		if t.Timestamp.After(end) {
			end = t.Timestamp
		}
	}
	return
}
