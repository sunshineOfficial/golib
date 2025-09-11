package db

import (
	"context"
	"fmt"

	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"

	_ "github.com/ClickHouse/clickhouse-go/v2" // clickhouse driver
	"github.com/jmoiron/sqlx"
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
		return nil, err
	}

	return db, nil
}
