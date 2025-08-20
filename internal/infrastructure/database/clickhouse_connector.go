package database

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseConnector struct {
	*BaseConnector
	conn clickhouse.Conn
}

func NewClickHouseConnector(cfg Config) (*ClickHouseConnector, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
	})
	if err != nil {
		return nil, err
	}
	return &ClickHouseConnector{
		BaseConnector: &BaseConnector{cfg: cfg},
		conn:          conn,
	}, nil
}

func (c *ClickHouseConnector) Ping(_ context.Context) error {
	return c.conn.Ping(context.Background())
}

func (c *ClickHouseConnector) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	rows, err := c.conn.Query(ctx, sql, args...)
	return clickhouseRows{rows}, err
}

func (c *ClickHouseConnector) Close() error {
	return c.conn.Close()
}

// simple wrapper
type clickhouseRows struct {
	clickhouse.Rows
}
