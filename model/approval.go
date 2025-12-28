package model

import "time"

type Approval struct {
	Id         int32     `db:"id"`
	ExpenseId  int32     `db:"expense_id"`
	ApproverId *int32    `db:"approver_id"`
	Status     string    `db:"status"`
	Notes      string    `db:"notes"`
	CreatedAt  time.Time `db:"created_at"`
}
