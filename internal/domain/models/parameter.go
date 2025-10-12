package models

// ParameterModel holds the statistical model of parameters for a set of SQL templates.
// It maps each template (by its GroupKey) to a model of its parameters.
type ParameterModel struct {
	// TemplateParameters maps a template's GroupKey to a map of its parameters.
	// Each parameter is then mapped to its value distribution.
	TemplateParameters map[string]map[string]*ValueDistribution
}

// NewParameterModel creates an empty parameter model.
func NewParameterModel() *ParameterModel {
	return &ParameterModel{
		TemplateParameters: make(map[string]map[string]*ValueDistribution),
	}
}

// ValueDistribution holds the observed values and their frequencies for a single parameter.
type ValueDistribution struct {
	// Values is a slice of unique parameter values.
	Values []interface{}
	// Frequencies is a slice of corresponding frequencies for each value.
	Frequencies []int
	// Total is the total number of observations for this parameter.
	Total int
}

// NewValueDistribution creates an empty value distribution.
func NewValueDistribution() *ValueDistribution {
	return &ValueDistribution{
		Values:      make([]interface{}, 0),
		Frequencies: make([]int, 0),
	}
}

// AddObservation records an observation of a parameter value.
// If the value has been seen before, its frequency is incremented.
// Otherwise, the new value is added with a frequency of 1.
func (vd *ValueDistribution) AddObservation(value interface{}) {
	vd.Total++
	for i, v := range vd.Values {
		if v == value {
			vd.Frequencies[i]++
			return
		}
	}
	vd.Values = append(vd.Values, value)
	vd.Frequencies = append(vd.Frequencies, 1)
}