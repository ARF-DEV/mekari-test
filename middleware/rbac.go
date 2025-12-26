package middleware

import (
	"net/http"

	"github.com/arf-dev/mekari-test/pkg/ctxutils"
	"github.com/arf-dev/mekari-test/pkg/httputils/apierror"
	"github.com/arf-dev/mekari-test/pkg/httputils/response"
	"github.com/rs/zerolog/log"
)

// RBAC
func (m *Middleware) AccessWithRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userData := ctxutils.GetUserDataFromCtx(r.Context())
			if userData.Role != role {
				log.Log().Msgf("error on RBAC middleware: user role is %s", userData.Role)
				response.Send(w, "", nil, apierror.ErrUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
