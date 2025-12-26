package middleware

import (
	"net/http"
	"strings"

	"github.com/arf-dev/mekari-test/config"
	"github.com/arf-dev/mekari-test/pkg/authentication"
	"github.com/arf-dev/mekari-test/pkg/ctxutils"
	"github.com/arf-dev/mekari-test/pkg/httputils/apierror"
	"github.com/arf-dev/mekari-test/pkg/httputils/response"
	"github.com/rs/zerolog/log"
)

type Middleware struct {
	config *config.Config
}

func New(config *config.Config) *Middleware {
	return &Middleware{
		config: config,
	}
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer ") {
			log.Log().Msgf("error Authenticate middleware: no Bearer tag")
			response.Send(w, "", nil, apierror.ErrUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := authentication.ParseClaimsFromToken(m.config.JWT_SECRET, token)
		if err != nil {
			response.Send(w, "", nil, apierror.ErrUnauthorized)
			return
		}
		r = r.WithContext(
			ctxutils.CtxWithUserData(
				r.Context(),
				ctxutils.UserData{
					Email: claims.Email,
					Role:  claims.Role,
				},
			),
		)
		next.ServeHTTP(w, r)
	})
}
