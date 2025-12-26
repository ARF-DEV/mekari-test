package model

import "time"

type Expense struct {
	Id          int32     `db:"id"`
	UserId      int32     `db:"user_id"`
	AmountIdr   int64     `db:"amount_idr"`
	Description string    `db:"description"`
	ReceiptUrl  string    `db:"receipt_url"`
	Status      string    `db:"status"`
	SubmittedAt time.Time `db:"submitted_at"`
	ProcessedAt time.Time `db:"processed_at"`
}
