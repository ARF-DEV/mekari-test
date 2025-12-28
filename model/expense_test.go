package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAutoApproved(t *testing.T) {
	t.Run("return false if amount >= 1000000", func(t *testing.T) {
		req := CreateExpenseRequest{
			AmountIdr: 1000000,
		}
		assert.Equal(t, false, req.IsAutoApproved())
	})
	t.Run("return true if amount < 1000000", func(t *testing.T) {
		req := CreateExpenseRequest{
			AmountIdr: 900000,
		}
		assert.Equal(t, true, req.IsAutoApproved())
	})
}
