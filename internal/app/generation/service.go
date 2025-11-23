package generation

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
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

	pm := s.paramSvc.BuildModel(models.TraceCollection{}, templates)

	wl, err := s.workloadSvc.GenerateWorkload(templates, pm, n)
	if err != nil {
		return nil, err
	}

	return wl, nil
}

type valFreq struct {
	val  interface{}
	freq int
}

// Generate implements the enhanced generation logic for Phase 2.
func (s *DefaultService) Generate(ctx context.Context, req GenerateRequest) ([]models.SQLTrace, error) {
	if len(req.SourceTraces) == 0 {
		return nil, fmt.Errorf("no source traces provided for generation")
	}

	// 1. Analyze parameters to get distributions
	analyzer := &services.ParameterAnalyzer{MaxCardinality: 10000} // Set limit
	paramStats := analyzer.Analyze(req.SourceTraces)

	// If Hotspots not provided, detect them
	if req.HotspotValues == nil {
		detector := &services.HotspotDetector{Threshold: 0.05}
		req.HotspotValues = make(map[string][]interface{})
		for paramName, stats := range paramStats {
			req.HotspotValues[paramName] = detector.Detect(stats)
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

	paramSamplers := make(map[string]*samplers.ZipfSampler)
	importServicesZipf := func(name string) *samplers.ZipfSampler {
		if s, ok := paramSamplers[name]; ok {
			return s
		}
		hotspots := req.HotspotValues[name]
		s := samplers.NewZipfSampler(1.1)
		s.HotspotValues = hotspots
		// We use default HotspotProb (0.3) for now, as currently there is no config input in GenerateRequest.
		// In a future update, GenerateRequest should support config overrides.
		paramSamplers[name] = s
		return s
	}

	// Pre-convert stats to ValueDistribution for sampling
	// IMPORTANT: Sort values by frequency descending so Zipf makes sense.
	paramDists := make(map[string]*models.ValueDistribution)
	for name, stats := range paramStats {
		vfList := make([]valFreq, 0, len(stats.ValueCounts))
		total := 0
		for val, count := range stats.ValueCounts {
			vfList = append(vfList, valFreq{val: val, freq: count})
			total += count
		}

		// Sort descending
		sort.Slice(vfList, func(i, j int) bool {
			return vfList[i].freq > vfList[j].freq
		})

		dist := models.NewValueDistribution()
		dist.Values = make([]interface{}, len(vfList))
		dist.Frequencies = make([]int, len(vfList))
		dist.Total = total

		for i, vf := range vfList {
			dist.Values[i] = vf.val
			dist.Frequencies[i] = vf.freq
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
			// Fallback timestamp logic?
			newTrace.Timestamp = time.Now()
		}

		generated = append(generated, newTrace)
	}

	return generated, nil
}
