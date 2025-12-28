package ctxutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsManager(t *testing.T) {
	t.Run("return false if user is not a manager", func(t *testing.T) {
		ctx := CtxWithUserData(context.Background(), UserData{
			UserId: 1,
			Role:   "user",
		})

		userData := GetUserDataFromCtx(ctx)
		assert.Equal(t, false, userData.IsManager())
	})
	t.Run("return true if user is a manager", func(t *testing.T) {
		ctx := CtxWithUserData(context.Background(), UserData{
			UserId: 1,
			Role:   "manager",
		})

		userData := GetUserDataFromCtx(ctx)
		assert.Equal(t, true, userData.IsManager())
	})

}
