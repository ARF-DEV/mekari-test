package expensesv

import (
	"context"

	"github.com/arf-dev/mekari-test/model"
)

type ExpenseRepository interface {
	Insert(ctx context.Context, newExpense model.Expense) (id int32, err error)
	SelectOneExpense(ctx context.Context, id int32) (expense model.Expense, err error)
	SelectExpense(ctx context.Context, page int64, size int64) (expenses []model.Expense, err error)
}
