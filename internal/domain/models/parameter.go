package models

import (
	"math/rand"
	"time"

	"github.com/turtacn/SQLTraceBench/pkg/types"
)

type ParameterModel struct {
	Templates   map[string]types.ParameterType `json:"templates"`
	Parameters  map[string]ParamStats          `json:"parameters"`
	GlobalStats GlobalStatistics               `json:"global_stats"`
	GeneratedAt time.Time                      `json:"generated_at"`
}

type ParamStats struct {
	Type           types.ParameterType
	Distribution   types.DistributionType
	ValueSet       []interface{}
	Frequencies    []int
	MinVal, MaxVal float64
}

type Distribution interface {
	Sample() interface{}
	GetMean() float64
}

type ZipfianDistribution struct {
	Alpha float64
	N     int
	r     *rand.Rand
}

func (z *ZipfianDistribution) Sample() interface{} {
	// Very basic zipfian
	return 1 + z.r.Intn(z.N)
}

func (z *ZipfianDistribution) GetMean() float64 {
	return float64(z.N) / 2
}

type GlobalStatistics struct {
	TotalQueries      int64
	UniqueParameters  int
	ParameterCoverage float64
	TimeRange         struct {
		Start, End time.Time
	}
}

func (pm *ParameterModel) UpdateStats() {
	// Incrementally update
	pm.GlobalStats.DataTimestamp = time.Now()
}

//Personal.AI order the ending
