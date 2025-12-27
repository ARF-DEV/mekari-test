package router

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/arf-dev/mekari-test/api"
	"github.com/arf-dev/mekari-test/config"
	"github.com/arf-dev/mekari-test/database"
	"github.com/arf-dev/mekari-test/middleware"
	"github.com/go-chi/chi/v5"
)

func New(config *config.Config) (*chi.Mux, error) {
	database, err := database.New(config.DB_MASTER)
	if err != nil {
		return nil, err
	}
	middlewareManager := middleware.New(config)

	api := api.New(config, database)
	chiMux := chi.NewMux()

	chiMux.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// login, health check (public)
			r.Get("/health", api.HealthCheckHandlr.Ping)
			r.Post("/auth/login", api.AuthHandlr.Login)
		})
		r.Group(func(r chi.Router) {
			r.Use(middlewareManager.Authenticate)
			// testing
			r.Get("/health/auth", api.HealthCheckHandlr.Ping)

			r.Get("/expenses", api.ExpenseHandlr.HandleGetExpenseList)
			r.Get("/expenses/{id}", api.ExpenseHandlr.HandleGetExpense)
			r.Post("/expenses", api.ExpenseHandlr.HandleCreateExpense)

			r.Group(func(r chi.Router) {
				r.Use(middlewareManager.AccessWithRole("user"))
				// endpoint for testing perpose
				r.Get("/health/auth/user", api.HealthCheckHandlr.Ping)
			})
			r.Group(func(r chi.Router) {
				r.Use(middlewareManager.AccessWithRole("manager"))
				// endpoint for testing perpose
				r.Get("/health/auth/manager", api.HealthCheckHandlr.Ping)
			})
			// expenses (private)
		})
	})
	printEndpoints(chiMux)
	return chiMux, nil
}

func printEndpoints(mux *chi.Mux) {
	chi.Walk(mux, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s] %s", method, route)
		if len(middlewares) > 0 {
			fmt.Printf(" (middlewares: %d)", len(middlewares))
		}
		if handler != nil {
			funcName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			fmt.Printf(" -> %s", funcName)
		}
		fmt.Println()
		return nil
	})
}
