package services

import (
	"fmt"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// ModelSampler defines the interface for a sampler bound to a specific parameter model.
type ModelSampler interface {
	Sample() (interface{}, error)
}

// Synthesizer is responsible for filling SQL templates with synthetic data
// based on statistical parameter models.
type Synthesizer struct {
	// samplers maps GroupKey -> ParamName -> ModelSampler
	samplers map[string]map[string]ModelSampler
}

// BoundZipfSampler adapts the generic ZipfSampler to the ModelSampler interface
// by binding it to a specific ParameterModel.
type BoundZipfSampler struct {
	sampler *ZipfSampler
	model   *models.ParameterModel
}

func (s *BoundZipfSampler) Sample() (interface{}, error) {
	return s.sampler.Sample(s.model)
}

// BoundWeightedSampler adapts the generic WeightedRandomSampler to the ModelSampler interface.
type BoundWeightedSampler struct {
	sampler *WeightedRandomSampler
	model   *models.ParameterModel
}

func (s *BoundWeightedSampler) Sample() (interface{}, error) {
	return s.sampler.Sample(s.model)
}

// BoundUniformSampler adapts a uniform sampling strategy.
type BoundUniformSampler struct {
	sampler *WeightedRandomSampler
	model   *models.ParameterModel
}

func (s *BoundUniformSampler) Sample() (interface{}, error) {
	return s.sampler.Sample(s.model)
}

// NewSynthesizer creates a new Synthesizer initialized with the provided workload parameter models.
func NewSynthesizer(workloadModel *models.WorkloadParameterModel) *Synthesizer {
	s := &Synthesizer{
		samplers: make(map[string]map[string]ModelSampler),
	}

	zipfSvc := NewZipfSampler(1.001) // Default s, will be overridden by model.ZipfS
	weightedSvc := NewWeightedRandomSampler()

	for groupKey, params := range workloadModel.TemplateParameters {
		s.samplers[groupKey] = make(map[string]ModelSampler)
		for paramName, model := range params {
			var sampler ModelSampler

			switch model.DistType {
			case models.DistZipfian:
				sampler = &BoundZipfSampler{sampler: zipfSvc, model: model}
			case models.DistUniform:
				sampler = &BoundUniformSampler{sampler: weightedSvc, model: model}
			default: // Fallback to empirical/weighted
				sampler = &BoundWeightedSampler{sampler: weightedSvc, model: model}
			}
			s.samplers[groupKey][paramName] = sampler
		}
	}

	return s
}

// FillParameters generates values for the template's parameters and returns the list of arguments.
func (s *Synthesizer) FillParameters(tmpl *models.SQLTemplate) ([]interface{}, error) {
	groupSamplers, ok := s.samplers[tmpl.GroupKey]
	if !ok {
		// No specific model for this template group, return default values or error
		return make([]interface{}, len(tmpl.Parameters)),
			fmt.Errorf("no parameter model found for template group key: %s", tmpl.GroupKey)
	}

	args := make([]interface{}, len(tmpl.Parameters))
	for i, paramName := range tmpl.Parameters {
		var val interface{} = "DEFAULT" // Fallback
		var err error

		if sampler, exists := groupSamplers[paramName]; exists {
			val, err = sampler.Sample()
			if err != nil {
				// On sampling error, you might want to use a fallback or log the error
				val = fmt.Sprintf("ERR_SAMPLING_%s", paramName)
			}
		}
		args[i] = val
	}

	return args, nil
}
