package expensesv

import (
	"context"

	"github.com/arf-dev/mekari-test/model"
)

type ExpenseRepository interface {
	Insert(ctx context.Context, newExpense model.Expense) (expense model.Expense, err error)
	SelectOneExpense(ctx context.Context, id int32) (expense model.Expense, err error)
	SelectExpense(ctx context.Context, page int64, size int64, status string, userId *int32) (expenses []model.Expense, err error)
	Update(ctx context.Context, id int32, updateFunc func(expense *model.Expense)) error
}

type ApprovalRepository interface {
	Insert(ctx context.Context, newApproval model.Approval) (id int32, err error)
	SelectOneApproval(ctx context.Context, id int32) (approval model.Approval, err error)
	Update(ctx context.Context, id int32, updateFunc func(approval *model.Approval)) error
}
