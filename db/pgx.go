package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// NewPgx creates sqlx.DB instance with pgx driver
func NewPgx(ctx context.Context, connectionString string, options ...Option) (*sqlx.DB, error) {
	pool, err := NewPgxPool(ctx, connectionString, options...)
	if err != nil {
		return nil, fmt.Errorf("can't create pgx pool: %w", err)
	}

	return WrapPgxPool(ctx, pool)
}

func NewPgxPool(ctx context.Context, connectionString string, options ...Option) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("can't parse connection string: %w", err)
	}
	if poolConfig == nil {
		return nil, errors.New("parsed config is nil")
	}

	handler := applyOptions(options...)
	if handler.traces {
		poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	}

	return pgxpool.NewWithConfig(ctx, poolConfig)
}

func WrapPgxPool(ctx context.Context, pool *pgxpool.Pool) (*sqlx.DB, error) {
	db := sqlx.NewDb(stdlib.OpenDBFromPool(pool), "pgx")
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("can't ping server: %w", err)
	}

	return db, nil
}
