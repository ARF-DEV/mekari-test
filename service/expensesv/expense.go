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

type Service struct {
	config      *config.Config
	expenseRepo ExpenseRepository
}

func New(config *config.Config, expenseRepo ExpenseRepository) *Service {
	return &Service{
		config:      config,
		expenseRepo: expenseRepo,
	}
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

func (service *Service) CreateExpense(ctx context.Context, req model.CreateExpenseRequest) (int32, error) {
	userData := ctxutils.GetUserDataFromCtx(ctx)
	userId := userData.UserId
	status := "awaiting-approval"
	isAutoApproved := false
	if req.AmountIdr < 1000000 {
		status = "approved"
		isAutoApproved = true
	}

	expenseId, err := service.expenseRepo.Insert(
		ctx,
		model.Expense{
			UserId:         userId,
			AmountIdr:      req.AmountIdr,
			Description:    req.Description,
			ReceiptUrl:     req.ReceiptUrl,
			Status:         status,
			SubmittedAt:    time.Now(),
			ProcessedAt:    time.Time{}, // zero value
			IsAutoApproved: isAutoApproved,
		},
	)
	if err != nil {
		return 0, err
	}

	if isAutoApproved {
		go func(expenseId int32) {
			noCancelCtx := context.WithoutCancel(ctx)
			err := service.doPayment(noCancelCtx, model.PaymentRequest{
				Amount:     req.AmountIdr,
				ExternalId: uuid.NewString(),
			})
			if err != nil {
				log.Log().Err(err).Msgf("error when doing payment to a 3rd party for expense with id %d", expenseId)
				return
			}

			if err = service.expenseRepo.Update(
				noCancelCtx,
				expenseId,
				func(expense *model.Expense) {
					expense.Status = "completed"
					expense.ProcessedAt = time.Now()
				},
			); err != nil {
				log.Log().Err(err).Msgf("error when updating expense with id %d", expenseId)
				return
			}
		}(expenseId)
	}

	return expenseId, nil
}

func (service *Service) doPayment(ctx context.Context, paymentRequest model.PaymentRequest) error {
	const (
		max_retry int           = 2
		timeoff   time.Duration = time.Second * 5
	)
	retryCount := max_retry
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
				continue
			}
			return fmt.Errorf("unexpected status: %s", resp.Status)
		}

		// if successfull then break out of retry
		break
	}
	return nil
}
