package api

import (
	"github.com/arf-dev/mekari-test/config"
	"github.com/arf-dev/mekari-test/database"
	"github.com/arf-dev/mekari-test/handler/authhandlr"
	"github.com/arf-dev/mekari-test/handler/expensehandlr"
	"github.com/arf-dev/mekari-test/handler/healthhandlr"
	"github.com/arf-dev/mekari-test/repository/expenserepo"
	"github.com/arf-dev/mekari-test/repository/userrepo"
	"github.com/arf-dev/mekari-test/service/authsv"
	"github.com/arf-dev/mekari-test/service/expensesv"
)

type API struct {
	HealthCheckHandlr *healthhandlr.Handler
	AuthHandlr        *authhandlr.Handler
	ExpenseHandlr     *expensehandlr.Handler
}

func New(config *config.Config, database *database.Database) *API {
	userRepo := userrepo.New(database)
	expenseRepo := expenserepo.New(database)

	authServ := authsv.New(config, userRepo)
	expensServ := expensesv.New(config, expenseRepo)

	healthHandlr := healthhandlr.New()
	authHandlr := authhandlr.New(authServ)
	expenseHandlr := expensehandlr.New(expensServ)

	return &API{
		HealthCheckHandlr: healthHandlr,
		AuthHandlr:        authHandlr,
		ExpenseHandlr:     expenseHandlr,
	}
}
