package reporters

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

// HTMLReporter generates HTML reports for validation results.
type HTMLReporter struct {
	TemplatePath string
}

// NewHTMLReporter creates a new HTMLReporter.
func NewHTMLReporter(templatePath string) *HTMLReporter {
	return &HTMLReporter{
		TemplatePath: templatePath,
	}
}

// GenerateReport renders the validation report into an HTML file.
func (r *HTMLReporter) GenerateReport(
	report *services.ValidationReport,
	outputPath string,
) error {
	// Parse the template
	tmpl, err := template.ParseFiles(r.TemplatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create the output file
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Render the template
	if err := tmpl.Execute(file, report); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
