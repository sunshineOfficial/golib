package db

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
)

const clickhouseDriver = "clickhouse"

func NewClickhouse(ctx context.Context, connectionString string, options ...Option) (*sqlx.DB, error) {
	handler := applyOptions(options...)

	var (
		db  *sqlx.DB
		err error
	)

	if handler.traces {
		db, err = otelsqlx.Open(clickhouseDriver, connectionString)
	} else {
		db, err = sqlx.Open(clickhouseDriver, connectionString)
	}

	if err != nil {
		return nil, fmt.Errorf("can't open connection: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("can't ping connection: %w", err)
	}

	return db, nil
}

func NewNativeClickhouse(ctx context.Context, clickhouseOptions *clickhouse.Options) (driver.Conn, error) {
	conn, err := clickhouse.Open(clickhouseOptions)
	if err != nil {
		return nil, fmt.Errorf("can't open connection: %w", err)
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("can't ping connection: %w", err)
	}

	return conn, nil
}

func NewClickhouseOptions(host string, port int, database, username, password string) *clickhouse.Options {
	return &clickhouse.Options{
		Protocol: clickhouse.Native,
		Addr:     []string{fmt.Sprintf("%s:%d", host, port)},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
	}
}
