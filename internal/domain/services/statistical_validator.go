package services

import (
	"math"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

// ValidationResult represents the result of a statistical test.
type ValidationResult struct {
	TestName string
	PValue   float64
	Passed   bool
	Details  map[string]interface{}
}

// StatisticalValidator performs statistical tests on datasets.
type StatisticalValidator struct {
	KSThreshold        float64
	ChiSquareThreshold float64
}

// NewStatisticalValidator creates a new StatisticalValidator.
func NewStatisticalValidator(ksThreshold, chiSquareThreshold float64) *StatisticalValidator {
	return &StatisticalValidator{
		KSThreshold:        ksThreshold,
		ChiSquareThreshold: chiSquareThreshold,
	}
}

// KolmogorovSmirnovTest performs the Two-Sample Kolmogorov-Smirnov test.
func (v *StatisticalValidator) KolmogorovSmirnovTest(observed, expected []float64) *ValidationResult {
	// Calculate KS statistic using gonum/stat.
	// We assume unweighted samples (nil weights).
	ksStat := stat.KolmogorovSmirnov(observed, nil, expected, nil)
    n1 := float64(len(observed))
    n2 := float64(len(expected))

    // P-value calculation (Smirnov approximation)
    // Ref: Numerical Recipes in C, section 14.3
    en := math.Sqrt(n1 * n2 / (n1 + n2))
    lambda := (en + 0.12 + 0.11/en) * ksStat
    pValue := kolmogorovQ(lambda)

    passed := pValue >= v.KSThreshold

	return &ValidationResult{
		TestName: "KS Test",
		PValue:   pValue,
		Passed:   passed,
		Details: map[string]interface{}{
			"ks_statistic": ksStat,
			"sample_size_1":  len(observed),
            "sample_size_2":  len(expected),
		},
	}
}

// kolmogorovQ calculates the complementary cumulative distribution function of the Kolmogorov distribution.
func kolmogorovQ(lambda float64) float64 {
    if lambda < 1.1e-16 {
        return 1.0
    }

    term := 0.0
    for j := 1; j <= 100; j++ {
        t := math.Exp(-2.0 * lambda * lambda * float64(j*j))
        if j%2 == 1 {
            term += t
        } else {
            term -= t
        }
        if t <= 0.001*term { // Convergence check
            return 2.0 * term
        }
    }
    return 2.0 * term
}


// ChiSquareTest performs the Chi-Square Goodness of Fit test.
func (v *StatisticalValidator) ChiSquareTest(observedFreq, expectedFreq []int) *ValidationResult {
	chiSquare := 0.0
	for i := range observedFreq {
        if expectedFreq[i] == 0 {
            continue // Avoid division by zero
        }
		diff := float64(observedFreq[i] - expectedFreq[i])
		chiSquare += (diff * diff) / float64(expectedFreq[i])
	}

	df := float64(len(observedFreq) - 1)
    var pValue float64
    if df > 0 {
	    chiDist := distuv.ChiSquared{K: df}
	    pValue = 1 - chiDist.CDF(chiSquare)
    } else {
        pValue = 0.0 // Not enough degrees of freedom
        if chiSquare == 0 {
            pValue = 1.0
        }
    }

	return &ValidationResult{
		TestName: "Chi-Square Test",
		PValue:   pValue,
		Passed:   pValue >= v.ChiSquareThreshold,
		Details: map[string]interface{}{
			"chi_square": chiSquare,
			"df":         df,
		},
	}
}

// JensenShannonDivergence calculates the Jensen-Shannon Divergence between two probability distributions.
func (v *StatisticalValidator) JensenShannonDivergence(p, q []float64) float64 {
	if len(p) != len(q) {
		return math.Inf(1)
	}

	m := make([]float64, len(p))
	for i := range p {
		m[i] = 0.5 * (p[i] + q[i])
	}

	kl1 := kullbackLeibler(p, m)
	kl2 := kullbackLeibler(q, m)

	return 0.5 * (kl1 + kl2)
}

func kullbackLeibler(p, q []float64) float64 {
	kl := 0.0
	for i := range p {
		if p[i] > 0 && q[i] > 0 {
			kl += p[i] * math.Log(p[i]/q[i])
		}
	}
	return kl
}
