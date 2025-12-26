package api

import (
	"github.com/arf-dev/mekari-test/config"
	"github.com/arf-dev/mekari-test/database"
	"github.com/arf-dev/mekari-test/handler/authhandlr"
	"github.com/arf-dev/mekari-test/handler/healthhandlr"
	"github.com/arf-dev/mekari-test/repository/userrepo"
	"github.com/arf-dev/mekari-test/service/authsv"
)

type API struct {
	HealthCheckHandlr *healthhandlr.Handler
	AuthHandlr        *authhandlr.Handler
}

func New(config *config.Config, database *database.Database) *API {
	userRepo := userrepo.New(database)
	authServ := authsv.New(config, userRepo)

	healthHandlr := healthhandlr.New()
	authHandlr := authhandlr.New(authServ)
	return &API{
		HealthCheckHandlr: healthHandlr,
		AuthHandlr:        authHandlr,
	}
}
