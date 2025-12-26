package expenserepo

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

func (repo *Repository) Insert(ctx context.Context, newExpense model.Expense) (id int32, err error) {
	builder := squirrel.Insert("expenses").
		Columns(
			"user_id",
			"amount_idr",
			"description",
			"receipt_url",
			"status",
			"submitted_at",
			"processed_at",
		).
		Values(
			newExpense.UserId,
			newExpense.AmountIdr,
			newExpense.Description,
			newExpense.ReceiptUrl,
			newExpense.Status,
			newExpense.SubmittedAt,
			newExpense.ProcessedAt,
		)
	builder = builder.Prefix(
		"RETURNING id",
	)

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return id, err
	}

	if err := repo.querier.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return id, err
	}

	return
}

func (repo *Repository) SelectOneExpense(ctx context.Context, id int32) (expense model.Expense, err error) {
	builder := squirrel.Select(
		"id",
		"user_id",
		"amount_idr",
		"description",
		"receipt_url",
		"status",
		"submitted_at",
		"processed_at",
	).
		From("expenses")
	builder = builder.Where("id = ?", id)
	builder = builder.Limit(1)

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return expense, err
	}

	if err := repo.querier.QueryRow(ctx, query, args...).StructScan(&expense); err != nil {
		return expense, err
	}

	return
}
