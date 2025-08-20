package models

import (
	"context"
	"sync"
	"time"
)

type BenchmarkWorkload struct {
	WorkloadID string
	Queries    []WorkloadQuery
	Config     ExecutionConfig
	Schedule   QuerySchedule
}

type WorkloadQuery struct {
	ID           string
	SQL          string
	Weight       float64
	ExpectedTime float64 // ms
}

type ExecutionConfig struct {
	TargetQPS      float64
	MaxConcurrency int
	Duration       time.Duration
	Warmup         time.Duration
	DatabaseType   string // "starrocks" | "clickhouse"
}

type QuerySchedule struct {
	mu      sync.RWMutex
	queries []WorkloadQuery
	closed  bool
}

func (qs *QuerySchedule) GetNextQuery() (WorkloadQuery, bool) {
	qs.mu.RLock()
	defer qs.mu.RUnlock()
	if qs.closed || len(qs.queries) == 0 {
		return WorkloadQuery{}, false
	}
	return qs.queries[0], true
}

func (qs *QuerySchedule) Close() {
	qs.mu.Lock()
	qs.closed = true
	qs.mu.Unlock()
}

func (w *BenchmarkWorkload) Execute(ctx context.Context) {}

//Personal.AI order the ending
