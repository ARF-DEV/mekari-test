package expensesv

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	config       *config.Config
	expenseRepo  ExpenseRepository
	approvalRepo ApprovalRepository
}

func New(config *config.Config, expenseRepo ExpenseRepository, approvalRepo ApprovalRepository) *Service {
	return &Service{
		config:       config,
		expenseRepo:  expenseRepo,
		approvalRepo: approvalRepo,
	}
}

func (service *Service) GetExpense(ctx context.Context, req model.GetExpenseRequest) (resp model.GetExpenseResponse, err error) {
	expense, err := service.expenseRepo.SelectOneExpense(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return resp, apierror.ErrResourceNotFound
		}
		log.Log().Err(err).Msg("error on GetExpense.SelectOneExpense")
		return resp, apierror.ErrInternalServer
	}
	resp.Data = expense
	return resp, nil
}

func (service *Service) GetExpenseList(ctx context.Context, req model.GetExpenseListRequest) (resp model.GetExpenseListResponse, err error) {
	expenses, err := service.expenseRepo.SelectExpense(ctx, req.Page, req.Size, req.Status)
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

func (service *Service) CreateExpense(ctx context.Context, req model.CreateExpenseRequest) (int32, error) {
	userData := ctxutils.GetUserDataFromCtx(ctx)
	userId := userData.UserId
	status := "pending"
	isAutoApproved := false
	if req.AmountIdr < 1000000 {
		status = "approved"
		isAutoApproved = true
	}

	now := timeNow()
	expenseId, err := service.expenseRepo.Insert(
		ctx,
		model.Expense{
			UserId:         userId,
			AmountIdr:      req.AmountIdr,
			Description:    req.Description,
			ReceiptUrl:     req.ReceiptUrl,
			Status:         status,
			SubmittedAt:    now,
			ProcessedAt:    now,
			IsAutoApproved: isAutoApproved,
		},
	)
	if err != nil {
		log.Log().Err(err).Msgf("error on CreateExpense.Insert")
		return 0, err
	}

	if isAutoApproved {
		go service.processedPayment(ctx, expenseId, req.AmountIdr)
	}

	return expenseId, nil
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
		ApproverId: userData.UserId,
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
	err := service.doPayment(noCancelCtx, model.PaymentRequest{
		Amount:     amount,
		ExternalId: uuid.NewString(),
	})
	if err != nil {
		log.Log().Err(err).Msgf("error when doing payment to a 3rd party for expense with id %d", expenseId)
		return
	}
}
func (service *Service) doPayment(ctx context.Context, paymentRequest model.PaymentRequest) error {
	const (
		max_retry int           = 2
		timeoff   time.Duration = time.Second * 5
	)
	retryCount := max_retry
	retryTimeoff := timeoff
	for ; retryCount > 0; retryCount-- {
		paymentEndpoint := service.config.PAYMENT_GATEWAY_URL + "/v1/payments"
		requestBody, _ := json.Marshal(paymentRequest)
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			paymentEndpoint,
			bytes.NewBuffer(requestBody),
		)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 30 * time.Second}

		log.Log().Msgf("hit endpoint %s", paymentEndpoint)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			if resp.StatusCode == 429 || resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != 501) {
				// retry if status code is 429, 0, or >= 500 expect 501
				time.Sleep(retryTimeoff)
				retryTimeoff *= 2
				continue
			}
			return fmt.Errorf("unexpected status: %s", resp.Status)
		}

		// if successful then break out of retry
		break
	}
	return nil
}
