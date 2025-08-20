package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/turtacn/SQLTraceBench/pkg/types"
)

type SQLTrace struct {
	Timestamp       time.Time `json:"timestamp"`
	Query           string    `json:"query"`
	ExecutionTimeMs float64   `json:"execution_time_ms"`
	RowsReturned    int64     `json:"rows_returned"`
	RowsScanned     int64     `json:"rows_scanned,omitempty"`
	DatabaseName    string    `json:"database_name,omitempty"`
	UserName        string    `json:"user_name,omitempty"`
	ClientIP        string    `json:"client_ip,omitempty"`
	QueryHash       string    `json:"query_hash,omitempty"`
}

func (t *SQLTrace) Validate() *types.SQLTraceBenchError {
	if t.Query == "" {
		return types.NewError(types.ErrInvalidInput, "query field cannot be empty")
	}
	if t.Timestamp.IsZero() {
		return types.NewError(types.ErrInvalidInput, "timestamp field cannot be zero")
	}
	return nil
}

func (t *SQLTrace) GetQueryHash() string {
	if t.QueryHash == "" {
		hash := sha256.Sum256([]byte(t.Query))
		t.QueryHash = hex.EncodeToString(hash[:16])
	}
	return t.QueryHash
}

func (t *SQLTrace) String() string {
	return fmt.Sprintf("[%s] %s (%.2fms, %d rows)", t.Timestamp.Format("15:04:05"), t.Query, t.ExecutionTimeMs, t.RowsReturned)
}

type TraceCollection struct {
	Traces []SQLTrace
}

func (tc *TraceCollection) Add(trace SQLTrace) {
	tc.Traces = append(tc.Traces, trace)
}

func (tc *TraceCollection) Count() int {
	return len(tc.Traces)
}

//Personal.AI order the ending
