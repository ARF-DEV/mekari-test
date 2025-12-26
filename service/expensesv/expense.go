package expensesv

type Service struct {
	expenseRepo ExpenseRepository
}

func New(expenseRepo ExpenseRepository) *Service {
	return &Service{
		expenseRepo: expenseRepo,
	}
}
