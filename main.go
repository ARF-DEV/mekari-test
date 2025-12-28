package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arf-dev/mekari-test/config"
	"github.com/arf-dev/mekari-test/router"
	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

func main() {
	env, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error when loading the config")
	}
	mux, err := router.New(env)
	if err != nil {
		log.Fatal().Err(err).Msg("error when initializing mux")
	}

	server := http.Server{
		Addr:         ":20000",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	runServer := func(errChan chan<- error) {
		log.Log().Msgf("server running on port 20000")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- err
		}
	}

	serverErrorChan := make(chan error, 1)
	go runServer(serverErrorChan)

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrorChan:
		log.Fatal().Err(err).Msg("error from server start")
	case <-exitChan:
		log.Log().Msg("signal received, shutting down now!")
	}

	exitCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	if err := server.Shutdown(exitCtx); err != nil {
		log.Fatal().Err(err).Msg("error from server server shutdown")
	}

	log.Log().Msg("server has gracefully shutdown")
}
