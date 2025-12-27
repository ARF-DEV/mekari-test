package expensesv

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/arf-dev/mekari-test/model"
	"github.com/arf-dev/mekari-test/pkg/ctxutils"
	"github.com/arf-dev/mekari-test/pkg/httputils/apierror"
	"github.com/rs/zerolog/log"
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
	userData := ctxutils.GetUserDataFromCtx(ctx)
	userId := userData.UserId
	status := "pending"
	if req.AmountIdr < 1000000 {
		status = "approved"
		// TODO: hit 3rd party api for auto-approval case
	}

	return service.expenseRepo.Insert(
		ctx,
		model.Expense{
			UserId:      userId,
			AmountIdr:   req.AmountIdr,
			Description: req.Description,
			ReceiptUrl:  req.ReceiptUrl,
			Status:      status,
			SubmittedAt: time.Now(),
			ProcessedAt: time.Time{}, // zero
		},
	)
}

func (service *Service) GetExpense(ctx context.Context, req model.GetExpenseRequest) (resp model.GetExpenseResponse, err error) {
	expense, err := service.expenseRepo.SelectOneExpense(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return resp, apierror.ErrResourceNotFound
		}
		log.Log().Err(err).Msg("error on SelectOneExpense")
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
		log.Log().Err(err).Msg("error on SelectExpense")
		return resp, apierror.ErrInternalServer
	}
	resp.Data = expenses
	return resp, nil
}
