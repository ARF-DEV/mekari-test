package userrepo

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/arf-dev/mekari-test/database"
	"github.com/arf-dev/mekari-test/model"
)

type Repository struct {
	querier database.Querier
}

func New(querier database.Querier) *Repository {
	return &Repository{
		querier: querier,
	}
}

func (repo *Repository) SelectOneUser(ctx context.Context, email string) (user model.User, err error) {
	builder := squirrel.Select("id", "email", "name", "role", "created_at").From("users")
	builder = builder.Where("email = ?", email)
	builder = builder.Limit(1)

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return user, err
	}

	if err := repo.querier.QueryRow(ctx, query, args...).StructScan(&user); err != nil {
		return user, err
	}
	return
}
