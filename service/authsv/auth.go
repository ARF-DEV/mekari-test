package authsv

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/arf-dev/mekari-test/config"
	"github.com/arf-dev/mekari-test/model"
	"github.com/arf-dev/mekari-test/pkg/authentication"
	"github.com/arf-dev/mekari-test/pkg/httputils/apierror"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type Service struct {
	config   *config.Config
	userRepo UserRepository
}

func New(config *config.Config, userRepo UserRepository) *Service {
	return &Service{
		config:   config,
		userRepo: userRepo,
	}
}

// if user exists, then create token
// otherwise, return error
func (service *Service) AuthenticateUser(ctx context.Context, req model.LoginRequest) (token string, err error) {
	user, err := service.userRepo.SelectOneUser(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return token, apierror.Error("unauthenticated", "Authentication failed", http.StatusUnauthorized)
		}
		log.Log().Err(err).Msg("error on SelectOneUser")
		return token, apierror.Error("internal_server_error", "Internal Server Error", http.StatusInternalServerError)
	}

	claims := authentication.Claims{
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ems-api",
			Subject:   "ems-fe",
		},
	}
	token, err = authentication.GenerateToken(service.config.JWT_SECRET, claims)
	if err != nil {
		log.Log().Err(err).Msg("error when generating token")
		return
	}
	return
}
