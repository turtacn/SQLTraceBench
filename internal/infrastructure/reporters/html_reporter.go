package reporters

import (
	"embed"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

//go:embed templates/report.html
var reportTemplate embed.FS

type HTMLReporter struct {
	template *template.Template
}

func NewHTMLReporter() (*HTMLReporter, error) {
	t, err := template.New("report.html").Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	}).ParseFS(reportTemplate, "templates/report.html")
	if err != nil {
		return nil, err
	}
	return &HTMLReporter{template: t}, nil
}

func (r *HTMLReporter) GenerateReport(data *models.ValidationReport, outputPath string) error {
	// If the output path is an existing directory, create the report file inside it.
	info, err := os.Stat(outputPath)
	if err == nil && info.IsDir() {
		outputPath = filepath.Join(outputPath, "validation_report.html")
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return r.template.Execute(f, data)
}
