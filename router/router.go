package router

import (
	"github.com/arf-dev/mekari-test/api"
	"github.com/go-chi/chi/v5"
)

func New() *chi.Mux {
	api := api.New()
	chiMux := chi.NewMux()

	chiMux.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// login, health check (public)
			r.Get("/health", api.HealthCheck.Ping)
			r.Get("/auth/login", api.Auth.Login)
		})
		r.Group(func(r chi.Router) {
			// expenses (private)
		})
	})
	return chiMux
}
