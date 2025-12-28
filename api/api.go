package api

import (
	"github.com/arf-dev/mekari-test/config"
	"github.com/arf-dev/mekari-test/database"
	"github.com/arf-dev/mekari-test/handler/authhandlr"
	"github.com/arf-dev/mekari-test/handler/expensehandlr"
	"github.com/arf-dev/mekari-test/handler/healthhandlr"
	"github.com/arf-dev/mekari-test/repository/approvalrepo"
	"github.com/arf-dev/mekari-test/repository/expenserepo"
	"github.com/arf-dev/mekari-test/repository/userrepo"
	"github.com/arf-dev/mekari-test/service/authsv"
	"github.com/arf-dev/mekari-test/service/expensesv"
)

//	@title			Expense Management API
//	@version		1.0
//	@description	This is a expense management API docs.

//	@host		localhost:20000
//	@BasePath	/api

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Bearer token
type API struct {
	HealthCheckHandlr *healthhandlr.Handler
	AuthHandlr        *authhandlr.Handler
	ExpenseHandlr     *expensehandlr.Handler
}

func New(config *config.Config, database *database.Database) *API {
	userRepo := userrepo.New(database)
	expenseRepo := expenserepo.New(database)
	approvalRepo := approvalrepo.New(database)

	authServ := authsv.New(config, userRepo)
	expensServ := expensesv.New(config, expenseRepo, approvalRepo)

	healthHandlr := healthhandlr.New()
	authHandlr := authhandlr.New(authServ)
	expenseHandlr := expensehandlr.New(expensServ)

	return &API{
		HealthCheckHandlr: healthHandlr,
		AuthHandlr:        authHandlr,
		ExpenseHandlr:     expenseHandlr,
	}
}
