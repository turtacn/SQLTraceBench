package generation

import (
	"context"
	"encoding/json"
	"os"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

type Service interface {
	GenerateWorkload(ctx context.Context, yamlTplPath, yamlOutPath string) error
}

type DefaultService struct {
	workloadSvc services.WorkloadService
}

func NewService(ws services.WorkloadService) Service { return &DefaultService{workloadSvc: ws} }

func (s *DefaultService) GenerateWorkload(ctx context.Context, yamlTplPath, yamlOutPath string) error {
	var tpls []models.SQLTemplate
	file, _ := os.Open(yamlTplPath)
	defer file.Close()
	_ = json.NewDecoder(file).Decode(map[string][]models.SQLTemplate{"templates": tpls})

	pm := &models.ParameterModel{} // dummy
	wl := s.workloadSvc.GenerateWorkload(ctx, tpls, pm)

	f, _ := os.Create(yamlOutPath)
	defer f.Close()
	return json.NewEncoder(f).Encode(map[string]interface{}{"query_count": len(wl.Queries)})
}
