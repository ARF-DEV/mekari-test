package expensesv

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/arf-dev/mekari-test/config"
	"github.com/arf-dev/mekari-test/model"
	"github.com/arf-dev/mekari-test/pkg/ctxutils"
	"github.com/arf-dev/mekari-test/pkg/httputils/apierror"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var timeNow = time.Now

type Service struct {
	config          *config.Config
	expenseRepo     ExpenseRepository
	approvalRepo    ApprovalRepository
	paymentOutbound PaymentOutbound
}

func New(config *config.Config, expenseRepo ExpenseRepository, approvalRepo ApprovalRepository, paymentOutbound PaymentOutbound) *Service {
	return &Service{
		config:          config,
		expenseRepo:     expenseRepo,
		approvalRepo:    approvalRepo,
		paymentOutbound: paymentOutbound,
	}
}

func (service *Service) GetExpense(ctx context.Context, req model.GetExpenseRequest) (resp model.GetExpenseResponse, err error) {
	userData := ctxutils.GetUserDataFromCtx(ctx)

	expense, err := service.expenseRepo.SelectOneExpense(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return resp, apierror.ErrResourceNotFound
		}
		log.Log().Err(err).Msg("error on GetExpense.SelectOneExpense")
		return resp, apierror.ErrInternalServer
	}

	if !userData.IsManager() && expense.UserId != userData.UserId {
		return resp, apierror.ErrUnauthorized
	}

	resp.Data = expense
	return resp, nil
}

func (service *Service) GetExpenseList(ctx context.Context, req model.GetExpenseListRequest) (resp model.GetExpenseListResponse, err error) {
	var userId *int32 = nil

	userData := ctxutils.GetUserDataFromCtx(ctx)
	if !userData.IsManager() {
		userId = &userData.UserId
	}

	expenses, err := service.expenseRepo.SelectExpense(ctx, req.Page, req.Size, req.Status, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return resp, apierror.ErrResourceNotFound
		}
		log.Log().Err(err).Msg("error on GetExpenseList.SelectExpense")
		return resp, apierror.ErrInternalServer
	}
	resp.Data = expenses
	return resp, nil
}

func (service *Service) CreateExpense(ctx context.Context, req model.CreateExpenseRequest) (model.CreateExpenseResponseData, error) {
	userData := ctxutils.GetUserDataFromCtx(ctx)
	userId := userData.UserId

	status := "pending"
	now := timeNow()
	processedAt := time.Time{} // zero value

	isAutoApproved := req.IsAutoApproved()
	if isAutoApproved {
		status = "auto-approved"
		isAutoApproved = true
		processedAt = now
	}

	expense, err := service.expenseRepo.Insert(
		ctx,
		model.Expense{
			UserId:      userId,
			AmountIdr:   req.AmountIdr,
			Description: req.Description,
			ReceiptUrl:  req.ReceiptUrl,
			Status:      status,
			SubmittedAt: now,
			ProcessedAt: processedAt,
		},
	)
	if err != nil {
		log.Log().Err(err).Msgf("error on CreateExpense.expenseRepo.Insert")
		return model.CreateExpenseResponseData{}, err
	}

	if isAutoApproved {
		if _, err = service.approvalRepo.Insert(ctx, model.Approval{
			ExpenseId:  expense.Id,
			ApproverId: nil,
			Status:     status,
			CreatedAt:  timeNow(),
			Notes:      "auto approved",
		}); err != nil {
			log.Log().Err(err).Msgf("error on CreateExpense.approvalRepo.Insert")
			return model.CreateExpenseResponseData{}, err
		}
		go service.processedPayment(ctx, expense.Id, expense.AmountIdr)
	}

	return model.CreateExpenseResponseData{
		Id:               expense.Id,
		AmountIdr:        expense.AmountIdr,
		Description:      expense.Description,
		Status:           expense.Status,
		RequiresApproval: !isAutoApproved,
		AutoApproved:     isAutoApproved,
	}, nil
}

func (service *Service) UpdateExpense(ctx context.Context, req model.UpdateExpenseRequest) error {
	userData := ctxutils.GetUserDataFromCtx(ctx)

	status := ""
	switch req.Status {
	case "reject":
		status = "rejected"
	case "approve":
		status = "approved"
	default:
		log.Log().Msgf("cannot update expense to status=%s", req.Status)
		return apierror.ErrBadRequest
	}

	expense, err := service.expenseRepo.SelectOneExpense(ctx, req.Id)
	if err != nil {
		log.Log().Err(err).Msgf("error on UpdateExpense.SelectOneExpense")
		return err
	}

	if err = service.expenseRepo.Update(ctx, req.Id, func(expense *model.Expense) {
		expense.Status = status
		expense.ProcessedAt = timeNow()
	}); err != nil {
		log.Log().Err(err).Msgf("error on UpdateExpense.Update")
		return err
	}

	if _, err = service.approvalRepo.Insert(ctx, model.Approval{
		ExpenseId:  expense.Id,
		ApproverId: &userData.UserId,
		Status:     status,
		CreatedAt:  timeNow(),
		Notes:      req.Notes,
	}); err != nil {
		log.Log().Err(err).Msgf("error on UpdateExpense.Insert")
		return err
	}

	if status == "approved" {
		go service.processedPayment(ctx, expense.Id, expense.AmountIdr)
	}
	return err
}

func (service *Service) processedPayment(ctx context.Context, expenseId int32, amount int64) {
	noCancelCtx := context.WithoutCancel(ctx)
	err := service.paymentOutbound.DoPayment(noCancelCtx, model.PaymentRequest{
		Amount:     amount,
		ExternalId: uuid.NewString(),
	})
	if err != nil {
		log.Log().Err(err).Msgf("error when doing payment to a 3rd party for expense with id %d", expenseId)
		return
	}
}
