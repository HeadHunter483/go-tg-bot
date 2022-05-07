package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// New creates a new connection pool to postgres db instance.
func New(ctx context.Context, dbInfo string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, dbInfo)
}
