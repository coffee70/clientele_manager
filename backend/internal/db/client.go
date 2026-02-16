package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a PostgreSQL connection pool from the given connection string.
func NewPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	if connString == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}
