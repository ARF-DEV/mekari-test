package model

import (
	"time"

	"github.com/arf-dev/mekari-test/pkg/httputils/response"
)

type Expense struct {
	Id          int32     `db:"id" json:"id"`
	UserId      int32     `db:"user_id" json:"user_id"`
	AmountIdr   int64     `db:"amount_idr" json:"amount_idr"`
	Description string    `db:"description" json:"description"`
	ReceiptUrl  string    `db:"receipt_url" json:"receipt_url"`
	Status      string    `db:"status" json:"status"`
	SubmittedAt time.Time `db:"submitted_at" json:"submitted_at"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at"`
}

type CreateExpenseRequest struct {
	UserId      int32     `json:"user_id"`
	AmountIdr   int64     `json:"amount_idr"`
	Description string    `json:"description"`
	ReceiptUrl  string    `json:"receipt_url"`
	SubmittedAt time.Time `json:"submitted_at"`
	ProcessedAt time.Time `json:"processed_at"`
}

type CreateExpenseResponse struct {
	response.BaseResponse
	Data int32 `json:"data"`
}

type GetExpenseRequest struct {
	Id int32 `path:"id"`
}

type GetExpenseResponse struct {
	response.BaseResponse
	Data Expense `json:"data"`
}

type GetExpenseListRequest struct {
	Page int64 `schema:"page"`
	Size int64 `schema:"size"`
}

type GetExpenseListResponse struct {
	response.BaseResponse
	Data []Expense `json:"data"`
}
