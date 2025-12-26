package authsv

import (
	"context"

	"github.com/arf-dev/mekari-test/model"
)

// define service needs in interface
type UserRepository interface {
	SelectOneUser(ctx context.Context, email string) (model.User, error)
}
type Service struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (service *Service) AuthenticateUser(ctx context.Context, req model.LoginRequest) (token string, err error) {
	return
}
