package validation

import (
	"context"
	"encoding/json"
	"os"

	"github.com/turtacn/SQLTraceBench/internal/app/execution"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
)

type Service interface {
	Validate(ctx context.Context, execOutPath string) error
}

type DefaultService struct {
	execService   execution.Service
	validationSvc services.ValidationService
	log           *utils.Logger
}

func NewService(execService execution.Service, valSvc services.ValidationService) Service {
	return &DefaultService{
		execService:   execService,
		validationSvc: valSvc,
		log:           utils.GetGlobalLogger(),
	}
}

func (s *DefaultService) Validate(ctx context.Context, execOutPath string) error {
	var synth models.PerformanceMetrics
	f, _ := os.Open(execOutPath)
	defer f.Close()
	_ = json.NewDecoder(f).Decode(&synth)

	report, _ := s.validationSvc.Validate(ctx, &models.PerformanceMetrics{}, &synth)
	s.log.Info("validation done", utils.Field{Key: "passed", Value: report.Passed})
	out, _ := os.Create("validation.json")
	defer out.Close()
	return json.NewEncoder(out).Encode(report)
}
