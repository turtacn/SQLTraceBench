package services

import (
	"context"
	"hash/fnv"
	"regexp"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/repositories"
)

type TemplateService interface {
	ExtractTemplates(ctx context.Context, traces []models.SQLTrace) ([]models.SQLTemplate, error)
}

type DefaultTemplateService struct {
	repo repositories.TemplateRepository
}

func NewTemplateService(repo repositories.TemplateRepository) TemplateService {
	return &DefaultTemplateService{repo: repo}
}

func (s *DefaultTemplateService) ExtractTemplates(_ context.Context, traces []models.SQLTrace) ([]models.SQLTemplate, error) {
	templateMap := make(map[string]*models.SQLTemplate)
	placeholderRegex := regexp.MustCompile(`'[^']*'|\d+`)

	for _, t := range traces {
		tempSQL := placeholderRegex.ReplaceAllString(t.Query, "{{param}}")
		key := hashString(tempSQL)
		if tpl, ok := templateMap[key]; ok {
			tpl.Frequency++
		} else {
			templateMap[key] = &models.SQLTemplate{
				TemplateID:    key,
				OriginalQuery: t.Query,
				TemplateQuery: tempSQL,
				Frequency:     1,
				CreatedAt:     time.Now(),
			}
		}
	}
	var out []models.SQLTemplate
	for _, v := range templateMap {
		out = append(out, *v)
	}
	return out, nil
}

func hashString(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return string(h.Sum(nil))
}
