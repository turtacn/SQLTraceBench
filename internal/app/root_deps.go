package app

import (
	"github.com/turtacn/SQLTraceBench/internal/app/conversion"
	"github.com/turtacn/SQLTraceBench/internal/app/execution"
	"github.com/turtacn/SQLTraceBench/internal/app/validation"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/database"
)

type Root struct {
	Config database.Config

	// domain services
	TemplateSvc   services.TemplateService
	ParameterSvc  services.ParameterService
	WorkloadSvc   services.WorkloadService
	SchemaSvc     services.SchemaService
	ValidationSvc services.ValidationService

	// app services
	Conversion conversion.Service
	Execution  execution.Service
	Validation validation.Service
}

func NewRoot(cfg database.Config) *Root {
	// simple instantiation
	svc := &Root{
		Config: cfg,
	}
	svc.SchemaSvc = services.NewSchemaService()
	svc.TemplateSvc = services.NewTemplateService(nil) // repo nil for MVP
	svc.ParameterSvc = services.NewParameterService(nil)
	svc.WorkloadSvc = services.NewWorkloadService(svc.ParameterSvc)
	svc.ValidationSvc = services.NewValidationService()

	// app
	svc.Conversion = conversion.NewService(svc.TemplateSvc, svc.SchemaSvc)
	svc.Execution = execution.NewService()
	svc.Validation = validation.NewService(svc.Execution, svc.ValidationSvc)

	return svc
}
