package models

// DistributionType defines the type of statistical distribution.
type DistributionType string

const (
	DistUniform DistributionType = "uniform"
	DistZipfian DistributionType = "zipfian"
	DistNormal  DistributionType = "normal"
	DistEmpirical DistributionType = "empirical" // Fallback for arbitrary distribution
)

// ParameterModel holds the detailed statistical model for a single parameter.
type ParameterModel struct {
	ParamName      string           `json:"param_name"`
	DataType       string           `json:"data_type"` // e.g., "INT", "VARCHAR"
	DistType       DistributionType `json:"dist_type"` // Detected distribution

	// Zipf specific
	ZipfS          float64 `json:"zipf_s,omitempty"` // Skewness parameter
	ZipfV          float64 `json:"zipf_v,omitempty"` // Value range parameter (usually 1)

	// Stats
	Cardinality    int     `json:"cardinality"`
	HotspotRatio   float64 `json:"hotspot_ratio"` // e.g., 0.8 means 20% items get 80% traffic

	// Sample Values (for reproduction)
	// We store TopValues and TopFrequencies to support empirical sampling of hotspots.
	TopValues      []interface{} `json:"top_values"`
	TopFrequencies []int         `json:"top_frequencies"`
}

// WorkloadParameterModel holds the statistical model of parameters for a set of SQL templates.
// It maps each template (by its GroupKey) to a model of its parameters.
type WorkloadParameterModel struct {
	// TemplateParameters maps a template's GroupKey to a map of its parameters.
	// Each parameter is then mapped to its parameter model.
	TemplateParameters map[string]map[string]*ParameterModel
}

// NewWorkloadParameterModel creates an empty workload parameter model.
func NewWorkloadParameterModel() *WorkloadParameterModel {
	return &WorkloadParameterModel{
		TemplateParameters: make(map[string]map[string]*ParameterModel),
	}
}

// ValueDistribution holds the observed values and their frequencies for a single parameter.
// Deprecated: This struct is kept for temporary compatibility but should be replaced by ParameterModel.
// Or we use this to accumulate data before converting to ParameterModel.
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
