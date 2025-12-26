package authsv

import (
	"context"

	"github.com/arf-dev/mekari-test/model"
)

type UserRepository interface {
	SelectOneUser(ctx context.Context, email string) (model.User, error)
}
