package userrepo

import "github.com/arf-dev/mekari-test/database"

type Repository struct {
	querier database.Querier
}

func New(querier database.Querier) *Repository {
	return &Repository{
		querier: querier,
	}
}
