package database

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConnector struct {
	*BaseConnector
	db *sql.DB
}

func NewMySQLConnector(cfg Config) (*MySQLConnector, error) {
	conn := NewBaseConnector(cfg)
	conn.cfg.Driver = "mysql"

	db, err := sql.Open("mysql", conn.connectionString())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return &MySQLConnector{BaseConnector: conn, db: db}, nil
}

func (c *MySQLConnector) Ping(_ context.Context) error {
	return c.db.Ping()
}

func (c *MySQLConnector) Query(_ context.Context, sql string, args ...interface{}) (Rows, error) {
	return c.db.Query(sql, args...)
}

func (c *MySQLConnector) Close() error {
	return c.db.Close()
}
