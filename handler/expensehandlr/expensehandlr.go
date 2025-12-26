package expensehandlr

import "github.com/arf-dev/mekari-test/service/expensesv"

type Handler struct {
	expenseServ *expensesv.Service
}

func New(expenseServ *expensesv.Service) *Handler {
	return &Handler{
		expenseServ: expenseServ,
	}
}
