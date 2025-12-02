package starrocks

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
)

type BenchmarkExecutor struct {
	conn *sql.DB
}

func NewBenchmarkExecutor(conn *sql.DB) *BenchmarkExecutor {
	return &BenchmarkExecutor{conn: conn}
}

func (e *BenchmarkExecutor) ExecuteQuery(ctx context.Context, req *proto.ExecuteQueryRequest) (*proto.ExecuteQueryResponse, error) {
	start := time.Now()
	_, err := e.conn.ExecContext(ctx, req.Sql, req.Args)
	duration := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &proto.ExecuteQueryResponse{
		DurationMicros: duration.Microseconds(),
	}, nil
}
