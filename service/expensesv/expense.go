package expensesv

import (
	"context"
	"database/sql"
	"errors"

	"github.com/arf-dev/mekari-test/model"
	"github.com/arf-dev/mekari-test/pkg/httputils/apierror"
)

type Service struct {
	expenseRepo ExpenseRepository
}

func New(expenseRepo ExpenseRepository) *Service {
	return &Service{
		expenseRepo: expenseRepo,
	}
}

func (service *Service) CreateExpense(ctx context.Context, req model.CreateExpenseRequest) (int32, error) {
	// TODO: Auto-approval logic
	// TODO: hit 3rd party api for auto-approval case
	status := "pending"
	return service.expenseRepo.Insert(
		ctx,
		model.Expense{
			UserId:      req.UserId,
			AmountIdr:   req.AmountIdr,
			Description: req.Description,
			ReceiptUrl:  req.ReceiptUrl,
			Status:      status,
			SubmittedAt: req.SubmittedAt,
			ProcessedAt: req.ProcessedAt,
		},
	)
}

func (service *Service) GetExpense(ctx context.Context, req model.GetExpenseRequest) (resp model.GetExpenseResponse, err error) {
	expense, err := service.expenseRepo.SelectOneExpense(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return resp, apierror.ErrResourceNotFound
		}
		return resp, apierror.ErrInternalServer
	}
	resp.Data = expense
	return resp, nil
}

func (service *Service) GetExpenseList(ctx context.Context, req model.GetExpenseListRequest) (resp model.GetExpenseListResponse, err error) {
	expenses, err := service.expenseRepo.SelectExpense(ctx, req.Page, req.Size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return resp, apierror.ErrResourceNotFound
		}
		return resp, apierror.ErrInternalServer
	}
	resp.Data = expenses
	return resp, nil
}
