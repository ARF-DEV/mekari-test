package api

import (
	"github.com/arf-dev/mekari-test/handler/authhandlr"
	"github.com/arf-dev/mekari-test/handler/healthhandlr"
)

type API struct {
	HealthCheck *healthhandlr.Handler
	Auth        *authhandlr.Handler
}

func New() *API {
	return &API{
		HealthCheck: healthhandlr.New(),
		Auth:        authhandlr.New(),
	}
}
