package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

func New(url string) (*Database, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}
func (manager *Database) QueryRow(ctx context.Context, query string, args ...any) *sqlx.Row {
	return manager.db.QueryRowxContext(ctx, query, args...)
}

type Querier interface {
	QueryRow(ctx context.Context, query string, args ...any) *sqlx.Row
	// Query(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
}
