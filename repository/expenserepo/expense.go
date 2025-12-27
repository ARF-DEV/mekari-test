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
	builder = builder.Suffix(
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

func (repo *Repository) SelectExpense(ctx context.Context, page, size int64) (expenses []model.Expense, err error) {
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
	builder = builder.OrderBy("submitted_at desc")
	builder = builder.Offset(uint64((page - 1) * size))
	builder = builder.Limit(uint64(size))

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return expenses, err
	}

	rows, err := repo.querier.Query(ctx, query, args...)
	if err != nil {
		return expenses, err
	}
	defer rows.Close()

	for rows.Next() {
		var expense model.Expense
		if err = rows.StructScan(&expense); err != nil {
			return expenses, err
		}
		expenses = append(expenses, expense)
	}
	return
}
