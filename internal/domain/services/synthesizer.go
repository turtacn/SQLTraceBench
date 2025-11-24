package services

import (
	"fmt"
	"strings"

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
				sampler = &BoundZipfSampler{
					sampler: zipfSvc,
					model:   model,
				}
			case models.DistUniform:
				sampler = &BoundUniformSampler{
					sampler: weightedSvc,
					model:   model,
				}
			default:
				// Default to Weighted/Empirical
				sampler = &BoundWeightedSampler{
					sampler: weightedSvc,
					model:   model,
				}
			}
			s.samplers[groupKey][paramName] = sampler
		}
	}

	return s
}

// FillParameters generates values for the template's parameters and returns the filled SQL.
// It returns the SQL string with values inlined (for logging) and the list of arguments (for execution).
func (s *Synthesizer) FillParameters(tmpl *models.SQLTemplate) (string, []interface{}, error) {
	groupSamplers, ok := s.samplers[tmpl.GroupKey]

	// We need to generate args in the order of tmpl.Parameters
	args := make([]interface{}, 0, len(tmpl.Parameters))

	// Create a map for value replacement in string
	paramValues := make(map[string]interface{})

	for _, paramName := range tmpl.Parameters {
		var val interface{} = "DEFAULT" // Fallback
		var err error

		if ok && groupSamplers != nil {
			if sampler, exists := groupSamplers[paramName]; exists {
				val, err = sampler.Sample()
				if err != nil {
					val = "ERR_SAMPLE"
				}
			}
		}

		args = append(args, val)
		paramValues[paramName] = val
	}

	filledSQL := tmpl.RawSQL
	for _, pName := range tmpl.Parameters {
		val := paramValues[pName]
		valStr := fmt.Sprintf("%v", val)

		// Simple quoting for strings (naive)
		if _, isString := val.(string); isString {
			valStr = fmt.Sprintf("'%s'", strings.ReplaceAll(valStr, "'", "''"))
		}

		filledSQL = strings.ReplaceAll(filledSQL, pName, valStr)
	}

	return filledSQL, args, nil
}
