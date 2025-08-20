package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/turtacn/SQLTraceBench/pkg/utils"
)

type Config struct {
	Driver                     string // mysql / clickhouse / starrocks
	Host                       string
	Port                       int
	User, Password, Database   string
	MaxOpenConns, MaxIdleConns int
	ConnTimeout                time.Duration
}

type Connector interface {
	Ping(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)
	Close() error
}

type Rows interface {
	Close() error
	Next() bool
	Scan(dest ...interface{}) error
}

type BaseConnector struct {
	log *utils.Logger
	cfg Config
}

func NewBaseConnector(cfg Config) *BaseConnector {
	return &BaseConnector{cfg: cfg, log: utils.GetGlobalLogger()}
}

func (c *BaseConnector) Ping(ctx context.Context) error {
	return errors.New("must implement by driver")
}

func (c *BaseConnector) Close() error {
	return nil
}

func (c *BaseConnector) connectionString() string {
	switch c.cfg.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s",
			c.cfg.User, c.cfg.Password, c.cfg.Host, c.cfg.Port, c.cfg.Database,
			c.cfg.ConnTimeout)
	default:
		return ""
	}
}
