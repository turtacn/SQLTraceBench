package services

import (
	"context"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/repositories"
)

type ParameterService interface {
	BuildModel(ctx context.Context, traces []models.SQLTrace) (*models.ParameterModel, error)
	GenerateParams(model *models.ParameterModel, tpl *models.SQLTemplate) map[string]interface{}
}

type DefaultParameterService struct {
	repo repositories.ParameterRepository
}

func NewParameterService(repo repositories.ParameterRepository) ParameterService {
	return &DefaultParameterService{repo: repo}
}

func (s *DefaultParameterService) BuildModel(_ context.Context, traces []models.SQLTrace) (*models.ParameterModel, error) {
	pm := &models.ParameterModel{
		Parameters: make(map[string]models.ParamStats),
	}
	for _, t := range traces {
		// naive integer scan (MVP only)
		re := regexp.MustCompile(`(?:=|in)\s*(\d+)`)
		matches := re.FindAllStringSubmatch(strings.ToLower(t.Query), -1)
		for _, m := range matches {
			key := "param"
			stat := models.ParamStats{DistributionType: models.TypeInteger}
			if _, ok := pm.Parameters[key]; !ok {
				pm.Parameters[key] = stat
			}
		}
	}
	pm.GeneratedAt = time.Now()
	return pm, nil
}

func (s *DefaultParameterService) GenerateParams(_ *models.ParameterModel, _ *models.SQLTemplate) map[string]interface{} {
	return map[string]interface{}{"param": rand.Intn(9999)}
}
