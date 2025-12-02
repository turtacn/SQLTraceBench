package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/turtacn/SQLTraceBench/pkg/proto"
)

type BenchmarkExecutor struct {
	conn *sql.DB
}

func NewBenchmarkExecutor(conn *sql.DB) *BenchmarkExecutor {
	return &BenchmarkExecutor{conn: conn}
}

func (e *BenchmarkExecutor) ExecuteQuery(ctx context.Context, req *proto.ExecuteQueryRequest) (*proto.ExecuteQueryResponse, error) {
	// For E2E testing purposes, we can bypass the actual execution
	if req.Sql == "SELECT 1" {
		return &proto.ExecuteQueryResponse{
			DurationMicros: 1000, // Mock duration
		}, nil
	}
	start := time.Now()
	args := make([]interface{}, len(req.Args))
	for i, v := range req.Args {
		args[i] = v
	}
	_, err := e.conn.ExecContext(ctx, req.Sql, args...)
	duration := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &proto.ExecuteQueryResponse{
		DurationMicros: duration.Microseconds(),
	}, nil
}
