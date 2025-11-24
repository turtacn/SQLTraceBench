package services

import (
	"context"
	"math/rand"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// GenerationConfig holds configuration for workload generation.
type GenerationConfig struct {
	TotalQueries int
	ScaleFactor  float64
}

// WorkloadService is responsible for generating a benchmark workload.
type WorkloadService struct {
	sampler Sampler // Kept for backward compatibility if needed, but primary logic moves to Synthesizer
}

// NewWorkloadService creates a new WorkloadService.
func NewWorkloadService(sampler Sampler) *WorkloadService {
	return &WorkloadService{sampler: sampler}
}

// GenerateWorkload creates a BenchmarkWorkload from SQL templates and a parameter model.
// It uses the Synthesizer for "Smart Generation".
func (s *WorkloadService) GenerateWorkload(
	ctx context.Context,
	templates []models.SQLTemplate,
	pm *models.WorkloadParameterModel,
	config GenerationConfig,
) (*models.BenchmarkWorkload, error) {
	var wl models.BenchmarkWorkload

	// 1. Initialize Synthesizer
	synthesizer := NewSynthesizer(pm)

	// 2. Prepare Template Selector (Weighted Random)
	var totalWeight int
	var weightedTemplates []models.SQLTemplate

	for _, t := range templates {
		w := t.Weight
		if w <= 0 {
			w = 1 // Default weight
		}
		t.Weight = w // update locally
		totalWeight += w
		weightedTemplates = append(weightedTemplates, t)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 3. Generation Loop
	targetCount := config.TotalQueries
	if targetCount <= 0 {
		return &wl, nil
	}

	wl.Queries = make([]models.QueryWithArgs, 0, targetCount)

	for i := 0; i < targetCount; i++ {
		// Select Template
		var selectedTmpl *models.SQLTemplate
		r := rng.Intn(totalWeight)
		currentW := 0
		for idx := range weightedTemplates {
			currentW += weightedTemplates[idx].Weight
			if r < currentW {
				selectedTmpl = &weightedTemplates[idx]
				break
			}
		}
		if selectedTmpl == nil {
			selectedTmpl = &weightedTemplates[len(weightedTemplates)-1]
		}

		// Synthesize Parameters
		filledSQL, _, err := synthesizer.FillParameters(selectedTmpl)
		if err != nil {
			continue
		}

		wl.Queries = append(wl.Queries, models.QueryWithArgs{
			Query: filledSQL,
			Args:  nil, // Use inlined values
		})
	}

	return &wl, nil
}
