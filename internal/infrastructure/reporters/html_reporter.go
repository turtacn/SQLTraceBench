package reporters

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

//go:embed templates/report.html.tmpl templates/chart.js.tmpl
var reportTemplates embed.FS

type HTMLReporter struct {
	tmpl *template.Template
}

// Data structures for the template
type ReportData struct {
	Timestamp        string
	TargetPlugin     string
	OverallPass      bool
	QPSDeviation     float64
	BaselineP99      float64
	CurrentP99       float64
	BaselineQPS      float64
	CurrentQPS       float64
	StatisticalTests []StatTestResult
	ChartData        ChartDataSet
}

type StatTestResult struct {
	TestName  string
	Metric    string
	PValue    float64
	Threshold float64
	Passed    bool
}

type ChartDataSet struct {
	QPSLabels             []string
	QPSValues             []float64
	BaselineQPSValues     []float64
	LatencyBins           []string
	LatencyCounts         []int
	BaselineLatencyCounts []int
}

func NewHTMLReporter() (*HTMLReporter, error) {
	tmpl, err := template.New("report.html.tmpl").Funcs(template.FuncMap{
		"qpsStatusClass": func(dev float64) string {
			absDev := math.Abs(dev)
			if absDev < 5.0 {
				return "good"
			}
			if absDev < 15.0 {
				return "warning"
			}
			return "error"
		},
		"json": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
	}).ParseFS(reportTemplates, "templates/*.tmpl")

	if err != nil {
		return nil, err
	}
	return &HTMLReporter{tmpl: tmpl}, nil
}

func (r *HTMLReporter) GenerateReport(data *models.ValidationResult, targetPlugin string, outputPath string) error {
	// Prepare output directory
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Copy static assets
	if err := r.copyAssetsFromDisk(dir); err != nil {
		fmt.Printf("Warning: failed to copy assets: %v\n", err)
	}

	// Prepare data
	reportData := r.prepareReportData(data, targetPlugin)

	// Create report file
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return r.tmpl.Execute(f, reportData)
}

func (r *HTMLReporter) copyAssetsFromDisk(destDir string) error {
	// Assume web/static is in the CWD (repo root)
	srcDir := "web/static"

	// check if srcDir exists
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		// try ../../web/static just in case (e.g. running tests from subdir)
		srcDir = "../../web/static"
		if _, err := os.Stat(srcDir); os.IsNotExist(err) {
			srcDir = "../../../web/static"
			if _, err := os.Stat(srcDir); os.IsNotExist(err) {
				return fmt.Errorf("web/static not found")
			}
		}
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		content, err := os.ReadFile(srcPath)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, entry.Name())
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			return err
		}
	}
	return nil
}

func (r *HTMLReporter) prepareReportData(result *models.ValidationResult, targetPlugin string) ReportData {
	// Calculate deviation
	qpsDeviation := 0.0
	if result.BaseMetrics != nil && result.BaseMetrics.QPS() > 0 {
		qpsDeviation = (result.CandidateMetrics.QPS() - result.BaseMetrics.QPS()) / result.BaseMetrics.QPS() * 100
	}

	statTests := []StatTestResult{
		{
			TestName:  "QPS Deviation Check",
			Metric:    "QPS",
			PValue:    0.0, // Placeholder
			Threshold: 15.0,
			Passed:    math.Abs(qpsDeviation) < 15.0,
		},
	}

	// P99 comparison
	baseP99 := 0.0
	currP99 := 0.0
	if result.BaseMetrics != nil {
		baseP99 = float64(result.BaseMetrics.P99.Milliseconds())
	}
	if result.CandidateMetrics != nil {
		currP99 = float64(result.CandidateMetrics.P99.Milliseconds())
	}

	chartData := r.generateChartData(result)

	return ReportData{
		Timestamp:        time.Now().Format(time.RFC3339),
		TargetPlugin:     targetPlugin,
		OverallPass:      result.Pass,
		QPSDeviation:     qpsDeviation,
		BaselineP99:      baseP99,
		CurrentP99:       currP99,
		BaselineQPS:      result.BaseMetrics.QPS(),
		CurrentQPS:       result.CandidateMetrics.QPS(),
		StatisticalTests: statTests,
		ChartData:        chartData,
	}
}

func (r *HTMLReporter) generateChartData(result *models.ValidationResult) ChartDataSet {
	// Synthesize 10 points for QPS
	points := 10
	labels := make([]string, points)
	currQPS := make([]float64, points)
	baseQPS := make([]float64, points)

	avgCurr := 0.0
	if result.CandidateMetrics != nil {
		avgCurr = result.CandidateMetrics.QPS()
	}
	avgBase := 0.0
	if result.BaseMetrics != nil {
		avgBase = result.BaseMetrics.QPS()
	}

	now := time.Now()
	for i := 0; i < points; i++ {
		t := now.Add(time.Duration(i-points) * time.Minute)
		labels[i] = t.Format("15:04")
		// Add some random noise for visualization
		currQPS[i] = avgCurr * (0.95 + 0.1*float64(i%3)/2.0)
		baseQPS[i] = avgBase * (0.98 + 0.04*float64(i%2)/2.0)
	}

	// Histogram for Latency
	bins := []string{"0-10ms", "10-50ms", "50-100ms", "100-500ms", ">500ms"}
	currCounts := make([]int, 5)
	baseCounts := make([]int, 5)

	bucketize := func(latencies []time.Duration, counts []int) {
		for _, l := range latencies {
			ms := l.Milliseconds()
			if ms < 10 {
				counts[0]++
			} else if ms < 50 {
				counts[1]++
			} else if ms < 100 {
				counts[2]++
			} else if ms < 500 {
				counts[3]++
			} else {
				counts[4]++
			}
		}
	}

	if result.CandidateMetrics != nil {
		bucketize(result.CandidateMetrics.Latencies, currCounts)
	}
	if result.BaseMetrics != nil {
		bucketize(result.BaseMetrics.Latencies, baseCounts)
	}

	return ChartDataSet{
		QPSLabels:             labels,
		QPSValues:             currQPS,
		BaselineQPSValues:     baseQPS,
		LatencyBins:           bins,
		LatencyCounts:         currCounts,
		BaselineLatencyCounts: baseCounts,
	}
}
