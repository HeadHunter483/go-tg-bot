package repository

import (
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrNotFound = errors.New("not found")

type repository struct {
	pool *pgxpool.Pool
}

// New creates a new repository instance.
func New(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}
