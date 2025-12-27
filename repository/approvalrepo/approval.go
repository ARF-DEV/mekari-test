package approvalrepo

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

func (repo *Repository) Insert(ctx context.Context, newApproval model.Approval) (id int32, err error) {
	builder := squirrel.Insert("approvals").
		Columns(
			"expense_id",
			"approver_id",
			"status",
			"notes",
			"created_at",
		).
		Values(
			newApproval.ExpenseId,
			newApproval.ApproverId,
			newApproval.Status,
			newApproval.Notes,
			newApproval.CreatedAt,
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

func (repo *Repository) SelectOneApproval(ctx context.Context, id int32) (approval model.Approval, err error) {
	builder := squirrel.Select(
		"id",
		"expense_id",
		"approver_id",
		"status",
		"notes",
		"created_at",
	).
		From("approvals")
	builder = builder.Where("id = ?", id)
	builder = builder.Limit(1)

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return approval, err
	}

	if err := repo.querier.QueryRow(ctx, query, args...).StructScan(&approval); err != nil {
		return approval, err
	}

	return
}

func (repo *Repository) Update(ctx context.Context, id int32, updateFunc func(expense *model.Approval)) error {
	expense, err := repo.SelectOneApproval(ctx, id)
	if err != nil {
		return err
	}

	updateFunc(&expense)
	builder := squirrel.Update("approvals").SetMap(
		map[string]interface{}{
			"status": expense.Status,
			"notes":  expense.Notes,
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
