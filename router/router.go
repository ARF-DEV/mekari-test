package router

import (
	"github.com/arf-dev/mekari-test/api"
	"github.com/arf-dev/mekari-test/config"
	"github.com/arf-dev/mekari-test/database"
	"github.com/go-chi/chi/v5"
)

func New(config *config.Config) (*chi.Mux, error) {
	database, err := database.New(config.DB_MASTER)
	if err != nil {
		return nil, err
	}

	api := api.New(config, database)
	chiMux := chi.NewMux()

	chiMux.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// login, health check (public)
			r.Get("/health", api.HealthCheckHandlr.Ping)
			r.Post("/auth/login", api.AuthHandlr.Login)
		})
		r.Group(func(r chi.Router) {
			// expenses (private)
		})
	})
	return chiMux, nil
}
