package database

import (
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

type Querier interface {
}
