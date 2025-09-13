package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseConnector struct {
	*BaseConnector
	db *sql.DB
}

func NewClickHouseConnector(cfg Config) (*ClickHouseConnector, error) {
	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
	})

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &ClickHouseConnector{
		BaseConnector: &BaseConnector{cfg: cfg},
		db:            db,
	}, nil
}

func (c *ClickHouseConnector) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c *ClickHouseConnector) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *ClickHouseConnector) Close() error {
	return c.db.Close()
}
