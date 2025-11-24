package generation

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/samplers"
)

// GenerateRequest encapsulates parameters for workload generation.
type GenerateRequest struct {
	SourceTraces    []models.SQLTrace
	Count           int
	HotspotValues   map[string][]interface{}
	TemporalPattern *services.TemporalPattern
}

// Service is the interface for the workload generation service.
type Service interface {
	GenerateWorkload(ctx context.Context, tplPath string, n int) (*models.BenchmarkWorkload, error)
	// Generate performs advanced generation using traces and parameter modeling.
	Generate(ctx context.Context, req GenerateRequest) ([]models.SQLTrace, error)
	// SetSampler allows configuring the sampler used for generation.
	SetSampler(sampler services.Sampler)
}

// DefaultService is the default implementation of the workload generation service.
type DefaultService struct {
	workloadSvc *services.WorkloadService
	paramSvc    *services.ParameterService
}

// NewService creates a new DefaultService.
func NewService() Service {
	sampler := services.NewWeightedRandomSampler()
	return &DefaultService{
		workloadSvc: services.NewWorkloadService(sampler),
		paramSvc:    services.NewParameterService(),
	}
}

// SetSampler updates the sampler used by the workload service.
func (s *DefaultService) SetSampler(sampler services.Sampler) {
	s.workloadSvc = services.NewWorkloadService(sampler)
}

// GenerateWorkload creates a benchmark workload from a set of SQL templates.
func (s *DefaultService) GenerateWorkload(ctx context.Context, tplPath string, n int) (*models.BenchmarkWorkload, error) {
	file, err := os.Open(tplPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var templates []models.SQLTemplate
	if err := json.NewDecoder(file).Decode(&templates); err != nil {
		return nil, err
	}

	// NOTE: P04 Implementation Assumption
	// Ideally, we should load a pre-computed WorkloadParameterModel here (e.g., from a JSON file).
	// Currently, the interface only accepts the template path.
	// We build a model from scratch, which implies using default distributions if no trace data is supplied.
	// In a full integration, `BuildModel` should be called with the original traces that generated these templates.

	pm := s.paramSvc.BuildModel(models.TraceCollection{}, templates)

	config := services.GenerationConfig{
		TotalQueries: n,
		ScaleFactor:  1.0,
	}

	wl, err := s.workloadSvc.GenerateWorkload(ctx, templates, pm, config)
	if err != nil {
		return nil, err
	}

	return wl, nil
}

// Generate implements the enhanced generation logic.
// This is kept for compatibility with existing tests/tools but should eventually
// use the Synthesizer approach if possible.
func (s *DefaultService) Generate(ctx context.Context, req GenerateRequest) ([]models.SQLTrace, error) {
	if len(req.SourceTraces) == 0 {
		return nil, fmt.Errorf("no source traces provided for generation")
	}

	// 1. Analyze parameters to get distributions
	analyzer := services.NewParameterAnalyzer()
	paramStats := analyzer.Analyze(req.SourceTraces)

	// If Hotspots not provided, extract them from stats
	if req.HotspotValues == nil {
		req.HotspotValues = make(map[string][]interface{})
		for paramName, model := range paramStats {
			req.HotspotValues[paramName] = model.TopValues
		}
	}

	// If Patterns not provided, extract them
	if req.TemporalPattern == nil {
		extractor := &services.TemporalPatternExtractor{Window: time.Hour}
		req.TemporalPattern = extractor.Extract(req.SourceTraces)
	}

	// 2. Prepare Samplers
	var temporalSampler *samplers.TemporalSampler
	if req.TemporalPattern != nil {
		minTime, _ := services.FindTimeRange(req.SourceTraces)
		temporalSampler = samplers.NewTemporalSampler(req.TemporalPattern, minTime)
	}

	// Helper to get sampler
	paramSamplers := make(map[string]*samplers.ZipfSampler)
	importServicesZipf := func(name string) *samplers.ZipfSampler {
		if s, ok := paramSamplers[name]; ok {
			return s
		}
		hotspots := req.HotspotValues[name]
		s := samplers.NewZipfSampler(1.1)
		s.HotspotValues = hotspots
		paramSamplers[name] = s
		return s
	}

	paramDists := make(map[string]*models.ValueDistribution)
	for name, model := range paramStats {
		dist := models.NewValueDistribution()
		// Reconstruct from TopValues/TopFrequencies
		if len(model.TopFrequencies) > 0 {
			dist.Values = make([]interface{}, len(model.TopValues))
			dist.Frequencies = make([]int, len(model.TopFrequencies))
			copy(dist.Values, model.TopValues)
			copy(dist.Frequencies, model.TopFrequencies)
			for _, f := range model.TopFrequencies {
				dist.Total += f
			}
		} else {
			dist.Values = model.TopValues
			dist.Frequencies = make([]int, len(model.TopValues))
			for i := range dist.Values {
				dist.Frequencies[i] = len(dist.Values) - i
				dist.Total += dist.Frequencies[i]
			}
		}
		paramDists[name] = dist
	}

	generated := make([]models.SQLTrace, 0, req.Count)
	sourceCount := len(req.SourceTraces)

	for i := 0; i < req.Count; i++ {
		srcTrace := req.SourceTraces[i%sourceCount]

		newTrace := models.SQLTrace{
			Query:      srcTrace.Query,
			Parameters: make(map[string]interface{}),
		}

		for pName := range srcTrace.Parameters {
			sampler := importServicesZipf(pName)
			dist, ok := paramDists[pName]
			if !ok || dist.Total == 0 {
				continue
			}

			val, err := sampler.Sample(dist)
			if err == nil {
				newTrace.Parameters[pName] = val
			}
		}

		if temporalSampler != nil {
			newTrace.Timestamp = temporalSampler.SampleTimestamp()
		} else {
			newTrace.Timestamp = time.Now()
		}

		generated = append(generated, newTrace)
	}

	return generated, nil
}
