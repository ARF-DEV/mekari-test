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

func (repo *Repository) Insert(ctx context.Context, newExpense model.Expense) (expense model.Expense, err error) {
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
		"RETURNING id, user_id, amount_idr, description, receipt_url, status, submitted_at, processed_at",
	)

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return expense, err
	}

	if err := repo.querier.QueryRow(ctx, query, args...).StructScan(&expense); err != nil {
		return expense, err
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

func (repo *Repository) SelectExpense(ctx context.Context, page, size int64, status string) (expenses []model.Expense, err error) {
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
	if status != "" {
		builder = builder.Where("status = ?", status)
	}

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

func (repo *Repository) Update(ctx context.Context, id int32, updateFunc func(expense *model.Expense)) error {
	expense, err := repo.SelectOneExpense(ctx, id)
	if err != nil {
		return err
	}

	updateFunc(&expense)

	builder := squirrel.Update("expenses").SetMap(
		map[string]interface{}{
			"amount_idr":   expense.AmountIdr,
			"description":  expense.Description,
			"receipt_url":  expense.ReceiptUrl,
			"status":       expense.Status,
			"submitted_at": expense.SubmittedAt,
			"processed_at": expense.ProcessedAt,
		},
	).Where("id = ?", id)

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	if _, err := repo.querier.Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil
}
